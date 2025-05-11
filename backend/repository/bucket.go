package repository

import (
	"cloudflare-r2-viewer/backend/entity"
	"cloudflare-r2-viewer/backend/repository/data"
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

type bucket struct {
	configData *data.BucketConfigFiledata
}

type BucketRepository interface {
	GetBucketsInfo() ([]*entity.Bucket, error)
	GetBucketInfo(bucketName string) (*entity.Bucket, error)
	SetBucketInfo(bucket *entity.Bucket) error
}

func NewBucketRepository() BucketRepository {
	b := &bucket{}
	b.init()
	return b
}

func (b *bucket) init() {
	file, err := os.ReadFile("config.json")
	if errors.Is(err, os.ErrNotExist) {
		b.configData = &data.BucketConfigFiledata{
			Buckets: []data.BucketConfig{},
		}
		return
	} else if err != nil {
		panic(fmt.Sprintf("failed to read config file: %v", err))
	}
	var bucketConfigFiledata data.BucketConfigFiledata
	err = json.Unmarshal(file, &bucketConfigFiledata)
	if err != nil {
		panic(fmt.Sprintf("failed to unmarshal config file: %v", err))
	}
	b.configData = &bucketConfigFiledata
}

func (b *bucket) saveConfig() error {
	file, err := json.MarshalIndent(b.configData, "", "  ")
	if err != nil {
		return err
	}
	err = os.WriteFile("config.json", file, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (b *bucket) GetBucketsInfo() ([]*entity.Bucket, error) {
	if b.configData == nil {
		b.init()
	}
	return b.configData.ConvertToEntities(), nil
}

func (b *bucket) GetBucketInfo(bucketName string) (*entity.Bucket, error) {
	if b.configData == nil {
		b.init()
	}

	for _, bucketConfig := range b.configData.Buckets {
		if bucketConfig.Name == bucketName {
			return &entity.Bucket{
				Name:          bucketConfig.Name,
				Description:   bucketConfig.Description,
				PublishDomain: bucketConfig.PublishDomain,
			}, nil
		}
	}
	return nil, ErrBucketConfigNotFound
}

func (b *bucket) SetBucketInfo(bucket *entity.Bucket) error {
	config := data.ConvertFromEntity(bucket)
	exist := false
	for i, bc := range b.configData.Buckets {
		if bc.Name == bucket.Name {
			exist = true
			b.configData.Buckets[i] = *config
			break
		}
	}

	if !exist {
		b.configData.Buckets = append(b.configData.Buckets, *config)
	}

	err := b.saveConfig()
	if err != nil {
		return err
	}

	return nil
}
