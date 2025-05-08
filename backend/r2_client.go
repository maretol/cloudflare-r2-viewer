package backend

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/pkg/errors"
)

type r2client struct {
	accountID       string
	accessKeyID     string
	secretAccessKey string
	client          *s3.Client
}

type R2Client interface {
}

func NewR2Client(ctx context.Context, accountID, accessKeyID, secretAccessKey string) R2Client {
	rc := &r2client{
		accountID:       accountID,
		accessKeyID:     accessKeyID,
		secretAccessKey: secretAccessKey,
	}
	rc.init(ctx)
	return rc
}

func (r *r2client) SetAccountID(ctx context.Context, accountID string) {
	r.accountID = accountID
	r.init(ctx)
}

func (r *r2client) SetAccessKeyID(ctx context.Context, accessKeyID string) {
	r.accessKeyID = accessKeyID
	r.init(ctx)
}

func (r *r2client) SetSecretAccessKey(ctx context.Context, secretAccessKey string) {
	r.secretAccessKey = secretAccessKey
	r.init(ctx)
}

func (r *r2client) init(ctx context.Context) {
	hasher := sha256.New()
	hasher.Write([]byte(r.secretAccessKey))
	hashedSecret := hex.EncodeToString(hasher.Sum(nil))

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion("auto"),
		config.WithBaseEndpoint(fmt.Sprintf("https://%s.r2.cloudflarestorage.com", r.accountID)),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(r.accessKeyID, hashedSecret, "")),
	)
	if err != nil {
		log.Fatal(errors.Wrap(err, "failed to load configuration"))
	}

	client := s3.NewFromConfig(cfg)
	r.client = client
}

func (r *r2client) GetBuckets(ctx context.Context) ([]string, error) {
	buckets, err := r.client.ListBuckets(ctx, &s3.ListBucketsInput{})
	if err != nil {
		return nil, errors.Wrap(err, "failed to list buckets")
	}

	bucketNames := make([]string, len(buckets.Buckets))
	for i, bucket := range buckets.Buckets {
		bucketNames[i] = *bucket.Name
	}

	return bucketNames, nil
}

func (r *r2client) GetObjects(ctx context.Context, bucketName string) ([]string, error) {
	objects, err := r.client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket: &bucketName,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to list objects")
	}

	objectKeys := make([]string, len(objects.Contents))
	for i, object := range objects.Contents {
		objectKeys[i] = *object.Key
	}

	return objectKeys, nil
}
