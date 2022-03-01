package main

import (
	"context"
	"golang.org/x/text/encoding/simplifiedchinese"
	"log"
	"os/exec"
)

type Charset string

const (
	UTF8    = Charset("UTF-8")
	GB18030 = Charset("GB18030")
)

func main() {
	cmd := exec.CommandContext(context.TODO(), "ping", "127.0.0.1", "-n", "2")
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Println(err)
		return
	}
	cmdRe := ConvertByte2String(output, "GB18030")
	log.Println(string(cmdRe))
}

func ConvertByte2String(byte []byte, charset Charset) string {
	var str string
	switch charset {
	case GB18030:
		var decodeBytes, _ = simplifiedchinese.GB18030.NewDecoder().Bytes(byte)
		str = string(decodeBytes)
	case UTF8:
		fallthrough
	default:
		str = string(byte)
	}
	return str
}
