package aliyun

import (
	"encoding/json"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v3/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/service/configuration"
	"github.com/pkg/errors"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
)

type SMS struct {
	AccessKeyID          string
	AccessKeySecret      string
	SignName             string
	ConfigurationService configuration.ConfigurationService
}

func NewSMS(oid dao.PrimaryKey) *SMS {
	sms := &SMS{}
	c := sms.ConfigurationService.GetConfigurations(oid, model.ConfigurationKeyAliyunAccessKeyID, model.ConfigurationKeyAliyunAccessKeySecret, model.ConfigurationKeyAliyunSMSSignName)
	sms.AccessKeyID = c[model.ConfigurationKeyAliyunAccessKeyID]
	sms.AccessKeySecret = c[model.ConfigurationKeyAliyunAccessKeySecret]
	sms.SignName = c[model.ConfigurationKeyAliyunSMSSignName]
	return sms
}

func (m *SMS) Send(templateCode string, phone string, data map[string]any) error {
	c := &openapi.Config{
		// 必填，您的 AccessKey ID
		AccessKeyId: core.String(m.AccessKeyID),
		// 必填，您的 AccessKey Secret
		AccessKeySecret: core.String(m.AccessKeySecret),
	}
	// Endpoint 请参考 https://api.aliyun.com/product/Dysmsapi
	c.Endpoint = tea.String("dysmsapi.aliyuncs.com")
	//_result := &client.Client{}
	newClient, err := dysmsapi20170525.NewClient(c)
	if err != nil {
		return err
	}

	dataBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	sendSmsRequest := &dysmsapi20170525.SendSmsRequest{
		PhoneNumbers:  tea.String(phone),
		SignName:      tea.String(m.SignName),
		TemplateCode:  tea.String(templateCode),
		TemplateParam: tea.String(string(dataBytes)),
	}

	runtime := &util.RuntimeOptions{}
	// 复制代码运行请自行打印 API 的返回值
	// 返回值为 Map 类型，可从 Map 中获得三类数据：响应体 body、响应头 headers、HTTP 返回的状态码 statusCode。
	result, err := newClient.SendSmsWithOptions(sendSmsRequest, runtime)
	if err != nil {
		return err
	}
	if *(result.Body.Code) != "OK" {
		return errors.New(*result.Body.Message)
	}
	return nil
}
