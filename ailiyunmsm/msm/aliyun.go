package msm

import (
	"errors"
	"fmt"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dysmsapi "github.com/alibabacloud-go/dysmsapi-20170525/v3/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
)

type ALiMessageClient struct {
	client *dysmsapi.Client
}

func NewClientMSM(accessKeyId string, accessKeySecret string) (*ALiMessageClient, error) {

	config := &openapi.Config{
		AccessKeyId:     tea.String(accessKeyId),
		AccessKeySecret: tea.String(accessKeySecret),
		Endpoint:        tea.String("dysmsapi.aliyuncs.com"),
	}
	client, err := dysmsapi.NewClient(config)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &ALiMessageClient{client: client}, nil
}
func (c *ALiMessageClient) Send(singerName string, code string, iphoneNums ...string) error {
	for _, v := range iphoneNums {
		requset := dysmsapi.SendSmsRequest{
			PhoneNumbers:  tea.String(v),
			SignName:      tea.String(singerName),
			TemplateCode:  tea.String("SMS_290778069"),
			TemplateParam: tea.String(fmt.Sprintf(`{"code":%s}`, code)),
		}
		result, err := c.client.SendSmsWithOptions(&requset, &util.RuntimeOptions{})
		if err != nil {
			fmt.Println(err)
			return err
		}
		if *result.Body.Code != "ok" {
			fmt.Println(result)
			return errors.New("系统错误")
		}
		fmt.Println(result.Body)

	}

	return nil
}
