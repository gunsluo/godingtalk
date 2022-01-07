package dingtalk

import (
	"context"
	"errors"
	"strings"
)

type SendMessageByTemplateRequest struct {
	AgentId    int
	TemplateId string
	UserIds    []string
	DeptIds    []string
	Data       string
}

type SendMessageByTemplateResponse struct {
	OpenAPIResponse
	TaskId    int    `json:"task_id"`
	RequestId string `json:"request_id"`
}

// 使用模板发送工作通知消息
func (dtc *DingTalkClient) SendMessageByTemplate(ctx context.Context, authCorpId string, req SendMessageByTemplateRequest) (SendMessageByTemplateResponse, error) {
	var data SendMessageByTemplateResponse
	if len(req.UserIds) == 0 && len(req.DeptIds) == 0 {
		return data, errors.New("user Ids and dept ids have at least one")
	}

	requestData := map[string]interface{}{
		"agent_id":    req.AgentId,
		"template_id": req.TemplateId,
	}
	if len(req.UserIds) > 0 {
		requestData["userid_list"] = strings.Join(req.UserIds, ",")
	}
	if len(req.DeptIds) > 0 {
		requestData["dept_id_list"] = strings.Join(req.DeptIds, ",")
	}
	if len(req.Data) > 0 {
		requestData["data"] = req.Data
	}

	err := dtc.httpTOP(ctx, "message/corpconversation/sendbytemplate", authCorpId, nil, requestData, &data)
	return data, err
}
