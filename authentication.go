package dingtalk

import (
	"net/url"
)

type GetCorpAccessTokenResponse struct {
	OpenAPIResponse
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

type GetAccessTokenResponse struct {
	OpenAPIResponse
	AccessToken string `json:"access_token"`
}

// 刷新并获取企业授权的凭证
func (dtc *DingTalkClient) IsvGetAndRefreshCorpAccessToken(authCorpId string) (string, error) {
	key := "corp_access_token_" + authCorpId
	corpAccessToken, err := dtc.cache.Get(key)
	if err == nil {
		return corpAccessToken.Value, nil
	}

	resp, err := dtc.IsvGetCorpAccessToken(authCorpId)
	if err != nil {
		return "", err
	}

	err = dtc.cache.Set(
		NewKVExpirable(key, resp.AccessToken, resp.ExpiresIn),
	)
	return resp.AccessToken, err
}

// 获取企业授权的凭证
func (dtc *DingTalkClient) IsvGetCorpAccessToken(authCorpId string) (GetCorpAccessTokenResponse, error) {
	var data GetCorpAccessTokenResponse
	requestData := map[string]string{
		"auth_corpid": authCorpId,
	}
	params := url.Values{}
	params.Set("accessKey", dtc.getAccessKey())
	err := dtc.httpIsv("service/get_corp_token", params, requestData, &data)
	return data, err
}

// 获取企业内部应用的access_token
func (dtc *DingTalkClient) GetCorpAccessToken() (GetCorpAccessTokenResponse, error) {
	var data GetCorpAccessTokenResponse
	params := url.Values{}
	params.Set("appkey", dtc.config.appKey)
	params.Set("appsecret", dtc.config.appSecret)
	err := dtc.httpRPC("gettoken", params, nil, &data)
	return data, err
}

// 获取微应用后台免登的access_token
func (dtc *DingTalkClient) GetAccessToken() (GetAccessTokenResponse, error) {
	var data GetAccessTokenResponse
	params := url.Values{}
	params.Set("corpid", dtc.config.corpId)
	params.Set("corpsecret", dtc.config.corpSecret)
	err := dtc.httpRPC("sso/gettoken", params, nil, &data)
	return data, err
}
