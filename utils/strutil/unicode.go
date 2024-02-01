package strutil

import (
	"bytes"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io/ioutil"
	"unicode/utf8"
)

func Byte2UTF8(bt []byte) (res string, err error) {

	// utf8 Valid方法判断格式
	if utf8.Valid(bt) {
		return string(bt), nil
	}
	bBt, err := GbkToUtf8(bt)
	if err != nil {
		return "", err
	}
	return string(bBt), err
}

func Str2UTF8(str string) (string, error) {

	bt := []byte(str)
	// utf8 Valid方法判断格式
	if utf8.Valid(bt) {
		return str, nil
	}
	res, err := GbkToUtf8(bt)
	if err != nil {
		return "", err
	}
	return string(res), err
}

// GbkToUtf8 utf8 格式
func GbkToUtf8(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}
