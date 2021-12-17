package dingtalk

type ActivateSuiteResponse struct {
	OpenAPIResponse
}

type SuiteAccessTokenResponse struct {
	OpenAPIResponse
	SuiteAccessToken string `json:"suite_access_token"`
	ExpiresIn        int    `json:"expires_in"`
}

// 刷新并获取第三方企业应用的suite_access_token
func (dtc *DingTalkClient) GetAndRefreshSuiteAccessToken() (string, error) {
	dtc.suiteAccessTokenLocker.Lock()
	defer dtc.suiteAccessTokenLocker.Unlock()

	var suiteAccessToken StringExpirable
	err := dtc.suiteAccessTokenCache.Get(&suiteAccessToken)
	if err == nil {
		return suiteAccessToken.Value, nil
	}

	resp, err := dtc.GetSuiteAccessToken()
	if err != nil {
		return "", err
	}

	err = dtc.suiteAccessTokenCache.Set(
		NewStringExpirable(resp.SuiteAccessToken, resp.ExpiresIn),
	)
	return suiteAccessToken.Value, err
}

// 获取第三方企业应用的suite_access_token
func (dtc *DingTalkClient) GetSuiteAccessToken() (SuiteAccessTokenResponse, error) {
	requestData := map[string]string{
		"suite_key":    dtc.config.suiteKey,
		"suite_secret": dtc.config.suiteSecret,
		"suite_ticket": dtc.GetSuiteTicket(),
	}

	var data SuiteAccessTokenResponse
	err := dtc.httpIsv("service/get_suite_token", nil, requestData, &data)

	return data, err
}

// 激活套件
func (dtc *DingTalkClient) IsvActivateSuite(authCorpId string, permanentCode string) (ActivateSuiteResponse, error) {
	var data ActivateSuiteResponse
	requestData := map[string]string{
		"suite_key":      dtc.config.suiteKey,
		"auth_corpid":    authCorpId,
		"permanent_code": permanentCode,
	}
	err := dtc.httpIsv("service/activate_suite", nil, requestData, &data)
	return data, err
}

// get suite ticket
func (dtc *DingTalkClient) GetSuiteTicket() string {
	var suiteTicket StringExpirable
	err := dtc.suiteTicketCache.Get(&suiteTicket)
	if err != nil {
		return ""
	}

	return suiteTicket.Value
}

// set suite ticket
func (dtc *DingTalkClient) SetSuiteTicket(suiteTicket string) error {
	return dtc.suiteTicketCache.Set(
		NewStringExpirable(suiteTicket, 18000),
	)
}
