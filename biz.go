package dingtalk

import "encoding/json"

type Biz2Data struct {
	SyncAction  string `json:"syncAction"`
	SuiteTicket string `json:"suiteTicket"`
	SyncSeq     string `json:"syncSeq"`
}

type Biz4Data struct {
	AuthUserInfo  BizAuthUserInfo `json:"auth_user_info"`
	AuthCorpInfo  BizAuthCorpInfo `json:"auth_corp_info"`
	PermanentCode string          `json:"permanent_code"`
	SyncAction    string          `json:"syncAction"`
	SyncSeq       string          `json:"syncSeq"`
	AuthInfo      BizAuthInfo     `json:"auth_info"`
	AuthScope     BizAuthScope    `json:"auth_scope"`
}

type Biz16Data struct {
	SyncAction string `json:"syncAction"`
	SyncSeq    string `json:"syncSeq"`
}

type BizAuthUserInfo struct {
	UserId string `json:"userId"`
}

type BizAuthCorpInfo struct {
	CorpLogoURL     string `json:"corp_logo_url"`
	CorpName        string `json:"corp_name"`
	CorpID          string `json:"corpid"`
	Industry        string `json:"industry"`
	InviteCode      string `json:"invite_code"`
	LicenseCode     string `json:"license_code"`
	AuthChannel     string `json:"auth_channel"`
	AuthChannelType string `json:"auth_channel_type"`
	IsAuthenticated bool   `json:"is_authenticated"`
	AuthLevel       int    `json:"auth_level"`
	InviteURL       string `json:"invite_url"`
}

type BizAuthInfo struct {
	Agents []BizAgent `json:"agent"`
}

type BizAgent struct {
	AdminList string `json:"admin_list"`
	AgentName string `json:"agent_name"`
	AgentId   string `json:"agentid"`
	AppId     string `json:"appid"`
	LogoUrl   string `json:"logo_url"`
}

type BizAuthScope struct {
	ErrCode         int              `json:""`
	ErrMsg          string           `json:""`
	ConditionFields []string         `json:"condition_field"`
	AuthUserFields  []string         `json:"auth_user_field"`
	AuthOrgScopes   BizAuthOrgScopes `json:"auth_org_scopes"`
}

type BizAuthOrgScopes struct {
	AuthedUser []string `json:"authed_user"`
	AuthedDept []int    `json:"authed_dept"`
}

var bizTypesUnmarshalJSONFuncs = map[int]func([]byte) (interface{}, error){
	2: func(data []byte) (interface{}, error) {
		var d Biz2Data
		if err := json.Unmarshal([]byte(data), &d); err != nil {
			return nil, err
		}
		return d, nil
	},
	4: func(data []byte) (interface{}, error) {
		var d Biz4Data
		if err := json.Unmarshal([]byte(data), &d); err != nil {
			return nil, err
		}
		return d, nil
	},
	16: func(data []byte) (interface{}, error) {
		var d Biz16Data
		if err := json.Unmarshal([]byte(data), &d); err != nil {
			return nil, err
		}
		return d, nil
	},
}
