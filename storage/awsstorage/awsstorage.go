package awsstorage

import (
	"encoding/json"
	"fmt"
	"os"
	"program/model"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/service/sqs"
	_ "gocloud.dev/blob/s3blob"
)

type AWSFs struct {
	awss3       *s3.S3
	sqsSess     *sqs.SQS
	bucketRead  string
	bucketWrite string
	bucketSQS   string
	urlQueue    string
}

func NewAwsStorage(region, id, secret, token string) (*AWSFs, error) {

	session, err := session.NewSession(
		&aws.Config{
			Region: &region,
			Credentials: credentials.NewStaticCredentials(
				id,
				secret,
				token,
			),
		})

	if err != nil {
		return nil, err
	}
	awss3 := s3.New(session)
	sqsSess := sqs.New(session)
	conn := &AWSFs{
		awss3:       awss3,
		sqsSess:     sqsSess,
		bucketRead:  os.Getenv("BUCKET_READ_NAME"),
		bucketWrite: os.Getenv("BUCKET_WRITE_NAME"),
		bucketSQS:   *aws.String("jokes-sqs-messages"),
		urlQueue:    *aws.String("https://sqs.eu-central-1.amazonaws.com/333746971525/JokesQueueSend"),
	}
	return conn, nil
}
func (a *AWSFs) UploadTos3(j model.Joke) error {
	filename := j.ID + ".json"
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	uploader := s3manager.NewUploaderWithClient(a.awss3)
	upParams := &s3manager.UploadInput{
		Bucket: &a.bucketRead,
		Key:    &filename,
		Body:   f,
	}
	_, err = uploader.Upload(upParams, func(u *s3manager.Uploader) {
		u.LeavePartsOnError = true
	})
	if err != nil {
		return err
	}

	return nil
}
func (a *AWSFs) ReadS3LambdaReport() ([]byte, error) {

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	client := lambda.New(sess, &aws.Config{Region: aws.String(os.Getenv("AWS_REGION"))})
	input := &lambda.InvokeInput{
		FunctionName: aws.String("Monthly_report"),
	}
	output, err := client.Invoke(input)
	if err != nil {
		return nil, err
	}
	result := output.Payload
	return result, nil

}
func (a *AWSFs) SendMsg(j model.Joke) (string, error) {
	str, err := json.Marshal(j)
	if err != nil {
		return "", err
	}
	// id := "id"
	res, err := a.sqsSess.SendMessage(&sqs.SendMessageInput{
		QueueUrl:    &a.urlQueue,
		MessageBody: aws.String(string(str)), ////joke
		// MessageGroupId: &id,
	})

	if err != nil {
		return "", err
	}
	msgId := *res.MessageId
	return msgId, nil
}
func (a *AWSFs) GetQueueUrl(queueName string) (string, error) {

	result, err := a.sqsSess.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: &queueName,
	})

	if err != nil {
		return "", err
	}

	return *result.QueueUrl, nil
}
func (a *AWSFs) GetMsg() (*sqs.ReceiveMessageOutput, error) {

	msgResult, err := a.sqsSess.ReceiveMessage(&sqs.ReceiveMessageInput{
		QueueUrl:            &a.urlQueue,
		MaxNumberOfMessages: aws.Int64(1),
		WaitTimeSeconds:     aws.Int64(10), /////timeout
	})
	// if len(msgResult.Messages) == 0 {
	// 	return nil, fmt.Errorf("len sqs message is 0: %v", err)
	// }
	if err != nil {
		return nil, fmt.Errorf("failed to fetch sqs message %v", err)
	}

	return msgResult, nil

}
func (a *AWSFs) DeleteMsg(messageHandle string) error {

	_, err := a.sqsSess.DeleteMessage(&sqs.DeleteMessageInput{
		QueueUrl:      &a.urlQueue,
		ReceiptHandle: &messageHandle,
	})

	if err != nil {
		return err
	}
	return nil
}
func (a *AWSFs) UploadMessageTos3(filename string) error {

	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	uploader := s3manager.NewUploaderWithClient(a.awss3)
	upParams := &s3manager.UploadInput{
		Bucket: &a.bucketSQS,
		Key:    &filename,
		Body:   f,
	}
	_, err = uploader.Upload(upParams, func(u *s3manager.Uploader) {
		u.LeavePartsOnError = true
	})
	if err != nil {
		return err
	}

	return nil
}
