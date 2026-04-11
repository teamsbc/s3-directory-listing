package s3

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Client struct {
	s3Client *s3.Client
	bucket   string
}

func NewClient(ctx context.Context, bucket string, profile string) (*Client, error) {
	// Load config with custom options
	optFns := []func(*config.LoadOptions) error{}

	// Use specified profile if provided
	if profile != "" {
		optFns = append(optFns, config.WithSharedConfigProfile(profile))
	}

	// If region is not set or is "auto", default to us-east-1
	// (CloudFlare R2 uses "auto" but AWS SDK needs a valid region)
	region := os.Getenv("AWS_REGION")
	if region == "" || region == "auto" {
		region = "us-east-1"
	}
	optFns = append(optFns, config.WithRegion(region))

	cfg, err := config.LoadDefaultConfig(ctx, optFns...)
	if err != nil {
		return nil, err
	}

	// Configure S3 client options
	s3Opts := []func(*s3.Options){
		// Use path-style addressing (required for R2 and some S3-compatible services)
		func(o *s3.Options) {
			o.UsePathStyle = true
		},
	}

	return &Client{
		s3Client: s3.NewFromConfig(cfg, s3Opts...),
		bucket:   bucket,
	}, nil
}

type Object struct {
	Key  string
	Size int64
}

func (c *Client) ListAllObjects(ctx context.Context) ([]Object, error) {
	var objects []Object
	paginator := s3.NewListObjectsV2Paginator(c.s3Client, &s3.ListObjectsV2Input{
		Bucket: aws.String(c.bucket),
	})

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, obj := range page.Contents {
			objects = append(objects, Object{
				Key:  aws.ToString(obj.Key),
				Size: aws.ToInt64(obj.Size),
			})
		}
	}

	return objects, nil
}
