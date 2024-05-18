package service

import (
	"context"
	"halo-suster/internal/pkg/dto"
	"halo-suster/internal/pkg/errs"
	"mime/multipart"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func (s *Service) UploadImageProcess(ctx context.Context, file *multipart.FileHeader) errs.Response {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(s.Config().S3Config.Region),
		Credentials: credentials.NewStaticCredentials(
			s.Config().S3Config.ID,
			s.Config().S3Config.Secret,
			"",
		),
	})
	if err != nil {
		return errs.NewInternalError("failed to create aws session", err)
	}

	fileContent, err := file.Open()
	if err != nil {
		return errs.NewInternalError("failed to open file", err)
	}
	defer fileContent.Close()

	output, err := s3.New(sess).PutObject(&s3.PutObjectInput{
		ACL:    aws.String("public-read"),
		Body:   fileContent,
		Bucket: aws.String(s.Config().S3Config.Bucket),
		Key:    aws.String(file.Filename),
	})
	if err != nil {
		return errs.NewInternalError("failed to upload image", err)
	}

	// get the URL of the uploaded file
	_ = output

	return errs.Response{
		Code:    http.StatusOK,
		Message: "File uploaded successfully",
		Data:    dto.ImageResponse{
			// ImageUrl: s.Config().S3Config.BaseURL + file.Filename,
		},
	}
}
