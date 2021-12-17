package dingtalk

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"reflect"
)

func (dtc *DingTalkClient) httpIsv(path string, params url.Values, requestData interface{}, responseData Unmarshallable) error {
	if params == nil {
		params = url.Values{}
	}

	switch path {
	case "service/get_corp_token":
	case "service/get_suite_token":
	default:
		suiteAccessToken, err := dtc.GetAndRefreshSuiteAccessToken()
		if err != nil {
			return err
		}
		if params.Get("suite_access_token") == "" {
			params.Set("suite_access_token", suiteAccessToken)
		}
	}
	/*
		switch path {
		case "service/get_corp_token":
			if params == nil {
				params = url.Values{}
			}
			if params.Get("accessKey") == "" {
				params.Set("accessKey", dtc.DTConfig.SuiteKey)
			}

			var suiteTicket string
			if suiteTicket != "" {
				if params.Get("suiteTicket") == "" {
					params.Set("suiteTicket", suiteTicket)
				}

				timestamp := fmt.Sprintf("%d", time.Now().UnixNano()/1000000)
				signature := signatureNoSign(timestamp+"\n"+suiteTicket, dtc.DTConfig.SNSSecret)
				params.Set("signature", signature)
			}
		case "service/get_permanent_code", "service/activate_suite", "service/get_auth_info":
			if dtc.SuiteAccessToken != "" {
				if params == nil {
					params = url.Values{}
				}
				if params.Get("suite_access_token") == "" {
					params.Set("suite_access_token", dtc.SuiteAccessToken)
				}
			}
		default:
			cur := isvGetCInfo[0]
			switch v := cur.(type) {
			case *DTIsvGetCompanyInfo:
				if v.AuthAccessToken != "" {
					if params == nil {
						params = url.Values{}
					}
					if params.Get("access_token") == "" {
						params.Set("access_token", v.AuthAccessToken)
					}
				}
			default:
				panic(errors.New("ERROR: *DTIsvGetCompanyInfo Error"))
			}
		}
	*/

	return dtc.httpRequest("oapi", path, params, requestData, responseData)
}

func (dtc *DingTalkClient) httpRequest(tagType string, path interface{}, params url.Values, requestData interface{}, responseData interface{}) error {
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

	if tagType == "oapi" {
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
	}
	if tagType == "tapi" {
		if content, err := ioutil.ReadAll(resp.Body); err == nil {
			v := reflect.ValueOf(responseData)
			v = v.Elem()
			v.SetBytes(content)
		}
	}
	return err
}
