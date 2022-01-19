package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

// Response apigateway에서 받는 이벤트
type Response events.APIGatewayProxyResponse

// RequestData string
type RequestData struct {
	// Subject 제목
	Subject string
	// Message 내용
	Message string
	// Recipient  받는이
	Recipient string
}

const (
	// Sender 자신의 메일주소 입력
	Sender = "\"doyoon Lee\" <leedy@mz.co.kr>"
	// CharSet 인코딩
	CharSet = "UTF-8"
)

// Handler 핸들러
func Handler(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {
	fmt.Printf("Request: %+v\n", request)
	fmt.Printf("Processing request data for request %s.\n", request.RequestContext.RequestID)
	fmt.Printf("Body size = %d.\n", len(request.Body))
	var requestData RequestData
	json.Unmarshal([]byte(request.Body), &requestData)
	fmt.Printf("RequestData: %+v", requestData)
	var result string
	if len(requestData.Subject) > 0 && len(requestData.Message) > 0 {
		result, _ = send(requestData.Subject, requestData.Message, requestData.Recipient)
	}
	resp := events.APIGatewayProxyResponse{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            result,
		Headers: map[string]string{
			"Content-Type":           "application/json",
			"X-MyCompany-Func-Reply": "send-mail-handler",
		},
	}
	return resp, nil
}
func send(Subject string, Message string, Recipient string) (string, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("ap-northeast-2")},
	)
	svc := ses.New(sess)
	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			CcAddresses: []*string{},
			ToAddresses: []*string{
				aws.String(Recipient),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String(CharSet),
					Data:    aws.String(Message),
				},
				Text: &ses.Content{
					Charset: aws.String(CharSet),
					Data:    aws.String(Message),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String(CharSet),
				Data:    aws.String(Subject),
			},
		},
		Source: aws.String(Sender),
	}
	result, err := svc.SendEmail(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ses.ErrCodeMessageRejected:
				fmt.Println(ses.ErrCodeMessageRejected, aerr.Error())
			case ses.ErrCodeMailFromDomainNotVerifiedException:
				fmt.Println(ses.ErrCodeMailFromDomainNotVerifiedException, aerr.Error())
			case ses.ErrCodeConfigurationSetDoesNotExistException:
				fmt.Println(ses.ErrCodeConfigurationSetDoesNotExistException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			fmt.Println(err.Error())
		}
		return "there was an unexpected error", err
	}
	fmt.Println("Email Sent to address: " + Recipient)
	fmt.Println(result)
	return "sent!", err
}
func main() {
	lambda.Start(Handler)
}
