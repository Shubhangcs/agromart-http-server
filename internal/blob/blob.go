package blob

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/shubhangcs/agromart-server/internal/env"
)

type AWSS3 struct {
	accessKey  string
	secretKey  string
	region     string
	bucketName string
	client     *s3.Client
}

func Connect() (*AWSS3, error) {
	var (
		accessKey  = env.GetString("ACCESS_KEY", "")
		secretKey  = env.GetString("SECRET_KEY", "")
		region     = env.GetString("REGION", "ap-south-1")
		bucketName = env.GetString("BUCKET_NAME", "")
	)
	creds := credentials.NewStaticCredentialsProvider(
		accessKey,
		secretKey,
		"",
	)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	
	awsCfg, err := config.LoadDefaultConfig(
		ctx,
		config.WithRegion(region),
		config.WithCredentialsProvider(creds),
	)
	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(
		awsCfg,
	)
	return &AWSS3{
		accessKey:  accessKey,
		secretKey:  secretKey,
		client:     client,
		bucketName: bucketName,
		region:     region,
	}, nil
}

func (cfg *AWSS3) GenerateUploadPresignedURL(key string) (string, error) {

	presignClient := s3.NewPresignClient(cfg.client)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req, err := presignClient.PresignPutObject(
		ctx,
		&s3.PutObjectInput{
			Bucket: &cfg.bucketName,
			Key:    &key,
		},
		s3.WithPresignExpires(15*time.Minute),
	)
	if err != nil {
		return "", err
	}

	return req.URL, nil
}

func (cfg *AWSS3) DeleteImage(key string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	_, err := cfg.client.DeleteObject(
		ctx,
		&s3.DeleteObjectInput{
			Bucket: &cfg.bucketName,
			Key:    &key,
		},
	)

	if err != nil {
		return err
	}

	return nil
}
