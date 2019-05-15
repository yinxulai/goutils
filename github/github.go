package github

import (
	"encoding/json"
	"fmt"
	"net/url"

	httpd "github.com/yinxulai/grpc-services/toolbox/http"
)

const (
	redirectAuthURL = "https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s&scope=%s&state=%s&allow_signup=%s"
	accessTokenURL  = "https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code"
)

type Github struct {
	ClientID     string
	ClientSecret string
}

// New 实例化授权信息
func New() *Github {
	github := new(Github)
	return github
}

//GenerateAuthURL 获取跳转的url地址
func (oauth *Github) GenerateAuthURL(redirectURI, scope, state string) (string, error) {
	//url encode
	urlStr := url.QueryEscape(redirectURI)
	completeURL := fmt.Sprintf(redirectAuthURL, oauth.ClientID, urlStr, scope, state, true)
	return completeURL, nil
}

type ResAccessToken struct {
	CommonError
	AccessToken string `json:"access_token"`
	TokenType   int64  `json:"token_type"`
	Scope       string `json:"scope"`
}

// GetUserAccessToken 通过网页授权的code 换取access_token(区别于context中的access_token)
func (oauth *Github) GetUserAccessToken(code string) (result *ResAccessToken, err error) {
	urlStr := fmt.Sprintf(accessTokenURL, oauth.ClientID, oauth.ClientSecret, code)
	var response []byte
	response, err = httpd.Get(urlStr)
	if err != nil {
		return
	}
	err = json.Unmarshal(response, &result)
	if err != nil {
		return
	}
	if result.ErrCode != 0 {
		err = fmt.Errorf("GetUserAccessToken error : errcode=%v , errmsg=%v", result.ErrCode, result.ErrMsg)
		return
	}
	return
}
