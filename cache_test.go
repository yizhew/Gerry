package gerry

import (
	oauth2 "github.com/silenceper/wechat/oauth"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestSimplePutAndGet(t *testing.T) {
	DefaultCache.Put("key1", "value1")
	assert.Equal(t, "value1", DefaultCache.Get("key1"), "cache should not chang value")
}

func TestExpiration(t *testing.T) {
	DefaultCache.Set("key2", "value2", 10*time.Second)
	DefaultCache.Set("key3", "value3", 10*time.Second)
	DefaultCache.Set("key4", "value4", 10*time.Minute)

	time.Sleep(40 * time.Second)

	assert.False(t, DefaultCache.IsExist("key2"), "key2 should not be existing")
	assert.False(t, DefaultCache.IsExist("key3"), "key3 should not be existing")
	assert.Equal(t, "value4", DefaultCache.Get("key4"), "key4 should be found")
}

func TestUserInfoInCache(t *testing.T) {
	ui := oauth2.UserInfo{}
	ui.OpenID = "id1"
	ui.Nickname = "name1"

	DefaultCache.Put("userinfo", ui)
	ui2 := DefaultCache.Get("userinfo").(oauth2.UserInfo)

	assert.NotNil(t, ui2, "ui2 should be null")
	assert.Equal(t, "id1", ui2.OpenID, "openid should not be null")
	assert.Equal(t, "name1", ui2.Nickname, "nickname should not be null")
}
