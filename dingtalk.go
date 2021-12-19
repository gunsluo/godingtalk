package dingtalk

import (
	"fmt"
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
	CorpId     string
	CorpSecret string
	ApiToken   string
	AgentId    string
	AppKey     string
	AppSecret  string
}

type ISVConfig struct {
	CorpId      string
	CorpSecret  string
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
	CorpId     string
	CorpSecret string
	ApiToken   string
	MiniAppId  string
	AppId      string
	AppSecret  string
}

type config struct {
	corpId     string
	corpSecret string // SSOSecret
	apiToken   string

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
	clientType     ClientType

	cache                  Cache
	persist                Persist
	suiteAccessTokenLocker sync.Mutex
}

type Option func(*options)

func WithCache(cache Cache) Option {
	return func(o *options) {
		o.cache = cache
	}
}

func WithPersist(persist Persist) Option {
	return func(o *options) {
		o.persist = persist
	}
}

func newDingTalkClient(clientType ClientType, cfg config, opts ...Option) *DingTalkClient {
	o := defaultOptions(clientType)
	for _, opt := range opts {
		opt(o)
	}

	c := &DingTalkClient{
		config: cfg,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		clientType: clientType,
		cache:      o.cache,
		persist:    o.persist,
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

func NewISVClient(cfg ISVConfig, opts ...Option) *DingTalkClient {
	return newDingTalkClient(ISV, config{
		corpId:      cfg.CorpId,
		corpSecret:  cfg.CorpSecret,
		apiToken:    cfg.ApiToken,
		token:       cfg.Token,
		aesKey:      cfg.AESKey,
		miniAppId:   cfg.MiniAppId,
		appId:       cfg.AppId,
		suiteId:     cfg.SuiteId,
		suiteKey:    cfg.SuiteKey,
		suiteSecret: cfg.SuiteSecret,
	}, opts...)
}

func NewCorpClient(cfg CorpConfig, opts ...Option) *DingTalkClient {
	return newDingTalkClient(CORP, config{
		corpId:     cfg.CorpId,
		corpSecret: cfg.CorpSecret,
		apiToken:   cfg.ApiToken,
		agentId:    cfg.AgentId,
		appKey:     cfg.AppKey,
		appSecret:  cfg.AppSecret,
	}, opts...)
}

func NewPersonClient(cfg PersonConfig, opts ...Option) *DingTalkClient {
	return newDingTalkClient(PERSON, config{
		corpId:     cfg.CorpId,
		corpSecret: cfg.CorpSecret,
		apiToken:   cfg.ApiToken,
		miniAppId:  cfg.MiniAppId,
		appId:      cfg.AppId,
		appSecret:  cfg.AppSecret,
	}, opts...)
}

type options struct {
	cache   Cache
	persist Persist
}

func defaultOptions(clientType ClientType) *options {
	return &options{
		cache:   NewFileCache(fmt.Sprintf(".dingtalk_cache_%d", clientType)),
		persist: NewFileCache(fmt.Sprintf(".dingtalk_persist_%d", clientType)),
	}
}
