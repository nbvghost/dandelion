package wechat

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

/**
curl 'https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=3aa1af8a-8f1e-4cd9-b23e-b76c2265cb83' -H 'Content-Type: application/json' \
-d '{
        "msgtype": "text",
        "text": {
            "content": "流水线已完成，'${CI_SOURCE_NAME}'，'${PIPELINE_NAME}'"
        }
   }'

*/

var defaultWorkGroupBot = &WorkGroupBot{Key: "3aa1af8a-8f1e-4cd9-b23e-b76c2265cb83", Client: &http.Client{}}

type WorkGroupBot struct {
	Key    string
	Client *http.Client
}

type Result struct {
	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`
}

func (m *WorkGroupBot) baseRequest(data any) error {
	bodyByte, err := json.Marshal(data)
	if err != nil {
		return err
	}
	request, err := http.NewRequest("POST", fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=%s", m.Key), bytes.NewBuffer(bodyByte))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json")

	response, err := m.Client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	responseBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}
	result := &Result{}
	err = json.Unmarshal(responseBytes, result)
	if err != nil {
		return err
	}
	if result.Errcode > 0 {
		return errors.New(result.Errmsg)
	}
	log.Println(string(responseBytes))
	return nil
}
func (m *WorkGroupBot) SendText(text string) error {
	var body = struct {
		MsgType string         `json:"msgtype"`
		Text    map[string]any `json:"text"`
	}{
		MsgType: "text",
		Text:    map[string]any{"content": text+fmt.Sprintf("\n\n通知时间：%s\n",time.Now().Format(time.DateTime))},
	}

	return m.baseRequest(&body)
}
func (m *WorkGroupBot) SendMarkdown(text string) error {
	var body = struct {
		MsgType  string         `json:"msgtype"`
		Markdown map[string]any `json:"markdown"`
	}{
		MsgType:  "markdown",
		Markdown: map[string]any{"content": text},
	}

	return m.baseRequest(&body)
}

func SendText(text string) error {
	go func() {
		err := defaultWorkGroupBot.SendText(text)
		if err != nil {
			log.Println(err)
		}
	}()
	return nil
}

func SendMarkdown(text string) error {
	go func() {
		err := defaultWorkGroupBot.SendMarkdown(text)
		if err != nil {
			log.Println(err)
		}
	}()
	return nil
}
