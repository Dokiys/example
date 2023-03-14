package redis

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
)

func TestRedisHello(t *testing.T) {
	ctx := context.Background()
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	defer client.Close()

	// Set & Get
	client.Set(ctx, "key", "value", 1*time.Hour)
	v := client.Get(ctx, "key")
	assert.Equal(t, "value", v.Val())

	// Not found
	_, err := client.Get(ctx, "key2").Result()
	assert.True(t, errors.Is(err, redis.Nil))
}

// 背景：并发的情况下我们可能希望通过某些条件来判断是否将值设置（并非简单的 SETNX）
// 比如有一个名为 version 的key，我们需要尝试给它赋值，如果所对应的值不存在时可以
// 成功存入，但值存在时，值只能增不能减。每次操作都将返回操作后的结果，并刷新key的过
// 期时间。
// 方案：可以利用redis支持的lua脚本，来做判断和存值的原子操作。
func TestRedisHelloLua(t *testing.T) {
	ctx := context.Background()
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	defer client.Close()

	src := `
local key = tostring(KEYS[1])
local oldValue = tonumber(redis.call("GET", key))
local newValue = tonumber(ARGV[1])
local expire = tonumber(ARGV[2])	-- 任何情况下缓存都需要设置过期时间！

if(newValue < 0)
then
	return error("invalid value set")
end

-- 如果不存在直接存
if(oldValue == nil)
then
    redis.call("SETEX", key, expire, newValue)
    return newValue
end

-- 存入最新值
if(newValue > oldValue)
then
    redis.call("SETEX", key, expire, newValue)
    return newValue
else
    redis.call("SETEX", key, expire, oldValue)
	return oldValue
end
`
	// 现将脚本加载到redis内存，如果给定的脚本已经在缓存里面了，那么不执行任何操作。
	// 脚本可以在缓存中保留无限长的时间，直到执行 SCRIPT FLUSH 为止。
	sha, err := redis.NewScript(src).Load(ctx, client).Result()
	if err != nil {
		t.Fatal(err)
	}

	keys, expire := []string{"version"}, 10
	v2 := []interface{}{"2", expire}
	result := client.EvalSha(ctx, sha, keys, v2...).Val()
	assert.Equal(t, int64(2), result)

	// 存入失败，依然取到最新值2
	v1 := []interface{}{"1", expire}
	result = client.EvalSha(ctx, sha, keys, v1...).Val()
	assert.Equal(t, int64(2), result)

	// 传值错误

	v := []interface{}{"-1", expire}
	_, err = client.EvalSha(ctx, sha, keys, v...).Result()
	t.Log(err) // => ERR user_script:9: invalid value set script: 251f1bfe40b199887ccc90ee6d0e6ad6863e31e6, on @user_script:9.
}
