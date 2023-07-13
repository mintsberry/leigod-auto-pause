package api

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

var accessToken = "uNvAjLs1gKhPWB2WVPgsixGXSFj2oGgrEa0XFXOKfiOUzLrZCw1Ey0Ug9ZmTzCV8"

func Login(username, password string) error {
	url := "https://webapi.leigod.com/api/auth/login"
	method := "POST"

	hash := md5.Sum([]byte(password))
	md5Password := hex.EncodeToString(hash[:])

	payload := fmt.Sprintf(`{
		"username": "%s",
		"password": "%s",
		"user_type": "0",
		"src_channel": "guanwang",
		"country_code": "86",
		"lang": "zh_CN",
		"region_code": "1"
	}`, username, md5Password)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, strings.NewReader(payload))

	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return err
	}

	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return err
	}
	code := data["code"].(float64)
	if code != 0 {
		return errors.New(data["msg"].(string))
	} else {
		accessToken = data["data"].(map[string]interface{})["login_info"].(map[string]interface{})["account_token"].(string)
	}
	return nil
}

func Pause() {
	url := "https://webapi.leigod.com/api/user/pause"
	payload := fmt.Sprintf(`{
		"account_token": "%s",
		"lang": "zh_CN"
	}`, accessToken)

	req, err := http.NewRequest("POST", url, strings.NewReader(payload))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	req.Header.Set("Content-Type", "text/plain")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(decodeUnicode(body))
}

func decodeUnicode(b []byte) string {
	re := regexp.MustCompile(`\\u([0-9a-fA-F]{4})`)
	escapedStr := re.ReplaceAllStringFunc(string(b), func(m string) string {
		codePoint, _ := strconv.ParseInt(m[2:], 16, 32)
		return string(codePoint)
	})

	return escapedStr
}
