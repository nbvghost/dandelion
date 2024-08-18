package aliyun

import (
	"encoding/json"
	"errors"
	alimt20181012 "github.com/alibabacloud-go/alimt-20181012/v2/client"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/nbvghost/dandelion/service"
	"github.com/nbvghost/dandelion/service/serviceargument"
	"github.com/nbvghost/tool/object"
)

type Result struct {
	TranslatedText string `json:"translatedText"`
	Error          string `json:"error"`
}

type Index struct {
	config *serviceargument.AliyunConfig
}

var aliyunClient *alimt20181012.Client

func (m *Index) getClient() (*alimt20181012.Client, error) {
	if aliyunClient == nil {
		if m.config == nil {
			m.config = service.Configuration.GetAliyunConfiguration(0)
		}
		// 工程代码泄露可能会导致 AccessKey 泄露，并威胁账号下所有资源的安全性。以下代码示例仅供参考。
		// 建议使用更安全的 STS 方式，更多鉴权访问方式请参见：https://help.aliyun.com/document_detail/378661.html。
		config := &openapi.Config{
			// 必填，请确保代码运行环境设置了环境变量 ALIBABA_CLOUD_ACCESS_KEY_ID。
			AccessKeyId: tea.String(m.config.AccessKeyID),
			// 必填，请确保代码运行环境设置了环境变量 ALIBABA_CLOUD_ACCESS_KEY_SECRET。
			AccessKeySecret: tea.String(m.config.AccessKeySecret),
		}
		// Endpoint 请参考 https://api.aliyun.com/product/alimt
		config.Endpoint = tea.String("mt.aliyuncs.com")

		_result, _err := alimt20181012.NewClient(config)
		if _err != nil {
			return nil, _err
		}
		aliyunClient = _result
		return aliyunClient, nil
	}
	return aliyunClient, nil
}

func (m *Index) Translate(query []string, from, to string) (map[int]string, error) {
	translateMap := make(map[int]string)
	samllMap := make(map[string]string)
	samllLen := 0
	for i := range query {
		if len(query[i]) > 1000 {
			translateText, err := m.translateBase(query[i], from, to)
			if err != nil {
				return nil, err
			}
			translateMap[i] = translateText
		} else {
			if samllLen+len(query[i]) > 8000 || len(samllMap)+1 >= 50 || i == len(query)-1 {
				samllMap[object.ParseString(i)] = query[i]
				mTranslateMap, err := m.translateBatchBase(samllMap, from, to)
				if err != nil {
					return nil, err
				}
				for i2, s := range mTranslateMap {
					translateMap[i2] = s
				}
				samllMap = make(map[string]string)
				samllLen = 0
			} else {
				samllMap[object.ParseString(i)] = query[i]
				samllLen = samllLen + len(query[i])
			}
		}
	}
	return translateMap, nil
}
func (m *Index) translateBase(query string, from, to string) (string, error) {
	client, err := m.getClient()
	if err != nil {
		return "", err
	}
	getBatchTranslateRequest := &alimt20181012.TranslateGeneralRequest{
		FormatType:     tea.String("text"),
		Scene:          tea.String("general"),
		SourceLanguage: tea.String(from),
		SourceText:     tea.String(query),
		TargetLanguage: tea.String(to),
	}
	runtime := &util.RuntimeOptions{}
	res, err := client.TranslateGeneralWithOptions(getBatchTranslateRequest, runtime)
	if err != nil {
		return "", err
	}
	if *res.StatusCode != 200 {
		return "", errors.New("网络错误")
	}
	if *res.Body.Code != 200 {
		return "", errors.New(tea.StringValue(res.Body.Message))
	}
	return tea.StringValue(res.Body.Data.Translated), nil
}
func (m *Index) translateBatchBase(query map[string]string, from, to string) (map[int]string, error) {
	outArr := make(map[int]string)
	client, err := m.getClient()
	if err != nil {
		return outArr, err
	}
	queryJson, err := json.Marshal(query)
	if err != nil {
		return outArr, err
	}
	getBatchTranslateRequest := &alimt20181012.GetBatchTranslateRequest{
		ApiType:        tea.String("translate_standard"),
		FormatType:     tea.String("text"),
		Scene:          tea.String("general"),
		SourceLanguage: tea.String(from),
		SourceText:     tea.String(string(queryJson)),
		TargetLanguage: tea.String(to),
	}
	runtime := &util.RuntimeOptions{}
	res, err := client.GetBatchTranslateWithOptions(getBatchTranslateRequest, runtime)
	if err != nil {
		return outArr, err
	}
	if *res.StatusCode != 200 {
		return outArr, errors.New("网络错误")
	}
	if *res.Body.Code != 200 {
		return outArr, errors.New(tea.StringValue(res.Body.Message))
	}

	trans := res.Body.TranslatedList
	for i := range trans {
		item := trans[i]
		index := object.ParseInt(item["index"])
		translated := object.ParseString(item["translated"])
		outArr[index] = translated
	}
	return outArr, nil
}
func New() *Index {
	return &Index{}
}
