package volcengine

import (
	"encoding/json"
	"errors"
	"github.com/nbvghost/dandelion/service"
	"github.com/nbvghost/dandelion/service/serviceargument"
	"github.com/volcengine/volc-sdk-golang/base"
	"net/http"
	"net/url"
	"time"
)

func New() *Index {
	return &Index{}
}

type Index struct {
	config *serviceargument.Volcengine
	client *base.Client
}

type Request struct {
	SourceLanguage string   `json:"SourceLanguage"`
	TargetLanguage string   `json:"TargetLanguage"`
	TextList       []string `json:"TextList"`
}
type Response struct {
	TranslationList []struct {
		Translation            string `json:"Translation"`
		DetectedSourceLanguage string `json:"DetectedSourceLanguage"`
	} `json:"TranslationList"`
	ResponseMetadata struct {
		RequestId string `json:"RequestId"`
		Action    string `json:"Action"`
		Version   string `json:"Version"`
		Service   string `json:"Service"`
		Region    string `json:"Region"`
		Error     struct {
			Code    string
			Message string
		} `json:"Error"`
	} `json:"ResponseMetadata"`
}

func (m *Index) translateBase(query []string, from, to string) ([]string, error) {
	req := Request{
		SourceLanguage: from,
		TargetLanguage: to,
		TextList:       query,
	}
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	client, err := m.getClient()
	if err != nil {
		return nil, err
	}

	respBytes, _, err := client.Json("TranslateText", nil, string(body))
	if err != nil {
		return nil, err
	}
	var response Response
	err = json.Unmarshal(respBytes, &response)
	if err != nil {
		return nil, err
	}
	if response.ResponseMetadata.Error.Code != "" {
		return nil, errors.New(response.ResponseMetadata.Error.Message)
	}

	outArr := make([]string, 0)
	for i := range response.TranslationList {
		outArr = append(outArr, response.TranslationList[i].Translation)
	}
	return outArr, nil
}
func (m *Index) getClient() (*base.Client, error) {
	if m.client == nil {
		if m.config == nil {
			m.config = service.Configuration.GetVolcengineConfiguration(0)
		}

		var (
			ServiceInfo = &base.ServiceInfo{
				Timeout: 5 * time.Second,
				Host:    "translate.volcengineapi.com",
				Header: http.Header{
					"Accept": []string{"application/json"},
				},
				Credentials: base.Credentials{Region: base.RegionCnNorth1, Service: "translate"},
			}
			ApiInfoList = map[string]*base.ApiInfo{
				"TranslateText": {
					Method: http.MethodPost,
					Path:   "/",
					Query: url.Values{
						"Action":  []string{"TranslateText"},
						"Version": []string{"2020-06-01"},
					},
				},
			}
		)

		m.client = base.NewClient(ServiceInfo, ApiInfoList)
		m.client.SetAccessKey(m.config.AccessKeyID)
		m.client.SetSecretKey(m.config.AccessKeySecret)
	}
	return m.client, nil
}
func (m *Index) Translate(query []string, from, to string) (map[int]string, error) {
	translateMap := make(map[int]string)

	smallArr := make([]string, 0)
	smallLen := 0

	getNextItemLen := func(i int) int {
		ii := i + 1
		if ii > len(query)-1 {
			return 0
		} else {
			return len(query[ii])
		}
	}

	for i := range query {

		if len(query[i])+getNextItemLen(i) > 5000 {
			translates, err := m.translateBase([]string{query[i]}, from, to)
			if err != nil {
				return nil, err
			}
			translateMap[i] = translates[0]
			continue
		}

		smallArr = append(smallArr, query[i])
		smallLen = smallLen + len(query[i])

		if smallLen+getNextItemLen(i) > 5000 || len(smallArr) >= 16 || i == len(query)-1 {
			translates, err := m.translateBase(smallArr, from, to)
			if err != nil {
				return nil, err
			}
			for ii := range translates {
				translateMap[i-(len(translates)-1)+ii] = translates[ii]
			}
			smallArr = nil
			smallLen = 0
		}
	}
	return translateMap, nil
}
