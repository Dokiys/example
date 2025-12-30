package localcache

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"reflect"
	"slices"
	"strings"
	"sync"
	"time"

	lru "github.com/hashicorp/golang-lru/v2"
	"github.com/jmoiron/sqlx"
)

type IModel interface {
	TableName() string
	PK() string
}

type VersionManager interface {
	Version(ctx context.Context, key string) (string, error)
	IncrVersion(ctx context.Context, key string) (string, error)
}
type Options struct {
	maxSize    int
	status     any
	withStatus bool
	pkTags     []string
}
type LocalOption func(*Options)

func WithMaxSize(maxSize int) LocalOption {
	return func(o *Options) {
		if maxSize <= 0 {
			panic("LocalCache invalid maxSize")
		}
		o.maxSize = maxSize
	}
}

func WithStatus(status any) LocalOption {
	return func(o *Options) {
		o.status = status
		o.withStatus = true
	}
}

type PKTagsOptions func(*Options)

func PkDbTags(tags ...string) PKTagsOptions {
	return func(o *Options) {
		o.pkTags = tags
	}
}

type LocalCache[S comparable, T IModel] struct {
	db *sqlx.DB
	*Options
	initOnce *sync.Once

	lruCache *lru.Cache[S, T]
	noRows   map[S]struct{}
	vm       VersionManager

	tableName     string
	keyDbTags     []string
	tagFieldIdxes []int
}

const (
	defaultSize = 10000
)

// NewLocalCache
// TODO[Dokiy] 修改范型 (2025/12/19)
// 目前不知道这里如何限制S必须是[x]any的数组。而且也不知道如何让[2]any的长度直接跟keyDbTags长度挂钩
// _ = NewLocalCache[[2]any, MockIModel](nil, nil, []string{})
func NewLocalCache[S comparable, T IModel](db *sqlx.DB, vm VersionManager, pkTags PKTagsOptions, opts ...LocalOption) LocalCache[S, T] {
	var t = new(T)
	var option = &Options{
		maxSize: defaultSize,
	}
	pkTags(option)
	for _, opt := range opts {
		opt(option)
	}
	localCache := LocalCache[S, T]{
		db:            db,
		Options:       option,
		initOnce:      &sync.Once{},
		noRows:        make(map[S]struct{}),
		vm:            vm,
		tableName:     (*t).TableName(),
		keyDbTags:     option.pkTags,
		tagFieldIdxes: make([]int, len(option.pkTags)),
	}
	if localCache.maxSize <= 0 {
		panic("LocalCache init: invalid config")
	}
	if len(localCache.pkTags) <= 0 {
		panic("LocalCache init: invalid keyDbTags length")
	}

	var s S
	typ := reflect.TypeOf(s)
	if typ.Kind() != reflect.Array || !typ.Comparable() {
		panic("LocalCache init: invalid key type")
	}
	if typ.Len() != len(localCache.pkTags) {
		panic("LocalCache init: invalid key type length which must be equal to keyDbTags length")
	}
	localCache.indexTagField(*t)

	return localCache
}

func (l *LocalCache[S, T]) IncrVersion(ctx context.Context) (string, error) {
	return l.vm.IncrVersion(ctx, l.tableName)
}

func (l *LocalCache[S, T]) GetLocal(ctx context.Context, params ...any) (t *T, err error) {
	l.initOnce.Do(l.init)

	// 尝试从缓存获取数据
	var s S
	{
		if reflect.ValueOf(&s).Len() != len(params) {
			panic("field value type is not assignable/convertible to array element type")
		}
		sv := reflect.ValueOf(&s).Elem()
		for i, param := range params {
			elem := sv.Index(i)
			fv := reflect.ValueOf(param)
			if fv.Type().AssignableTo(elem.Type()) {
				elem.Set(fv)
			} else if fv.Type().ConvertibleTo(elem.Type()) {
				elem.Set(fv.Convert(elem.Type()))
			} else {
				panic("field value type is not assignable/convertible to array element type")
			}
		}

		if t, ok := l.lruCache.Peek(s); ok {
			return &t, nil
		}
		if _, ok := l.noRows[s]; ok {
			return nil, sql.ErrNoRows
		}
	}

	var mdl T
	var args = params
	where := make([]string, 0, len(l.keyDbTags))
	for _, tag := range l.keyDbTags {
		where = append(where, fmt.Sprintf(" t.%s = ?", tag))
	}
	var sqlStr = fmt.Sprintf("SELECT * FROM %s t WHERE %s", l.tableName, strings.Join(where, " AND "))
	if l.Options.withStatus {
		sqlStr = fmt.Sprintf("SELECT * FROM %s t WHERE %s AND t.status = ?", l.tableName, strings.Join(where, " AND "))
		args = append(args, l.Options.status)
	}
	if err := l.db.GetContext(ctx, &mdl, sqlStr, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			l.noRows[s] = struct{}{}
			return nil, err
		}

		slog.Error(fmt.Sprintf("刷新 %s LocalCache查询单条数据失败", l.tableName))
		return &mdl, err
	}

	l.add(mdl)
	return &mdl, nil
}

func (l *LocalCache[S, T]) indexTagField(t T) {
	tValue := reflect.ValueOf(t)
	if tValue.Type().Kind() != reflect.Struct {
		panic("LocalCache load: provided value is not a struct")
	}

	var hasStatusTag bool
	var s S
	for j := 0; j < tValue.NumField(); j++ {
		tag, ok := tValue.Type().Field(j).Tag.Lookup("db")
		if !ok {
			continue
		}

		if idx := slices.Index(l.keyDbTags, tag); idx != -1 {
			// 检验能否存入S
			sv := reflect.ValueOf(&s).Elem()
			if tValue.Field(j).Type().AssignableTo(sv.Index(idx).Type()) || tValue.Field(j).Type().ConvertibleTo(sv.Index(idx).Type()) {
				// 保证tagFieldIdxes的顺序，在插入和读取数据的时候都比较key，顺序不同导致key不同。
				l.tagFieldIdxes[idx] = j
				if len(l.tagFieldIdxes) == len(l.keyDbTags) && (!l.Options.withStatus || hasStatusTag) {
					return
				}
			} else {
				panic("LocalCache load: cannot assign field to type S")
			}
		}
		if tag == "status" && l.Options.withStatus {
			if reflect.TypeOf(l.Options.status).AssignableTo(tValue.Field(j).Type()) {
				if len(l.tagFieldIdxes) == len(l.keyDbTags) {
					return
				} else {
					hasStatusTag = true
				}
			} else {
				panic("LocalCache load: cannot assign status")
			}
		}
	}

	panic("LocalCache load: no field with db tag '" + strings.Join(l.keyDbTags, ",") + "' or 'status' found")
}
func (l *LocalCache[S, T]) init() {
	l.lruCache, _ = lru.New[S, T](l.maxSize)

	var ctx = context.Background()
	var version string
	var flash = func() error {
		var mdlList []T
		var pageSize = l.lruCache.Len()
		if pageSize <= 0 {
			pageSize = l.maxSize
		}

		var args []any
		var sqlStr = fmt.Sprintf("SELECT * FROM %s t ORDER BY t.updated_at DESC LIMIT 0, ?", l.tableName)
		if l.Options.withStatus {
			sqlStr = fmt.Sprintf("SELECT * FROM %s t WHERE t.status = ? ORDER BY t.updated_at DESC LIMIT 0, ?", l.tableName)
			args = append(args, l.Options.status)
		}
		args = append(args, pageSize)

		if err := l.db.SelectContext(ctx, &mdlList, sqlStr, args...); err != nil {
			slog.Error(fmt.Sprintf("刷新 %s LocalCache查询数据失败", l.tableName))
			return err
		}

		l.noRows = make(map[S]struct{})
		l.lruCache.Purge()
		for _, t := range mdlList {
			l.add(t)
		}
		return nil
	}

	var flushCh = make(chan struct{})
	var initDone = make(chan struct{})
	var flushInterval = time.NewTicker(1 * time.Minute)
	var forceFlushInterval = time.NewTicker(5 * time.Minute)
	var once sync.Once
	go func() {
		for {
			select {
			case <-flushCh:
				newVersion, err := l.vm.Version(ctx, l.tableName)
				if err == nil && version != newVersion {
					if err := flash(); err == nil {
						version = newVersion
					}
					forceFlushInterval.Reset(5 * time.Minute)
				}

				once.Do(func() {
					close(initDone)
				})
			case <-forceFlushInterval.C:
				_ = flash()
				flushInterval.Reset(1 * time.Minute)
			}
		}
	}()
	go func() {
		flushCh <- struct{}{}
		for range flushInterval.C {
			flushCh <- struct{}{}
		}
	}()

	<-initDone
}
func (l *LocalCache[S, T]) add(t T) {
	tValue := reflect.ValueOf(t)

	var s S
	sv := reflect.ValueOf(&s).Elem() // sv 是 S 的可设置值（array）
	for i, idx := range l.tagFieldIdxes {
		// 取出 T 中对应字段的值
		fv := tValue.Field(idx)
		// 将字段值写入数组第 i 个元素
		elem := sv.Index(i)
		if fv.Type().AssignableTo(elem.Type()) {
			elem.Set(fv)
		} else if fv.Type().ConvertibleTo(elem.Type()) {
			elem.Set(fv.Convert(elem.Type()))
		} else {
			panic("LocalCache add: cannot assign field to type S")
		}
	}
	l.lruCache.Add(s, t)
}
