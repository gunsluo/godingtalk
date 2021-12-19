package dingtalk

import (
	"encoding/json"
	"errors"
)

type PushNotification struct {
	Random       string    `json:"Random"`
	EventType    string    `json:"EventType"`
	TestSuiteKey string    `json:"TestSuiteKey"`
	BizItems     []BizItem `json:"bizData"`
}

type BizItem struct {
	BizType     int    `json:"biz_type"`
	GmtCreate   int    `json:"gmt_create"`
	OpenCursor  int    `json:"open_cursor"`
	SubscribeId string `json:"subscribe_id"`
	Id          int    `json:"id"`
	GmtModified int    `json:"gmt_modified"`
	BizId       string `json:"biz_id"`
	CorpId      string `json:"corp_id"`
	Status      int    `json:"status"`

	BizData interface{} `json:"biz_data"`
}

// UnmarshalJSON unmarshals a biz item.
func (b *BizItem) UnmarshalJSON(data []byte) error {
	var tmpBizItem = struct {
		BizType     int    `json:"biz_type"`
		GmtCreate   int    `json:"gmt_create"`
		OpenCursor  int    `json:"open_cursor"`
		SubscribeId string `json:"subscribe_id"`
		Id          int    `json:"id"`
		GmtModified int    `json:"gmt_modified"`
		BizId       string `json:"biz_id"`
		CorpId      string `json:"corp_id"`
		Status      int    `json:"status"`
		BizData     string `json:"biz_data"`
	}{}

	if err := json.Unmarshal(data, &tmpBizItem); err != nil {
		return err
	}

	b.BizType = tmpBizItem.BizType
	b.GmtCreate = tmpBizItem.GmtCreate
	b.OpenCursor = tmpBizItem.OpenCursor
	b.SubscribeId = tmpBizItem.SubscribeId
	b.Id = tmpBizItem.Id
	b.GmtModified = tmpBizItem.GmtModified
	b.BizId = tmpBizItem.BizId
	b.CorpId = tmpBizItem.CorpId
	b.Status = tmpBizItem.Status

	if fn, ok := bizTypesUnmarshalJSONFuncs[tmpBizItem.BizType]; ok {
		if d, err := fn([]byte(tmpBizItem.BizData)); err == nil {
			b.BizData = d
		} else {
			b.BizData = tmpBizItem.BizData
		}
	} else {
		b.BizData = tmpBizItem.BizData
	}

	return nil
}

func (dtc *DingTalkClient) Decrypt(signature, timestamp, nonce, secretMsg string) (string, error) {
	if dtc.pushCryptoSuit == nil {
		return "", errors.New("please set aes key and token")
	}

	return dtc.pushCryptoSuit.Decrypt(signature, timestamp, nonce, secretMsg)
}

func (dtc *DingTalkClient) Encrypt(msg, timestamp, nonce string) (string, string, error) {
	if dtc.pushCryptoSuit == nil {
		return "", "", errors.New("please set aes key and token")
	}

	return dtc.pushCryptoSuit.Encrypt(msg, timestamp, nonce)
}

func (dtc *DingTalkClient) DecryptAndUnmarshalPushNotification(signature, timestamp, nonce, secretMsg string) (*PushNotification, error) {
	plainTxt, err := dtc.Decrypt(signature, timestamp, nonce, secretMsg)
	if err != nil {
		return nil, err
	}

	notification := &PushNotification{}
	if err := json.Unmarshal([]byte(plainTxt), notification); err != nil {
		return nil, err
	}

	return notification, nil
}
