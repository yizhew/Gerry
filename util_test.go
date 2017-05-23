package gerry

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"unicode/utf8"
)

func TestMakeQRcode(t *testing.T) {
	str := "www.baidu.com"

	qrcontent := MakeQRCode(str)

	fmt.Println("BASE64:", qrcontent)
	assert.NotEmpty(t, qrcontent, "should create good base64 string")
}

func TestSimpleString(t *testing.T) {
	fullContent := "负责制定技改工程的预算及可行性报告，并对技改工程中的工程用料、工程进度、工程质量进行监督考核。责实施设备及管线的安装工程的对外承包。"

	uniFullContent := []rune(fullContent)
	fmt.Println("RUNE LEN:", len(uniFullContent))
	fmt.Println("UTF8 LEN:", utf8.RuneCountInString(fullContent))
	fmt.Println("LEN:", len(fullContent))

	simple := MakeSimpleString(fullContent+"hahahah", 40)
	fmt.Println(simple)
	assert.NotEmpty(t, simple, "succss to simplify string")
}
