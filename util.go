package gerry

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	// "unicode/utf8"

	sj "github.com/bitly/go-simplejson"
	rq "github.com/parnurzeal/gorequest"
	qr "github.com/skip2/go-qrcode"
)

const ERROR_DEFAULT_RESPONSE string = `{"result": {"state": "999", "message": "服务器发生错误！"}}`
const GENERAL_RESPONSE_TEMPLATE string = `{"result": {"state": "%s", "message": "%s"}}`

func GetCurrentPath() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return ""
	}
	return dir
}

func GetCurrentPaths() []string {
	path := GetCurrentPath()

	if strings.Contains(path, "/") {
		return strings.Split(path, "/")
	}

	return strings.Split(path, "\\")
}

func DoPost(url string, data interface{}, partialNode string) string {
	content := ""

	request := rq.New()
	if resp, body, _ := request.Post(url).Send(data).End(); resp.StatusCode != 200 {
		Logger.Error("URL:", url, " STATUS:", resp.StatusCode, " Data:", data)
		content = ERROR_DEFAULT_RESPONSE
	} else if partialNode != "" {
		raw, err := sj.NewJson([]byte(body))
		// fmt.Printf("\nRAW JSON: %s, ERROR: %t", raw, (err != nil))
		partial, err := raw.Get(partialNode).MarshalJSON()
		// fmt.Printf("\nPARTIAL JSON: %s, ERROR: %t", string(partial[:]), (err != nil))
		if err != nil {
			fmt.Printf("\nerror found: %s", err.Error())
		}
		content = string(partial[:])
	} else {
		content = body
	}

	Logger.Info(content)

	return content
}

func DoPostHTMLSafe(url string, data interface{}) string {
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	enc.Encode(data)
	payload := strings.NewReader(buf.String())

	req, _ := http.NewRequest("POST", url, payload)
	req.Header.Add("content-type", "application/json")
	req.Header.Add("cache-control", "no-cache")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	return string(body)
}

func DoPostRaw(url string, data string, partialNode string) string {
	payload := strings.NewReader(fmt.Sprintf(`"%s"`, data))
	// payload := strings.NewReader(data)
	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("content-type", "application/json")
	req.Header.Add("cache-control", "no-cache")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	content := string(body[:])

	Logger.Info(content)
	return content
}

func CheckRequestData(data map[string]string) (key string, ok bool) {
	ok = true
	for k, v := range data {
		if string(v) == "" || string(v) == "0" {
			ok = false
			key = k
			break
		}
	}

	return
}

func MakeErrorResponse(message string) string {
	return fmt.Sprintf(GENERAL_RESPONSE_TEMPLATE, "999", message)
}

func MakeSuccessResponse(message string) string {
	return fmt.Sprintf(GENERAL_RESPONSE_TEMPLATE, "0", message)
}

func DoGet(url string, data interface{}) string {
	return ""
}

func MakeQRCode(content string) string {
	prefix := "data:image/png;base64,"
	if png, err := qr.Encode(content, qr.Medium, 256); err != nil {
		return prefix
	} else {
		final := base64.StdEncoding.EncodeToString(png)
		return prefix + final
	}
}

func MakeSimpleString(content string, size int) string {
	uniContent := []rune(content)
	if len(uniContent) <= size {
		return content
	}

	newContent := string(uniContent[:size-3]) + "……"

	return newContent
}
