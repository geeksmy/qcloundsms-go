package sms

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/geeksmy/qcloundsms-go/util"
	"io/ioutil"
	"net/http"
	"time"
)

const sendFVoiceUrl = "https://cloud.tim.qq.com/v5/tlsvoicesvr/sendfvoice?sdkappid=%d&random=%d"

type FileVoiceSender struct {
	Base
}

func NewFileVoiceSender(appID int, appKey string) *FileVoiceSender {
	return &FileVoiceSender{
		Base{
			AppID:  appID,
			AppKey: appKey,
		},
	}
}

/*
 * 发送文件语音
 * nationCode 国家码，如 86 为中国
 * phoneNumber 不带国家码的手机号
 * fid 语音文件fid
 * playTimes 播放次数
 * ext 服务端原样返回的参数，可填空
 */
func (f *FileVoiceSender) Send(nationCode, phoneNumber, fid string, playTimes int, ext string) (*FileVoiceSenderResult, error) {
	random := util.GetRandom()
	now := util.GetCurrentTime()

	type Tel struct {
		NationCode string `json:"nationcode"`
		Mobile     string `json:"mobile"`
	}

	type Body struct {
		Tel       *Tel   `json:"tel"`
		Fid       string `json:"fid"`
		PlayTimes int    `json:"playtimes"`
		Sig       string `json:"sig"`
		Time      int64  `json:"time"`
		Ext       string `json:"ext,omitempty"`
	}

	body := new(Body)
	body.Tel = &Tel{
		NationCode: nationCode,
		Mobile:     phoneNumber,
	}
	body.Fid = fid
	body.PlayTimes = playTimes
	body.Sig = util.CalculateSignatureWithPhoneNumber(f.AppKey, random, now, phoneNumber)
	body.Time = now
	body.Ext = ext

	bodyJson, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	httpClient := &http.Client{
		Timeout: 60 * time.Second,
	}

	req, err := http.NewRequest("POST", fmt.Sprintf(sendFVoiceUrl, f.AppID, random), bytes.NewBuffer(bodyJson))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	result := new(FileVoiceSenderResult)
	err = result.ParseFromHTTPResponseBody(b)
	if err != nil {
		return nil, err
	}

	return result, nil
}
