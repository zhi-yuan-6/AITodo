package util

import (
	"AITodo/config"
	"AITodo/globals"
	"encoding/json"
	"fmt"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dysmsapi "github.com/alibabacloud-go/dysmsapi-20170525/v4/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"log"
)

type SMSService struct {
	client *dysmsapi.Client
	cfg    *config.SMSConfig
}

func NewSMSService(cfg *config.SMSConfig) (*SMSService, error) {
	client, err := dysmsapi.NewClient(&openapi.Config{
		AccessKeyId:     tea.String(cfg.AccessKeyID),
		AccessKeySecret: tea.String(cfg.AccessKeySecret),
		Endpoint:        tea.String(cfg.Endpoint),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create SMS client: %w", err)
	}

	return &SMSService{
		client: client,
		cfg:    cfg,
	}, nil
}

func (s *SMSService) SendVerificationCodeAsync(phone, code string) error {
	// 提交任务到协程池
	return globals.TaskPool.Submit(func() {
		err := s.SendVerificationCode(phone, code)
		if err != nil {
			// 异步记录错误（示例使用标准库，建议替换为结构化日志）
			log.Printf("异步短信发送失败: phone=%s error=%v", phone, err)
		}
	})
}

func (s *SMSService) SendVerificationCode(phone, code string) error {
	templateParam, _ := json.Marshal(map[string]string{"code": code})

	request := &dysmsapi.SendSmsRequest{
		PhoneNumbers:  tea.String(phone),
		SignName:      tea.String(s.cfg.SignName),
		TemplateCode:  tea.String(s.cfg.TemplateCode),
		TemplateParam: tea.String(string(templateParam)),
	}

	response, err := s.client.SendSmsWithOptions(request, &util.RuntimeOptions{})
	if err != nil {
		return handleSDKError(err)
	}

	if *response.Body.Code != "OK" {
		return fmt.Errorf("%s", *response.Body.Message)
	}

	return nil
}

func handleSDKError(err error) error {
	if sdkErr, ok := err.(*tea.SDKError); ok {
		var data map[string]interface{}
		if jsonErr := json.Unmarshal([]byte(tea.StringValue(sdkErr.Data)), &data); jsonErr == nil {
			if recommend, exists := data["Recommend"]; exists {
				return fmt.Errorf("%s (建议: %s)", tea.StringValue(sdkErr.Message), recommend)
			}
		}
	}
	return fmt.Errorf("sms service error: %w", err)
}

func VerifySMS(phoneNumber, userCode string) bool {
	// 从 Redis 获取存储的验证码
	storedCode, err := GetCode(phoneNumber)

	if err != nil {
		fmt.Printf("Error retrieving code from Redis: %v\n", err)
		return false
	}

	// 比对验证码
	return storedCode == userCode
}
