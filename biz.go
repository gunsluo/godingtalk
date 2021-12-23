package dingtalk

import (
	"encoding/json"
)

const (
	SyncActionOrgMicroAppRestore     = "org_micro_app_restore"
	SyncActionOrgMicroAppStop        = "org_micro_app_stop"
	SyncActionOrgMicroAppRemove      = "org_micro_app_remove"
	SyncActionOrgMicroAppScopeUpdate = "org_micro_app_scope_update"
	SyncActionOrgUpdate              = "org_update"
	SyncActionOrgRemove              = "org_remove"
)

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

type Biz7Data struct {
	SyncAction string `json:"syncAction"`
	AgentId    int    `json:"agentId"`
}

type Biz16Data struct {
	SyncAction      string `json:"syncAction"`
	Errcode         int    `json:"errcode"`
	Errmsg          string `json:"errmsg"`
	Corpid          string `json:"corpid"`
	CorpName        string `json:"corp_name"`
	AuthLevel       int    `json:"auth_level"`
	Industry        string `json:"industry"`
	IsAuthenticated bool   `json:"is_authenticated"`
	CorpLogoUrl     string `json:"corp_logo_url"`
}

type BizAuthUserInfo struct {
	UserId string `json:"userId"`
}

type BizAuthCorpInfo struct {
	CorpLogoURL     string `json:"corp_logo_url"`
	CorpType        int    `json:"corp_type"`
	CorpName        string `json:"corp_name"`
	FullCorpName    string `json:"full_corp_name"`
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
	AdminList []string `json:"admin_list"`
	AgentName string   `json:"agent_name"`
	AgentId   int      `json:"agentid"`
	AppId     int      `json:"appid"`
	LogoUrl   string   `json:"logo_url"`
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
	7: func(data []byte) (interface{}, error) {
		var d Biz7Data
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
