package dingtalk

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"time"
)

func (dtc *DingTalkClient) httpRPC(ctx context.Context, path string, params url.Values, requestData interface{}, responseData Unmarshallable) error {
	return dtc.httpRequest(ctx, "oapi", path, params, requestData, responseData)
}

func (dtc *DingTalkClient) httpIsv(ctx context.Context, path string, params url.Values, requestData interface{}, responseData Unmarshallable) error {
	if params == nil {
		params = url.Values{}
	}

	switch path {
	case "service/get_corp_token":
		timestamp := fmt.Sprintf("%d", time.Now().UnixNano()/1000000)
		suiteTicket := dtc.GetSuiteTicket(ctx)
		signature := signatureThirdParty(timestamp, suiteTicket, dtc.getAccessSecret())
		params.Set("timestamp", timestamp)
		params.Set("signature", signature)
		params.Set("suiteTicket", suiteTicket)
	case "service/get_suite_token":
		// nothing
	default:
		suiteAccessToken, err := dtc.GetAndRefreshSuiteAccessToken(ctx)
		if err != nil {
			return err
		}
		if params.Get("suite_access_token") == "" {
			params.Set("suite_access_token", suiteAccessToken)
		}
	}

	return dtc.httpRequest(ctx, "oapi", path, params, requestData, responseData)
}

func (dtc *DingTalkClient) httpSNS(ctx context.Context, path string, params url.Values, requestData interface{}, responseData Unmarshallable) error {
	switch path {
	case "sns/getuserinfo_bycode":
		timestamp := fmt.Sprintf("%d", time.Now().UnixNano()/1000000)
		signature := signatureNoSign(timestamp, dtc.getAccessSecret())
		params.Set("timestamp", timestamp)
		params.Set("signature", signature)
	default:
	}

	return dtc.httpRequest(ctx, "oapi", path, params, requestData, responseData)
}

func (dtc *DingTalkClient) httpTOP(ctx context.Context, path, authCorpId string, params url.Values, requestData interface{}, responseData interface{}) error {
	if params == nil {
		params = url.Values{}
	}

	corpAccessToken, err := dtc.IsvGetAndRefreshCorpAccessToken(ctx, authCorpId)
	if err != nil {
		return err
	}
	params.Set("access_token", corpAccessToken)

	return dtc.httpRequest(ctx, "tapi", path, params, requestData, responseData)
}

func (dtc *DingTalkClient) httpRequest(ctx context.Context, tagType string, path interface{}, params url.Values, requestData interface{}, responseData interface{}) error {
	var request *http.Request
	var requestUrl string
	client := dtc.httpClient

	if tagType == "oapi" {
		requestUrl = OAPIURL + path.(string) + "?" + params.Encode()
		fmt.Printf("requestUrl=%s\n", requestUrl)
		if requestData != nil {
			switch v := requestData.(type) {
			case *uploadFile:
				var b bytes.Buffer
				if v.Reader == nil {
					return errors.New("upload file is empty")
				}
				w := multipart.NewWriter(&b)
				fw, err := w.CreateFormFile(v.FieldName, v.FileName)
				if err != nil {
					return err
				}
				if _, err = io.Copy(fw, v.Reader); err != nil {
					return err
				}
				if err = w.Close(); err != nil {
					return err
				}
				request, _ = http.NewRequest("POST", requestUrl, &b)
				request.Header.Set("Content-Type", w.FormDataContentType())
			default:
				d, _ := json.Marshal(requestData)
				request, _ = http.NewRequest("POST", requestUrl, bytes.NewReader(d))
				request.Header.Set("Content-Type", typeJSON+"; charset=UTF-8")
			}
		} else {
			request, _ = http.NewRequest("GET", requestUrl, nil)
		}
	}
	if tagType == "tapi" {
		requestUrl = TOPAPIURL + path.(string) + "?" + params.Encode()
		d, _ := json.Marshal(requestData)
		request, _ = http.NewRequest("POST", requestUrl, bytes.NewReader(d))
		request.Header.Set("Content-Type", typeJSON+"; charset=UTF-8")
	}
	resp, err := client.Do(request)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return errors.New("Server Error: " + resp.Status)
	}
	defer resp.Body.Close()
	contentType := resp.Header.Get("Content-Type")
	pos := len(typeJSON)

	//if tagType == "oapi" {
	if len(contentType) >= pos && contentType[0:pos] == typeJSON {
		if content, err := ioutil.ReadAll(resp.Body); err == nil {
			json.Unmarshal(content, responseData)
			switch responseData.(type) {
			case Unmarshallable:
				resData := responseData.(Unmarshallable)
				return resData.checkError()
			}
		}
	} else {
		switch v := responseData.(type) {
		case *MediaDownloadFileResponse:
			io.Copy(v.Writer, resp.Body)
		}
	}
	//}

	return err
}

func (dtc *DingTalkClient) getAccessKey() string {
	switch dtc.clientType {
	case ISV:
		return dtc.config.suiteKey
	case CORP:
		return dtc.config.appKey
	}

	return ""
}

func (dtc *DingTalkClient) getAccessSecret() string {
	switch dtc.clientType {
	case ISV:
		return dtc.config.suiteSecret
	case CORP:
		return dtc.config.appSecret
	}

	return ""
}
