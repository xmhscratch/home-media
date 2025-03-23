package sys

import (
	"context"

	"github.com/golang/groupcache"
)

var SessionContext = context.Background()

var UUIDFormat = "%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x"
var UUIDNamespace = "6ba7b8109dad11d180b400c02fd430c8"

// RedisConnectionConfig ...
type RedisConnectionConfig struct {
	DB       string `json:"db"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// SQLiteConnectionConfig ...
type SQLiteConnectionConfig struct {
	Path string `json:"path"`
}

// ConfigFile comment
type ConfigFile struct {
	Debug    bool                   `json:"debug"`
	Port     string                 `json:"port"`
	HostName string                 `json:"hostName"`
	SQLite   SQLiteConnectionConfig `json:"sqlite"`
	Redis    RedisConnectionConfig  `json:"redis"`
	BinPath  string                 `json:"binPath"`
	TmpPath  string                 `json:"tmpPath"`
	RootPath string                 `json:"rootPath"`
	DataPath string                 `json:"dataPath"`
}

// Config comment
type Config struct {
	ConfigFile
	AppDir     string
	GroupCache *groupcache.Group
}
