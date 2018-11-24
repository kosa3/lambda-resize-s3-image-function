package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/nfnt/resize"
	"image"
	"image/jpeg"
	"log"
	"strconv"
)

const REGION = "ap-northeast-1"
const BUCKET_NAME = "resize-s3-images-dev-serverlessdeploymentbucket-3sws2q0iypsa"

type Event struct {
	FileName string `json:"filename"`
	Width    string `json:"width"`
	Height   string `json:"height"`
}

func Handler(event Event) (string, error) {
	fmt.Println(event.Width, event.FileName)
	log.Print(event.Width, event.FileName)
	sess, err := session.NewSession(aws.NewConfig().WithRegion(REGION))
	if err != nil {
		fmt.Println(err.Error())
	}

	svc := s3.New(sess)
	var key = event.FileName

	// get
	resp, err := svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(BUCKET_NAME),
		Key:    aws.String(key),
	})

	if err != nil {
		log.Println(err)
	}

	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(resp.Body); err != nil {
		log.Println(err)
	}

	imageBytes := buf.Bytes()
	decodeImg, _, err := image.Decode(bytes.NewReader(imageBytes))
	if err != nil {
		log.Println(err)
	}

	width, _ := strconv.ParseInt(event.Width, 10, 64)
	height, _ := strconv.ParseInt(event.Height, 10, 64)

	thumbnail := resize.Thumbnail(uint(width), uint(height), decodeImg, resize.Lanczos3)
	buf2 := new(bytes.Buffer)
	err = jpeg.Encode(buf2, thumbnail, nil)
	base64EncodingImgSrc := base64.StdEncoding.EncodeToString(buf2.Bytes())

	return base64EncodingImgSrc, nil
}

func main() {
	lambda.Start(Handler)
}
