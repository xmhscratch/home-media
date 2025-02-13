package session

import (
	"encoding/json"
	"testing"

	"github.com/sanity-io/litter"
)

// var testFSubDubJSON string = "[{\"stream_index\":1,\"codec_name\":\"eac3\",\"lang_code\":\"jpn\",\"lang_title\":\"Japanese (E-AC-3) (2.0) [AMZN]\"},{\"stream_index\":2,\"codec_name\":\"eac3\",\"lang_code\":\"eng\",\"lang_title\":\"English (E-AC-3) (2.0) [AMZN]\"}]"

var testFSubDubJSON string = "[{\"stream_index\":3,\"codec_name\":\"ass\",\"lang_code\":\"eng\",\"lang_title\":\"English (Full) (ASS) [Chihiro] (restyled)\"},{\"stream_index\":4,\"codec_name\":\"ass\",\"lang_code\":\"enm\",\"lang_title\":\"English (Full with Honorifics) (ASS) [CR]+[Chihiro]\"},{\"stream_index\":5,\"codec_name\":\"ass\",\"lang_code\":\"ara\",\"lang_title\":\"Arabic - Saudi Arabia (Full) (ASS) [CR]\"},{\"stream_index\":6,\"codec_name\":\"ass\",\"lang_code\":\"chi\",\"lang_title\":\"Chinese - China (Full) (ASS) [CR]\"},{\"stream_index\":7,\"codec_name\":\"ass\",\"lang_code\":\"chi\",\"lang_title\":\"Chinese - Hong Kong (Full) (ASS) [CR]\"},{\"stream_index\":8,\"codec_name\":\"ass\",\"lang_code\":\"fre\",\"lang_title\":\"French (Full) (ASS) [CR]\"},{\"stream_index\":9,\"codec_name\":\"ass\",\"lang_code\":\"ger\",\"lang_title\":\"German (Full) (ASS) [CR]\"},{\"stream_index\":10,\"codec_name\":\"ass\",\"lang_code\":\"ind\",\"lang_title\":\"Indonesian (Full) (ASS) [CR]\"},{\"stream_index\":11,\"codec_name\":\"ass\" ,\"lang_code\":\"ita\",\"lang_title\":\"Italian (Full) (ASS) [CR]\"},{\"stream_index\":12,\"codec_name\":\"ass\",\"lang_code\":\"may\",\"lang_title\":\"Malay (Full) (ASS) [CR]\"},{\"stream_index\":13,\"codec_name\":\"ass\",\"lang_code\":\"por\",\"lang_title\":\"Portuguese - Brazil (Full) (ASS) [CR]\"},{\"stream_index\":14,\"codec_name\":\"ass\",\"lang_code\":\"rus\",\"lang_title\":\"Russian (Full) (ASS) [CR]\"},{\"stream_index\":15,\"codec_name\":\"ass\",\"lang_code\":\"spa\",\"lang_title\":\"Spanish - Europe (Full) (ASS) [CR]\"},{\"stream_index\":16,\"codec_name\":\"ass\",\"lang_code\":\"spa\",\"lang_title\":\"Spanish - Latin America (Full) (ASS) [CR]\"},{\"stream_index\":17,\"codec_name\":\"ass\",\"lang_code\":\"tha\",\"lang_title\":\"Thai (Full) (ASS) [CR]\"},{\"stream_index\":18,\"codec_name\":\"ass\",\"lang_code\":\"vie\",\"lang_title\":\"Vietnamese (Full) (ASS) [CR]\"},{\"stream_index\":19,\"codec_name\":\"ass\",\"lang_code\":\"eng\",\"lang_title\":\"English (Signs & Songs) (ASS) [Chihiro]\"}]"

func TestFSubDubMarshalJSON(t *testing.T) {
	var (
		err         error
		b           []byte
		rawFSubDubs FSubDubList = FSubDubList{
			FSubDub{
				StreamIndex: 1,
				CodecName:   "eac3",
				LangCode:    "jpn",
				LangTitle:   "Japanese (E-AC-3) (2.0) [AMZN]",
			},
		}
	)
	// json.Marshal(rawFSubDubs)
	b, err = json.Marshal(rawFSubDubs)
	litter.D(err, string(b))
}

func TestFSubDubUnmarshalJSON(t *testing.T) {
	var (
		err         error
		rawFSubDubs FSubDubList
	)
	// json.Unmarshal([]byte(testFSubDubJSON), &rawFSubDubs)
	err = json.Unmarshal([]byte(testFSubDubJSON), &rawFSubDubs)
	litter.D(err, rawFSubDubs)
}
