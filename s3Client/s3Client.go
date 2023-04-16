package s3Client

import (
	"bytes"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"path/filepath"
	"s3-interaction/logger"
)

// S3Client is a wrapper around the AWS SDK
type S3Client struct {
	uploader   *manager.Uploader
	downloader *manager.Downloader
	bucketName string
	logger     *logger.CustomLogger
}

// NewS3Client returns a new S3Client object with the given
// credentials and region set in the config and the given bucket name
func NewS3Client(awsKey, awsSecret, awsRegion string) (*S3Client, error) {
	customLogger, err := logger.NewCustomLogger("router", true)
	if err != nil {
		return nil, err
	}

	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			awsKey,
			awsSecret,
			"",
		)),
		config.WithRegion(awsRegion),
	)
	if err != nil {
		return nil, err
	}

	s3Client := s3.NewFromConfig(cfg)

	return &S3Client{
		uploader:   manager.NewUploader(s3Client),
		downloader: manager.NewDownloader(s3Client),
		logger:     customLogger,
	}, nil
}

// UploadFile uploads a file to S3
func (s *S3Client) UploadFile(bucket string, fileContent []byte, key string) error {
	s.logger.Info(fmt.Sprintf("Uploading file to S3: %s", filepath.Join(bucket, key)))
	_, err := s.uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   bytes.NewReader(fileContent),
	})
	return err
}

// DownloadFile downloads a file from S3
func (s *S3Client) DownloadFile(bucket, key string) ([]byte, error) {
	s.logger.Info(fmt.Sprintf("Downloading file from S3: %s", filepath.Join(bucket, key)))
	buff := manager.NewWriteAtBuffer([]byte{})
	_, err := s.downloader.Download(context.TODO(), buff, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, err
	}
	return buff.Bytes(), nil
}
