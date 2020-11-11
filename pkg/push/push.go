package push

import (
	"bytes"
	"encoding/json"
	"net/http"
)

const (
	PROXY_FEISHU_GROUP_URL = "https://open.feishu.cn/open-apis/bot/hook/35d3307db9b543058c3e483992a2921c"
)

type feishu_message struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

func SendProxyPush(title, text string) error {
	var message = feishu_message{
		Title: title,
		Text:  text,
	}
	b, err := json.Marshal(message)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", PROXY_FEISHU_GROUP_URL, bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}
