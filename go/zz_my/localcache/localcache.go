package localcache

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"reflect"
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

type LocalCache[S comparable, T IModel] struct {
	db *sqlx.DB
	*Options
	initOnce *sync.Once

	lruCache *lru.Cache[S, T]
	noRows   map[S]struct{}
	vm       VersionManager

	// TODO[Dokiy] to be continued! (2025/11/26)
	// 添加联合唯一键查找
	tableName   string
	keyDbTag    string
	tagFieldIdx int
}

const (
	defaultSize = 10000
)

func NewLocalCache[S comparable, T IModel](db *sqlx.DB, vm VersionManager, keyDbTag string, opts ...LocalOption) LocalCache[S, T] {
	var t = new(T)
	var option = &Options{
		maxSize: defaultSize,
	}
	for _, opt := range opts {
		opt(option)
	}
	localCache := LocalCache[S, T]{
		db:          db,
		Options:     option,
		initOnce:    &sync.Once{},
		noRows:      make(map[S]struct{}),
		vm:          vm,
		tableName:   (*t).TableName(),
		keyDbTag:    keyDbTag,
		tagFieldIdx: -1,
	}
	if localCache.maxSize <= 0 {
		panic("LocalCache init: invalid config")
	}
	localCache.indexTagField(*t)

	return localCache
}

func (l *LocalCache[S, T]) IncrVersion(ctx context.Context) (string, error) {
	return l.vm.IncrVersion(ctx, l.tableName)
}
func (l *LocalCache[S, T]) GetLocal(ctx context.Context, s S) (t *T, err error) {
	l.initOnce.Do(l.init)
	if t, ok := l.lruCache.Peek(s); ok {
		return &t, nil
	}
	if _, ok := l.noRows[s]; ok {
		return nil, sql.ErrNoRows
	}

	var mdl T
	var args = []any{s}
	var sqlStr = fmt.Sprintf("SELECT * FROM %s t WHERE t.%s = ?", l.tableName, l.keyDbTag)
	if l.Options.withStatus {
		sqlStr = fmt.Sprintf("SELECT * FROM %s t WHERE t.%s = ? AND t.status = ?", l.tableName, l.keyDbTag)
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

	var hasIdxTag, hasStatusTag bool
	for j := 0; j < tValue.NumField(); j++ {
		tag, ok := tValue.Type().Field(j).Tag.Lookup("db")
		if !ok {
			continue
		}

		if tag == l.keyDbTag {
			if tValue.Field(j).Type().AssignableTo(reflect.TypeOf(*new(S))) {
				l.tagFieldIdx = j
				if !l.Options.withStatus || hasStatusTag {
					return
				} else {
					hasIdxTag = true
				}
			} else {
				panic("LocalCache load: cannot assign field to type S")
			}
		}
		if tag == "status" && l.Options.withStatus {
			if reflect.TypeOf(l.Options.status).AssignableTo(tValue.Field(j).Type()) {
				if hasIdxTag {
					return
				} else {
					hasStatusTag = true
				}
			} else {
				panic("LocalCache load: cannot assign status")
			}
		}
	}

	panic("LocalCache load: no field with db tag '" + l.keyDbTag + "' or 'status' found")
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
	l.lruCache.Add(tValue.Field(l.tagFieldIdx).Interface().(S), t)
}
