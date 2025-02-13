package session

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
)

func (ctx *File[T]) GetFileName(filepath string) string {
	re, err := regexp.Compile(`[\/]{0,1}([\w\W]+)+([\.]{1}[a-zA-Z0-9]+?)$`)
	matches := re.FindStringSubmatch(filepath)
	if err != nil {
		fmt.Println(err)
	}

	if len(matches) == 3 {
		return matches[1]
	} else {
		return "download-file"
	}
}

func (ctx *File[T]) GetFileExt(filepath string) string {
	re, err := regexp.Compile(`[\/]{0,1}([\w\W]+)+([\.]{1}[a-zA-Z0-9]+?)$`)
	matches := re.FindStringSubmatch(filepath)
	if err != nil {
		fmt.Println(err)
	}

	// mime.TypeByExtension()
	if len(matches) == 3 {
		return matches[2]
	} else {
		return ".bin"
	}
}

func ParseSourceType(s string) (FileSourceType, error) {
	switch s {
	case FILE_SOURCE_TYPE_DIRECT.String():
		return FILE_SOURCE_TYPE_DIRECT, nil
	case FILE_SOURCE_TYPE_TORRENT.String():
		return FILE_SOURCE_TYPE_TORRENT, nil
	default:
		return 0, errors.New("invalid type")
	}
}

func (ctx *FileMetaInfoList) ToArray() []interface{} {
	var rs []interface{} = []interface{}{}
	for k, v := range *ctx {
		rs = append(rs, k)
		rs = append(rs, v)
	}
	return rs
}

func (ctx *FileMetaInfoList) ToMap() map[string]FileMetaInfo {
	var rs map[string]FileMetaInfo = map[string]FileMetaInfo{}
	for k, v := range *ctx {
		rs[k] = v
	}
	return rs
}

func (ctx *FileMetaInfoList) GetValue(key string) FileMetaInfo {
	var rs FileMetaInfo
	for k, v := range *ctx {
		if k != key {
			continue
		}
		rs = v
		break
	}
	return rs
}

type fileMetaInfoAlias struct {
	Path string `json:"path"`
	Size int64  `json:"size"`
}

func (ctx FileMetaInfo) MarshalJSON() ([]byte, error) {
	return json.Marshal(&fileMetaInfoAlias{
		Path: ctx.Path,
		Size: ctx.Size,
	})
}

func (ctx *FileMetaInfo) UnmarshalJSON(data []byte) error {
	var err error
	_, err = fileMetaInfoUnmarshal(ctx, data)
	return err
}

func (ctx FileMetaInfo) MarshalBinary() ([]byte, error) {
	return ctx.MarshalJSON()
}

func (ctx *FileMetaInfo) UnmarshalBinary(data []byte) error {
	var err error
	_, err = fileMetaInfoUnmarshal(ctx, data)
	return err
}

func (list FileMetaInfoList) MarshalBinary() ([]byte, error) {
	return json.Marshal(list)
}

func (list *FileMetaInfoList) UnmarshalBinary(data []byte) error {
	return json.Unmarshal([]byte(data), &list)
}

func fileMetaInfoUnmarshal(ctx *FileMetaInfo, data []byte) (*FileMetaInfo, error) {
	var (
		err       error
		unescaped string
		a         *fileMetaInfoAlias
	)

doneParsing:
	for range [1]struct{}{} {
		if err = json.Unmarshal([]byte(data), &a); err != nil {
			if err = json.Unmarshal([]byte(data), &unescaped); err != nil {
				break doneParsing
			}
			if err = json.Unmarshal([]byte(unescaped), &a); err != nil {
				break doneParsing
			}
		}
	}

	if a != nil {
		err = nil
	}

	if ctx == nil {
		ctx = &FileMetaInfo{
			Path: a.Path,
			Size: a.Size,
		}
	} else {
		ctx.Path = a.Path
		ctx.Size = a.Size
	}

	return ctx, err
}
