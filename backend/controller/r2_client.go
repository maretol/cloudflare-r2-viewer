package controller

import (
	"cloudflare-r2-viewer/backend/customerror"
	"cloudflare-r2-viewer/backend/entity"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

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

type ClientConfig struct {
	AccountID       string `json:"account_id"`        // アカウントID
	AccessKeyID     string `json:"access_key_id"`     // アクセスキーID
	SecretAccessKey string `json:"secret_access_key"` // シークレットアクセスキー
}

type R2Client interface {
	SetAccountID(ctx context.Context, accountID string)
	SetAccessKeyID(ctx context.Context, accessKeyID string)
	SetSecretAccessKey(ctx context.Context, secretAccessKey string)
	GetBuckets(ctx context.Context) ([]string, error)
	GetObjects(ctx context.Context, bucketName string) ([]*entity.Object, error)
}

func NewR2Client() R2Client {
	r2client := &r2client{}
	r2client.readConfigFile()
	return r2client
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
	if r.accountID == "" || r.accessKeyID == "" || r.secretAccessKey == "" {
		return
	}
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
	r.writeConfigFile()
}

func (r *r2client) readConfigFile() {
	file, err := os.ReadFile("config.json")
	if errors.Is(err, os.ErrNotExist) {
		r.client = nil
		return
	} else if err != nil {
		log.Fatal(errors.Wrap(err, "failed to read config file"))
	}
	var clientConfig ClientConfig
	err = json.Unmarshal(file, &clientConfig)
	if err != nil {
		log.Fatal(errors.Wrap(err, "failed to unmarshal config file"))
	}
	r.accountID = clientConfig.AccountID
	r.accessKeyID = clientConfig.AccessKeyID
	r.secretAccessKey = clientConfig.SecretAccessKey
}

func (r *r2client) writeConfigFile() {
	clientConfig := ClientConfig{
		AccountID:       r.accountID,
		AccessKeyID:     r.accessKeyID,
		SecretAccessKey: r.secretAccessKey,
	}
	file, err := json.MarshalIndent(clientConfig, "", "  ")
	if err != nil {
		log.Fatal(errors.Wrap(err, "failed to marshal config file"))
	}
	err = os.WriteFile("config.json", file, 0644)
	if err != nil {
		log.Fatal(errors.Wrap(err, "failed to write config file"))
	}
}

func (r *r2client) GetBuckets(ctx context.Context) ([]string, error) {
	if r.client == nil {
		return nil, customerror.ErrR2ClientNotReady
	}

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

func (r *r2client) GetObjects(ctx context.Context, bucketName string) ([]*entity.Object, error) {
	objects, err := r.client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket: &bucketName,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to list objects")
	}

	ents := make([]*entity.Object, len(objects.Contents))
	for i, object := range objects.Contents {
		path := *object.Key
		splitPath := strings.Split(path, "/")
		name := splitPath[len(splitPath)-1]
		ents[i] = &entity.Object{
			Path: path,
			Name: name,
			Size: *object.Size,
		}
	}

	return ents, nil
}
