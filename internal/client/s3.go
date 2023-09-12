package client

import (
	"bytes"
	"context"
	"fmt"
	"git.miem.hse.ru/1206/material-library/internal/config"
	"git.miem.hse.ru/1206/material-library/internal/models/tplibrary"
	"github.com/gabriel-vasile/mimetype"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"path"
	"time"
)

const bucketPublicPolicyTemplate = `{
		"Version": "2012-10-17",
		"Statement": [
			{
				"Effect": "Allow",
				"Principal": "*",
				"Action": "s3:GetObject",
				"Resource": "arn:aws:s3:::%s/*"
			}
		]
	}`

type S3Client struct {
	*minio.Client
	cfg *config.S3
}

func NewS3Client(cfg *config.S3) (*S3Client, error) {
	cl, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKeyId, cfg.SecretKey, ""),
		Secure: !cfg.DisableSSL,
		Region: cfg.Region,
	})
	if err != nil {
		return nil, err
	}
	if cfg.UrlLifespan == 0 {
		cfg.UrlLifespan = 5 * time.Minute
	}

	return &S3Client{
		Client: cl,
		cfg:    cfg,
	}, nil
}

func (c *S3Client) PutFile(ctx context.Context, fName string, body []byte) (info *tplibrary.FileInfo, err error) {
	if err = c.createBucket(ctx, c.cfg.BucketName); err != nil {
		return nil, err
	}
	dst := &bytes.Buffer{}
	_, err = dst.Write(body)
	if err != nil {
		return nil, err
	}
	fKey := uuid.New().String()
	fType, fExt := c.getFileType(body)
	i, err := c.PutObject(ctx, c.cfg.BucketName, fKey, dst, int64(dst.Len()), minio.PutObjectOptions{
		ContentType:        fType,
		ContentDisposition: "attachment; filename=" + fName + fExt,
	})
	if err != nil {
		return nil, err
	}
	info = &tplibrary.FileInfo{
		Key:  i.Key,
		Size: i.Size,
		Url:  path.Join(c.cfg.PublicEndpoint, c.cfg.BucketName, fKey),
	}
	return info, nil
}

func (c *S3Client) RemoveFile(ctx context.Context, key string) error {
	if err := c.RemoveObject(ctx, c.cfg.BucketName, key, minio.RemoveObjectOptions{}); err != nil {
		return err
	}
	return nil
}

func (c *S3Client) createBucket(ctx context.Context, bucket string) error {
	if exists, err := c.BucketExists(ctx, bucket); err != nil {
		return err
	} else if exists {
		return nil
	}
	if err := c.MakeBucket(ctx, bucket, minio.MakeBucketOptions{
		Region:        c.cfg.Region,
		ObjectLocking: true,
	}); err != nil {
		return err
	}
	policy := fmt.Sprintf(bucketPublicPolicyTemplate, bucket)
	if err := c.SetBucketPolicy(ctx, bucket, policy); err != nil {
		return err
	}
	return nil
}

func (c *S3Client) getFileType(data []byte) (mime string, ext string) {
	m := mimetype.Detect(data)
	return m.String(), m.Extension()
}
