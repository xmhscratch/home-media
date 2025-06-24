package session

import (
	"encoding/json"
	"strconv"
)

type fStreamInfoAlias struct {
	StreamIndex    int64       `json:"stream_index"`
	CodecName      string      `json:"codec_name"`
	LangCode       string      `json:"lang_code"`
	LangTitle      string      `json:"lang_title"`
	BitRate        interface{} `json:"bps,omitempty"`
	Duration       string      `json:"duration,omitempty"`
	NumberOfFrames interface{} `json:"number_of_frames,omitempty"`
	NumberOfBytes  interface{} `json:"number_of_bytes,omitempty"`
}

func (ctx FStreamInfo) MarshalJSON() ([]byte, error) {
	return json.Marshal(&fStreamInfoAlias{
		StreamIndex:    ctx.StreamIndex,
		CodecName:      ctx.CodecName,
		LangCode:       ctx.LangCode,
		LangTitle:      ctx.LangTitle,
		BitRate:        ctx.BitRate,
		Duration:       ctx.Duration,
		NumberOfFrames: ctx.NumberOfFrames,
		NumberOfBytes:  ctx.NumberOfBytes,
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

	var (
		bps            int64
		numberOfFrames int64
		numberOfBytes  int64
	)
	switch a.BitRate.(type) {
	case string:
		if bps, err = strconv.ParseInt(a.BitRate.(string), 10, 64); err != nil {
			break
		}
	}
	switch a.NumberOfFrames.(type) {
	case string:
		if numberOfFrames, err = strconv.ParseInt(a.NumberOfFrames.(string), 10, 64); err != nil {
			break
		}
	}
	switch a.NumberOfBytes.(type) {
	case string:
		if numberOfBytes, err = strconv.ParseInt(a.NumberOfBytes.(string), 10, 64); err != nil {
			break
		}
	}

	if err != nil {
		return err
	}

	if ctx == nil {
		ctx = &FStreamInfo{
			StreamIndex:    a.StreamIndex,
			CodecName:      a.CodecName,
			LangCode:       a.LangCode,
			LangTitle:      a.LangTitle,
			BitRate:        bps,
			Duration:       a.Duration,
			NumberOfFrames: numberOfFrames,
			NumberOfBytes:  numberOfBytes,
		}
	} else {
		ctx.StreamIndex = a.StreamIndex
		ctx.CodecName = a.CodecName
		ctx.LangCode = a.LangCode
		ctx.LangTitle = a.LangTitle
		ctx.BitRate = bps
		ctx.Duration = a.Duration
		ctx.NumberOfFrames = numberOfFrames
		ctx.NumberOfBytes = numberOfBytes
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
		var (
			bps            int64
			numberOfFrames int64
			numberOfBytes  int64
		)
		switch f.BitRate.(type) {
		case string:
			if bps, err = strconv.ParseInt(f.BitRate.(string), 10, 64); err != nil {
				break
			}
		}
		switch f.NumberOfFrames.(type) {
		case string:
			if numberOfFrames, err = strconv.ParseInt(f.NumberOfFrames.(string), 10, 64); err != nil {
				break
			}
		}
		switch f.NumberOfBytes.(type) {
		case string:
			if numberOfBytes, err = strconv.ParseInt(f.NumberOfBytes.(string), 10, 64); err != nil {
				break
			}
		}

		if err != nil {
			return err
		}
		fsInf := FStreamInfo{
			StreamIndex:    f.StreamIndex,
			CodecName:      f.CodecName,
			LangCode:       f.LangCode,
			LangTitle:      f.LangTitle,
			BitRate:        bps,
			Duration:       f.Duration,
			NumberOfFrames: numberOfFrames,
			NumberOfBytes:  numberOfBytes,
		}
		ls = append(ls, fsInf)
	}
	*list = ls
	return err
}
