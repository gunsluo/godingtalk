package dingtalk

import "context"

type ActivateSuiteResponse struct {
	OpenAPIResponse
}

type SuiteAccessTokenResponse struct {
	OpenAPIResponse
	SuiteAccessToken string `json:"suite_access_token"`
	ExpiresIn        int    `json:"expires_in"`
}

type ListUnactivateSuitesResponse struct {
	OpenAPIResponse
	AppId    int      `json:""`
	CorpList []string `json:"corp_list"`
	HasMore  bool     `json:"has_more"`
}

// 刷新并获取第三方企业应用的suite_access_token
func (dtc *DingTalkClient) GetAndRefreshSuiteAccessToken(ctx context.Context) (string, error) {
	dtc.suiteAccessTokenLocker.Lock()
	defer dtc.suiteAccessTokenLocker.Unlock()

	suiteAccessToken, err := dtc.cache.Get(ctx, KeySuiteAccessToken)
	if err == nil {
		return suiteAccessToken.Value, nil
	}

	resp, err := dtc.GetSuiteAccessToken(ctx)
	if err != nil {
		return "", err
	}

	err = dtc.cache.Set(
		ctx,
		NewKVExpirable(KeySuiteAccessToken, resp.SuiteAccessToken, resp.ExpiresIn),
	)
	return resp.SuiteAccessToken, err
}

// 获取第三方企业应用的suite_access_token
func (dtc *DingTalkClient) GetSuiteAccessToken(ctx context.Context) (SuiteAccessTokenResponse, error) {
	requestData := map[string]string{
		"suite_key":    dtc.config.suiteKey,
		"suite_secret": dtc.config.suiteSecret,
		"suite_ticket": dtc.GetSuiteTicket(ctx),
	}

	var data SuiteAccessTokenResponse
	err := dtc.httpIsv(ctx, "service/get_suite_token", nil, requestData, &data)

	return data, err
}

// 激活套件
func (dtc *DingTalkClient) IsvActivateSuite(ctx context.Context, authCorpId string, permanentCode string) (ActivateSuiteResponse, error) {
	var data ActivateSuiteResponse
	requestData := map[string]string{
		"suite_key":      dtc.config.suiteKey,
		"auth_corpid":    authCorpId,
		"permanent_code": permanentCode,
	}
	err := dtc.httpIsv(ctx, "service/activate_suite", nil, requestData, &data)
	return data, err
}

// 获取应用未激活的企业列表
func (dtc *DingTalkClient) IsvListUnactivateSuites(ctx context.Context) (ListUnactivateSuitesResponse, error) {
	var data ListUnactivateSuitesResponse
	requestData := map[string]string{
		"app_id": dtc.config.appId,
	}
	err := dtc.httpIsv(ctx, "service/get_unactive_corp", nil, requestData, &data)
	return data, err
}

// get suite ticket
func (dtc *DingTalkClient) GetSuiteTicket(ctx context.Context) string {
	suiteTicket, err := dtc.persist.Get(ctx, KeySuiteTicket)
	if err != nil {
		return ""
	}

	return suiteTicket.Value
}

// set suite ticket
func (dtc *DingTalkClient) SetSuiteTicket(ctx context.Context, suiteTicket string) error {
	return dtc.persist.Set(
		ctx,
		NewKVExpirable(KeySuiteTicket, suiteTicket, 18300),
	)
}
