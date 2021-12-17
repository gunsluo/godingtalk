package dingtalk

import (
	"net/http"
	"sync"
	"time"
)

type ClientType int32

const (
	CORP   ClientType = 0
	ISV    ClientType = 1
	PERSON ClientType = 2
)

type CorpConfig struct {
	CorpId    string
	ApiToken  string
	AgentId   string
	AppKey    string
	AppSecret string
}

type ISVConfig struct {
	CorpId      string
	ApiToken    string
	MiniAppId   string
	AppId       string
	SuiteId     string
	SuiteKey    string
	SuiteSecret string
	AESKey      string
	Token       string
}

type PersonConfig struct {
	CorpId    string
	ApiToken  string
	MiniAppId string
	AppId     string
	AppSecret string
}

type config struct {
	corpId   string
	apiToken string

	agentId string
	appKey  string

	suiteId     string
	suiteKey    string
	suiteSecret string

	miniAppId string
	appId     string
	appSecret string

	token  string
	aesKey string
}

type DingTalkClient struct {
	config         config
	httpClient     *http.Client
	pushCryptoSuit *PushCryptoSuit

	/*
		AccessToken           string
		SSOAccessToken        string
		SNSAccessToken        string
		SuiteAccessToken      string
		AccessTokenCache      Cache
		TicketCache           Cache
		SSOAccessTokenCache   Cache
		SNSAccessTokenCache   Cache
	*/
	clientType ClientType

	suiteAccessTokenCache  Cache
	suiteAccessTokenLocker sync.Mutex

	suiteTicketCache Cache
}

func newDingTalkClient(clientType ClientType, cfg config) *DingTalkClient {
	c := &DingTalkClient{
		config: cfg,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		clientType:            clientType,
		suiteAccessTokenCache: NewFileCache(".suite_access_token"),
		suiteTicketCache:      NewFileCache(".suite_ticket"),
	}

	if cfg.aesKey != "" && cfg.token != "" {
		suit, err := NewPushCryptoSuit(cfg.token, cfg.aesKey, cfg.suiteKey)
		if err != nil {
			panic(err)
		}
		c.pushCryptoSuit = suit
	}

	return c
}

func NewISVClient(cfg ISVConfig) *DingTalkClient {
	return newDingTalkClient(ISV, config{
		corpId:      cfg.CorpId,
		apiToken:    cfg.ApiToken,
		token:       cfg.Token,
		aesKey:      cfg.AESKey,
		miniAppId:   cfg.MiniAppId,
		appId:       cfg.AppId,
		suiteId:     cfg.SuiteId,
		suiteKey:    cfg.SuiteKey,
		suiteSecret: cfg.SuiteSecret,
	})
}

func NewCorpClient(cfg CorpConfig) *DingTalkClient {
	return newDingTalkClient(CORP, config{
		corpId:    cfg.CorpId,
		apiToken:  cfg.ApiToken,
		agentId:   cfg.AgentId,
		appKey:    cfg.AppKey,
		appSecret: cfg.AppSecret,
	})
}

func NewPersonClient(cfg PersonConfig) *DingTalkClient {
	return newDingTalkClient(PERSON, config{
		corpId:    cfg.CorpId,
		apiToken:  cfg.ApiToken,
		miniAppId: cfg.MiniAppId,
		appId:     cfg.AppId,
		appSecret: cfg.AppSecret,
	})
}
