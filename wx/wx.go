package wx

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type WXClient struct {
	AppId     string `json:"appId"`
	AppSecret string `json:"appSecret"`
}

func NewWXClient(url, appId, appSecret string) *WXClient {
	return &WXClient{AppId: appId, AppSecret: appSecret}
}

type Jscode2sessionResp struct {
	Openid     string `json:"openid"`     //微信用户的唯一标识
	SessionKey string `json:"sessionKey"` //会话密钥
	Unionid    string `json:"unionid"`    //用户在微信开放平台的唯一标识符。本字段在满足一定条件的情况下才返回。
}

// 通过code 获取微信用户会话信息
func (w *WXClient) Jscode2session(code string) (result Jscode2sessionResp, err error) {
	wxUrl := fmt.Sprintf("https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code", w.AppId, w.AppSecret, code)
	resp, err := http.Get(wxUrl)
	if err != nil {
		return result, err
	}
	if resp.StatusCode != http.StatusOK {
		return result, fmt.Errorf("http request err. resp code=%d", resp.StatusCode)
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return result, err
	}
	err = json.Unmarshal(b, &result)
	if err != nil {
		return result, err
	}
	return result, nil
}
