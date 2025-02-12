package command

import (
	"bytes"
	"home-media/streaming/core"
	"io"

	"github.com/redis/go-redis/v9"
)

type SessionWriter struct {
	io.Writer
	SessionId string
	AttrName  string
	Config    *core.Config
	redis     *redis.Client
}

type CommandFrags struct {
	ExecBin     string `json:"executor"`
	Input       string `json:"input"`
	Start       string `json:"start"`
	Duration    string `json:"duration"`
	Output      string `json:"output"`
	DownloadURL string `json:"downloadUrl"`
	BaseURL     string `json:"baseURL"`
	RootDir     string `json:"rootDir"`
}

type CommandReader struct {
	io.Reader
	io.Writer
	SessionId string
	AttrName  string
	b         bytes.Buffer
	cmd       *CommandFrags
}
