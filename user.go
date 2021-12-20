package dingtalk

type GetCorpUserByUnionIdResponse struct {
	OpenAPIResponse
	RequestId string                     `json:"request_id"`
	UserInfo  GetCorpUserByUnionIdResult `json:"result"`
}

type GetCorpUserByUnionIdResult struct {
	ContactType int    `json:"contact_type"`
	UserId      string `json:"userid"`
}

type GetCorpUserDetailByUserIdResponse struct {
	OpenAPIResponse
	RequestId string                          `json:"request_id"`
	UserInfo  GetCorpUserDetailByUserIdResult `json:"result"`
}

type GetCorpUserDetailByUserIdResult struct {
	UserId               string       `json:"userid"`
	Unionid              string       `json:"unionid"`
	Name                 string       `json:"name"`
	Avatar               string       `json:"avatar"`
	StateCode            string       `json:"state_code"`
	ManagerUserid        string       `json:"manager_userid"`
	Mobile               string       `json:"mobile"`
	HideMobile           bool         `json:"hide_mobile"`
	Telephone            string       `json:"telephone"`
	Job_number           string       `json:"job_number"`
	Title                string       `json:"title"`
	Email                string       `json:"email"`
	WorkPlace            string       `json:"work_place"`
	Remark               string       `json:"remark"`
	LoginId              string       `json:"login_id"`
	ExclusiveAccountType string       `json:"exclusive_account_type"`
	ExclusiveAccount     bool         `json:"exclusive_account"`
	Extension            string       `json:"extension"`
	HiredDate            int          `json:"hired_date"`
	Active               bool         `json:"active"`
	RealAuthed           bool         `json:"real_authed"`
	OrgEmail             string       `json:"org_email"`
	OrgEmailType         string       `json:"org_email_type"`
	Nickname             string       `json:"nickname"`
	Senior               bool         `json:"senior"`
	Admin                bool         `json:"admin"`
	Boss                 bool         `json:"boss"`
	CorpId               string       `json:"corp_id"`
	DeptIdList           []int        `json:"dept_id_list"`
	DeptOrderList        []DeptOrder  `json:"dept_order_list"`
	LeaderInDept         []DeptLeader `json:"leader_in_dept"`
	RoleList             []UserRole   `json:"role_list"`
	UnionEmpExt          UnionEmpExt  `json:"union_emp_ext"`
}

type DeptOrder struct {
	DeptId int `json:"dept_id"`
	Order  int `json:"order"`
}

type DeptLeader struct {
	DeptId int  `json:"dept_id"`
	Leader bool `json:"leader"`
}

type UserRole struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	GroupName string `json:"group_name"`
}

type UnionEmpExt struct {
	UserId          string        `json:"userid"`
	UnionEmpMapList []UnionEmpMap `json:"union_emp_map_list"`
	CorpId          string        `json:"corp_id"`
}

type UnionEmpMap struct {
	UserId string `json:"userid"`
	CorpId string `json:"corp_id"`
}

type GetCorpUserInfoByCodeResponse struct {
	OpenAPIResponse
	RequestId string                      `json:"request_id"`
	UserInfo  GetCorpUserInfoByCodeResult `json:"result"`
}

type GetCorpUserInfoByCodeResult struct {
	UserId            string `json:"userid"`
	DeviceId          string `json:"device_id"`
	Sys               bool   `json:"sys"`
	SysLevel          int    `json:"sys_level"`
	AssociatedUnionid string `json:"associated_unionid"`
	Unionid           string `json:"unionid"`
	Name              string `json:"name"`
}

// 根据unionid获取组织用户userid
func (dtc *DingTalkClient) GetCorpUserByUnionId(authCorpId, unionId string) (GetCorpUserByUnionIdResponse, error) {
	var data GetCorpUserByUnionIdResponse
	requestData := map[string]string{
		"unionid": unionId,
	}

	err := dtc.httpTOP("user/getbyunionid", authCorpId, nil, requestData, &data)
	return data, err
}

// 通过免登码获取用户信息
func (dtc *DingTalkClient) GetCorpUserInfoByCode(authCorpId, code string) (GetCorpUserInfoByCodeResponse, error) {
	var data GetCorpUserInfoByCodeResponse
	requestData := map[string]string{
		"code": code,
	}

	err := dtc.httpTOP("v2/user/getuserinfo", authCorpId, nil, requestData, &data)
	return data, err
}

// 根据userId获取组织用户详情
func (dtc *DingTalkClient) GetCorpUserDetailByUserId(authCorpId, userId string, languages ...string) (GetCorpUserDetailByUserIdResponse, error) {
	var language string
	if len(languages) > 0 {
		language = languages[0]
	} else {
		language = "zh_CN"
	}
	var data GetCorpUserDetailByUserIdResponse
	requestData := map[string]string{
		"userid":   userId,
		"language": language,
	}

	err := dtc.httpTOP("v2/user/get", authCorpId, nil, requestData, &data)
	return data, err
}
