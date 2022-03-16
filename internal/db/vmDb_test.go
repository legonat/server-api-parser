package db

import (
	"awesomeProjectRucenter/pkg/erx"
	"bytes"
	"encoding/json"
	"os"
	"testing"
)

func BenchmarkPrepareInsertVmsReqConcat(b *testing.B) {
	byt, err := os.ReadFile("/home/legonat/mygo/awesomeProjectRucenter/data.txt")
	if err != nil {
		log.Error(erx.New(err))
	}
	c := new(bytes.Buffer)
	err = json.Compact(c, byt)
	if err != nil {
		log.Error(erx.New(err))
	}
	for n := 0; n < b.N; n++ {
		prepareInsertVmsRequest(c.Bytes())
	}
}
