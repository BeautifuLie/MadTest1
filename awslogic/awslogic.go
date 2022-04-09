package awslogic

import (
	"encoding/json"
	"fmt"
	"program/model"
	"program/storage"
)

type AwsServer struct {
	awsServer storage.AWSfuncs
}

func NewAwsServer(as storage.AWSfuncs) *AwsServer {
	s := &AwsServer{
		awsServer: as,
	}
	return s
}

func (a *AwsServer) Report() ([]byte, error) {

	str, err := a.awsServer.ReadS3LambdaReport()
	if err != nil {
		return nil, err
	}

	var rep []model.Report
	err = json.Unmarshal(str, &rep)
	if err != nil {
		return nil, err
	}
	res, err := json.MarshalIndent(rep, "", "   ")
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (a *AwsServer) SendMessage(j model.Joke) (string, error) {
	res, err := a.awsServer.SendMsg(j)
	if err != nil {
		return "", err
	}
	ret := fmt.Sprintf("Message was sended. Message id:%v", res)
	return ret, nil
}
