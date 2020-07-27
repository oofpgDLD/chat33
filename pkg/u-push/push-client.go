package u_push

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/33cn/chat33/pkg/http"
)

const (
	// The user agent
	USER_AGENT = "Mozilla/5.0"
	// The host
	host = "http://msg.umeng.com"
	// The upload path
	//uploadPath = "/upload"
	// The post path
	postPath = "/api/send"

	// The upload path
	uploadPath = "/upload"
)

type PushClient struct {
}

func (t *PushClient) Send(msg fieldSetter) error {
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	msg.SetPredefinedKeyValue("timestamp", timestamp)
	url := host + postPath
	postBody, err := msg.getPostBody()
	if err != nil {
		return err
	}

	secretStr := "POST" + url + string(postBody) + msg.getAppMasterSecret()
	w := md5.New()
	_, err = io.WriteString(w, secretStr)
	if err != nil {
		return err
	}
	sign := strings.ToLower(hex.EncodeToString(w.Sum(nil)))
	url = url + "?sign=" + sign

	body := bytes.NewBuffer(postBody)
	headers := make(map[string]string)
	headers["User-Agent"] = USER_AGENT
	byte, err := http.HTTPPostJSON(url, headers, body)
	if err != nil {
		var resp map[string]interface{}
		err = json.Unmarshal(byte, &resp)
		if err != nil {
			return err
		}

		if ret := resp["ret"]; ret == "FAIL" {
			if data, ok := resp["data"]; ok {
				if d, ok := data.(map[string]interface{}); ok {
					return errors.New("Failed to send the notification!:" + d["error_code"].(string) + d["error_msg"].(string))
				}
			}
		}
		return errors.New("Failed to send the notification!")
	}
	return nil
}

func (t *PushClient) UploadContents(appkey, appMasterSecret, contents string) (string, error) {
	// Construct the json string
	uploadJson := make(map[string]interface{})
	uploadJson["appkey"] = appkey
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	uploadJson["timestamp"] = timestamp
	uploadJson["content"] = contents
	// Construct the request
	url := host + uploadPath
	postBody, err := json.Marshal(uploadJson)
	if err != nil {
		return "", err
	}

	secretStr := "POST" + url + string(postBody) + appMasterSecret
	w := md5.New()
	_, err = io.WriteString(w, secretStr)
	if err != nil {
		return "", err
	}
	sign := strings.ToLower(hex.EncodeToString(w.Sum(nil)))
	url = url + "?sign=" + sign

	body := bytes.NewBuffer(postBody)
	headers := make(map[string]string)
	headers["User-Agent"] = USER_AGENT
	byte, err := http.HTTPPostJSON(url, headers, body)

	var resp map[string]interface{}
	err = json.Unmarshal(byte, &resp)
	if err != nil {
		return "", err
	}

	if ret := resp["ret"]; ret == "SUCCESS" {
		if data, ok := resp["data"]; ok {
			if d, ok := data.(map[string]interface{}); ok {
				if fileId, ok := d["file_id"]; ok {
					return fmt.Sprintf("%v", fileId), nil
				}
			}
		}
		return "", errors.New("Failed to upload file!: can not get file_id")
	} else {
		if data, ok := resp["data"]; ok {
			if d, ok := data.(map[string]interface{}); ok {
				return "", errors.New("Failed to upload file!:" + d["error_code"].(string) + d["error_msg"].(string))
			}
		}
		return "", errors.New("Failed to upload file!: unknow error")
	}
}
