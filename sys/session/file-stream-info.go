package session

import (
	"encoding/json"
)

type fStreamInfoAlias struct {
	StreamIndex int64  `json:"stream_index"`
	CodecName   string `json:"codec_name"`
	LangCode    string `json:"lang_code"`
	LangTitle   string `json:"lang_title"`
}

func (ctx FStreamInfo) MarshalJSON() ([]byte, error) {
	return json.Marshal(&fStreamInfoAlias{
		StreamIndex: ctx.StreamIndex,
		CodecName:   ctx.CodecName,
		LangCode:    ctx.LangCode,
		LangTitle:   ctx.LangTitle,
	})
}

func (ctx *FStreamInfo) UnmarshalJSON(data []byte) error {
	var (
		err       error
		unescaped string
		a         *fStreamInfoAlias
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
		ctx = &FStreamInfo{
			StreamIndex: a.StreamIndex,
			CodecName:   a.CodecName,
			LangCode:    a.LangCode,
			LangTitle:   a.LangTitle,
		}
	} else {
		ctx.StreamIndex = a.StreamIndex
		ctx.CodecName = a.CodecName
		ctx.LangCode = a.LangCode
		ctx.LangTitle = a.LangTitle
	}

	return err
}

func (list FStreamInfoList) MarshalJSON() ([]byte, error) {
	var rs []FStreamInfo = []FStreamInfo{}
	for _, f := range list {
		rs = append(rs, f)
	}
	return json.Marshal(rs)
}

func (list *FStreamInfoList) UnmarshalJSON(data []byte) error {
	var (
		err       error
		unescaped string
		rs        []fStreamInfoAlias = []fStreamInfoAlias{}
	)

doneParsing:
	for range [1]struct{}{} {
		if err = json.Unmarshal([]byte(data), &rs); err != nil {
			if err = json.Unmarshal([]byte(data), &unescaped); err != nil {
				break doneParsing
			}
			if err = json.Unmarshal([]byte(unescaped), &rs); err != nil {
				break doneParsing
			}
		}
	}

	if rs != nil {
		err = nil
	}

	var ls []FStreamInfo = []FStreamInfo{}
	for _, f := range rs {
		ls = append(ls, FStreamInfo(f))
	}
	*list = ls
	return err
}
