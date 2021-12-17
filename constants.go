package dingtalk

const (
	OAPIURL   = "https://oapi.dingtalk.com/"
	TOPAPIURL = "https://oapi.dingtalk.com/topapi/"
)

const (
	MessageTypeText       = "text"
	MessageTypeActionCard = "action_card"
	MessageTypeImage      = "image"
	MessageTypeVoice      = "voice"
	MessageTypeFile       = "file"
	MessageTypeLink       = "link"
	MessageTypeOA         = "oa"
	MessageTypeMarkdown   = "markdown"
)

const (
	signMD5         = "MD5"
	signHMAC        = "HMAC"
	topFormat       = "json"
	topV            = "2.0"
	topSimplify     = false
	topSecret       = "github.com/icepy"
	topSignMethod   = signMD5
	typeJSON        = "application/json"
	typeForm        = "application/x-www-form-urlencoded"
	typeMultipart   = "multipart/form-data"
	aesEncodeKeyLen = 43
)

const (
	CheckCreateSuiteURLEventType = "check_create_suite_url"
	CheckUpdateSuiteUrlEventType = "check_update_suite_url"
	CheckUrlEventType            = "check_url"
	SyncHTTPPushHighEventType    = "SYNC_HTTP_PUSH_HIGH"
	SyncHTTPPushMediumEventType  = "SYNC_HTTP_PUSH_MEDIUM"
)
