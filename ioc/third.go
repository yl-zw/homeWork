package ioc

import (
	"os"
	"webbook/ailiyunmsm/msm"
)

func NewThird() *msm.ALiMessageClient {
	AccessKeyId, ok := os.LookupEnv("AccessKeyId")
	if !ok {
		panic("AccessKeyId")
		return nil
	}
	AccessKeySecret, ok := os.LookupEnv("AccessKeySecret")
	if !ok {
		panic("AccessKeySecret")
		return nil
	}
	clientMSM, err := msm.NewClientMSM(AccessKeyId, AccessKeySecret)
	if err != nil {
		panic(err)
		return nil
	}
	return clientMSM
}
