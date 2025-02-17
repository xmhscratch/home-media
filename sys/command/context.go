package command

import (
	"bytes"
	"home-media/sys"
	"io"
)

type NullWriter struct {
	io.Writer
}

type sessionWriterAbstract struct {
	io.Writer
	SessionId string
	AttrName  string
	Config    *sys.Config
}

type SessionInfoWriter struct {
	*sessionWriterAbstract
}

type SessionFileWriter struct {
	*sessionWriterAbstract
	FileKey string
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
	StreamIndex string `json:"streamIndex"`
	LangCode    string `json:"langCode"`
}

type CommandReader struct {
	io.Reader
	io.Writer
	SessionId string
	AttrName  string
	b         bytes.Buffer
	cmd       *CommandFrags
}
