package webscan

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/twmb/murmur3"
)

func (e *Engine) GetHash(url string) (iconhash string) {
	content, err := e.FromURLGetContent(url)
	if err != nil {
		return
	}
	if len(content) > 0 {
		iconhash = Mmh3Hash32(StandBase64(content))
	}
	return
}

func (e *Engine) FromURLGetContent(url string) (content []byte, err error) {
	resp, err := e.reqClient.R().Get(url)
	if err != nil {
		return
	}
	if resp.StatusCode != 200 {
		err = fmt.Errorf("Req %v status is not 200", url)
	}
	content = resp.Bytes()
	return
}

func Mmh3Hash32(raw []byte) string {
	var h32 = murmur3.New32()
	h32.Write(raw)
	return fmt.Sprintf("%d", int32(h32.Sum32()))
}

func StandBase64(braw []byte) []byte {
	bckd := base64.StdEncoding.EncodeToString(braw)
	var buffer bytes.Buffer
	for i := 0; i < len(bckd); i++ {
		ch := bckd[i]
		buffer.WriteByte(ch)
		if (i+1)%76 == 0 {
			buffer.WriteByte('\n')
		}
	}
	buffer.WriteByte('\n')
	return buffer.Bytes()
}
