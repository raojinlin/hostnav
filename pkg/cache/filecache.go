package cache

import (
	"encoding/gob"
	"fmt"
	"os"
	"path"
	"time"
)

const cacheFileName = "cache.gob"

type CacheItem[T interface{}] struct {
	Content  T         `json:"content"`
	ExpireAt time.Time `json:"expire_at"`
}

type FileCache[T interface{}] struct {
	CacheDir   string `json:"cache_dir"`
	DefaultTTL time.Duration
	cache      map[string]*CacheItem[T]
}

func NewFileCache[T interface{}](cacheDir string, defaultTTL time.Duration) *FileCache[T] {
	return &FileCache[T]{
		CacheDir:   cacheDir,
		DefaultTTL: defaultTTL,
		cache:      make(map[string]*CacheItem[T]),
	}
}

func (c *FileCache[T]) Get(key string) (*T, error) {
	cacheItem, ok := c.cache[key]
	if !ok {
		return nil, fmt.Errorf("cache key %s not found", key)
	}

	if !cacheItem.ExpireAt.After(time.Now()) {
		c.Del(key)
		return nil, fmt.Errorf("cache key %s expired", key)
	}

	return &cacheItem.Content, nil
}

func (c *FileCache[T]) Set(key string, content T, ttl time.Duration) {
	if ttl <= 0 {
		ttl = c.DefaultTTL
	}

	cacheItem := &CacheItem[T]{Content: content, ExpireAt: time.Now().Add(ttl)}
	c.cache[key] = cacheItem
}

func (c *FileCache[T]) Del(key string) {
	delete(c.cache, key)
}

func (c *FileCache[T]) Save() error {
	cacheDir, err := os.OpenFile(c.CacheDir, os.O_WRONLY, os.ModeDir)
	if os.IsNotExist(err) {
		if err = os.Mkdir(c.CacheDir, os.ModeDir); err != nil {
			return err
		}
	}

	defer cacheDir.Close()

	cacheFilePath := path.Join(c.CacheDir, cacheFileName)
	cacheFile, err := os.OpenFile(cacheFilePath, os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}

	defer cacheFile.Close()

	encoder := gob.NewEncoder(cacheFile)
	if err = encoder.Encode(c.cache); err != nil {
		return err
	}

	return nil
}

func (c *FileCache[T]) Load() error {
	cacheFilePath := path.Join(c.CacheDir, cacheFileName)
	cacheFile, err := os.Open(cacheFilePath)
	if err != nil {
		// cache file not exists
		if os.IsNotExist(err) {
			return nil
		}

		return fmt.Errorf("cloud not open cache file: %v", err)
	}

	defer cacheFile.Close()
	decoder := gob.NewDecoder(cacheFile)
	if err = decoder.Decode(&c.cache); err != nil {
		return fmt.Errorf("cloud not decode cache file: %v", err)
	}

	return nil
}
