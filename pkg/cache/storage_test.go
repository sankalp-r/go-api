package cache

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInMemoryStore(t *testing.T) {
	testCache := NewStorage()
	testEtag := Etag{
		Key:  "etag1",
		Data: []byte("testdata"),
	}
	testCache.Set("url1", testEtag)
	assert.Equal(t, testCache.Get("url1"), &testEtag)

	testCache.Delete("url1")
	assert.True(t, testCache.Get("url1") == nil)
}
