package usecase

import (
	"cloudflare-r2-viewer/backend/controller"
	"cloudflare-r2-viewer/backend/entity"
	"cloudflare-r2-viewer/backend/repository"
	"context"
	"errors"
	"os"
)

type viewer struct {
	bucketRepository repository.BucketRepository
	objectRepository repository.ObjectRepository
	r2Client         controller.R2Client
}

type Viewer interface {
	GetBucketList(ctx context.Context) ([]*entity.Bucket, error)
	GetObjectList(ctx context.Context, bucketName string) ([]*entity.Object, error)
	SetBucketInfo(ctx context.Context, bucketName string, description, publishDomain string) error
	GetBucketInfo(ctx context.Context, bucketName string) (*entity.Bucket, error)
	FetchThumbnail(ctx context.Context, bucketName, objectPath string) (*[]byte, error)
	FetchImage(ctx context.Context, bucketName, objectPath string) (*[]byte, error)
}

func NewViewer(bucketRepository repository.BucketRepository, objectRepository repository.ObjectRepository, r2Client controller.R2Client) Viewer {
	return &viewer{
		bucketRepository: bucketRepository,
		objectRepository: objectRepository,
		r2Client:         r2Client,
	}
}

func (v *viewer) GetBucketList(ctx context.Context) ([]*entity.Bucket, error) {
	buckets, err := v.r2Client.GetBuckets(ctx)
	if err != nil {
		return nil, err
	}

	bucketInfo, err := v.bucketRepository.GetBucketsInfo()
	if err != nil {
		return nil, err
	}

	bucketList := make([]*entity.Bucket, len(buckets))
	for i, b := range buckets {
		for _, bi := range bucketInfo {
			if b == bi.Name {
				bucketList[i] = &entity.Bucket{
					Name:          b,
					Description:   bi.Description,
					PublishDomain: bi.PublishDomain,
				}
				break
			}
		}
		bucketList[i] = &entity.Bucket{
			Name:          b,
			Description:   "未設定",
			PublishDomain: "未設定",
		}
	}
	return bucketList, nil
}

func (v *viewer) GetObjectList(ctx context.Context, bucketName string) ([]*entity.Object, error) {
	objects, err := v.r2Client.GetObjects(ctx, bucketName)
	if err != nil {
		return nil, err
	}

	return objects, nil
}

func (v *viewer) GetBucketInfo(ctx context.Context, bucketName string) (*entity.Bucket, error) {
	bucketInfo, err := v.bucketRepository.GetBucketInfo(bucketName)
	if err != nil {
		return nil, err
	}

	return bucketInfo, nil
}

func (v *viewer) SetBucketInfo(ctx context.Context, bucketName string, description, publishDomain string) error {
	bucket := &entity.Bucket{
		Name:          bucketName,
		Description:   description,
		PublishDomain: publishDomain,
	}
	err := v.bucketRepository.SetBucketInfo(bucket)
	if err != nil {
		return err
	}

	return nil
}

func (v *viewer) FetchThumbnail(ctx context.Context, bucketName, objectPath string) (*[]byte, error) {
	blob, err := v.objectRepository.FetchLocalThumbnail(objectPath)
	if errors.Is(err, os.ErrNotExist) {
		bucketInfo, err := v.bucketRepository.GetBucketInfo(bucketName)
		if err != nil {
			return nil, err
		}
		publishDomain := bucketInfo.PublishDomain
		blob, err := v.objectRepository.FetchCDNThumbnail(objectPath, publishDomain)
		if err != nil {
			return nil, err
		}
		return blob, nil
	} else if err != nil {
		return nil, err
	}

	return blob, nil
}

func (v *viewer) FetchImage(ctx context.Context, bucketName, objectPath string) (*[]byte, error) {
	blob, err := v.objectRepository.FetchLocalImage(objectPath)
	if errors.Is(err, os.ErrNotExist) {
		bucketInfo, err := v.bucketRepository.GetBucketInfo(bucketName)
		if err != nil {
			return nil, err
		}
		publishDomain := bucketInfo.PublishDomain
		blob, err := v.objectRepository.FetchCDNImage(objectPath, publishDomain)
		if err != nil {
			return nil, err
		}
		return blob, nil
	} else if err != nil {
		return nil, err
	}

	return blob, nil
}
