package awslogic

import (
	"encoding/json"
	"fmt"
	"program/model"
	"program/storage"
	"program/tools"
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
func (a *AwsServer) UploadJokesAndMessagesTos3(j model.Joke) error {
	msgId, err := a.awsServer.SendMsg(j.ID)
	if err != nil {
		return err
	}
	fmt.Printf("Message was sended. Message id:%v", msgId)
	err = tools.CreateAndSaveJokes(j)
	if err != nil {
		return err
	}

	res, err := a.awsServer.GetMsg()
	if err != nil {
		return err
	}
	filename, err := tools.CreateAndSaveMessages(res.GoString())
	if err != nil {
		return err
	}
	err = a.awsServer.UploadMessageTos3(filename)
	if err != nil {
		return err
	}

	err = a.awsServer.UploadTos3(j)
	if err != nil {
		return err
	}
	err = a.awsServer.DeleteMsg(*res.Messages[0].ReceiptHandle)
	if err != nil {
		return err
	}
	return nil
}
func (a *AwsServer) RecieveMessage() (string, error) {
	res, err := a.awsServer.GetMsg()
	if err != nil {
		return "", err
	}
	msg := res.GoString()
	return msg, nil
}
func (a *AwsServer) SendMessage(j model.Joke) (string, error) {
	res, err := a.awsServer.SendMsg(j.ID)
	if err != nil {
		return "", err
	}
	ret := fmt.Sprintf("Message was sended. Message id:%v", res)
	return ret, nil
}
func (a *AwsServer) DeleteMessage(messageHandle string) error {
	err := a.awsServer.DeleteMsg(messageHandle)
	if err != nil {
		return err
	}

	return nil
}
