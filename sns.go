package dingtalk

import "net/url"

type SNSGetUserInfoByCodeResponse struct {
	OpenAPIResponse
	UserInfo SNSGetUserInfoByCode `json:"user_info"`
}

type SNSGetUserInfoByCode struct {
	Nick          string `json:"nick"`
	OpenId        string `json:"openid"`
	UnionId       string `json:"unionid"`
	AuthHighLevel bool   `json:"main_org_auth_high_level"`
}

func (dtc *DingTalkClient) SNSGetUserInfoByCode(tmpAuthCode string) (SNSGetUserInfoByCodeResponse, error) {
	var data SNSGetUserInfoByCodeResponse
	params := url.Values{}
	params.Add("accessKey", dtc.getAccessKey())
	requestData := map[string]string{
		"tmp_auth_code": tmpAuthCode,
	}
	err := dtc.httpSNS("sns/getuserinfo_bycode", params, requestData, &data)
	return data, err
}
