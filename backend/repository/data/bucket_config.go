package data

import "cloudflare-r2-viewer/backend/entity"

type BucketConfig struct {
	Name          string `json:"name"`           // バケット名
	Description   string `json:"description"`    // 説明
	PublishDomain string `json:"publish_domain"` // 公開ドメイン
}

type BucketConfigFiledata struct {
	Buckets []BucketConfig `json:"buckets"` // バケットの配列
}

func (bcf *BucketConfigFiledata) ConvertToEntities() []*entity.Bucket {
	buckets := make([]*entity.Bucket, len(bcf.Buckets))
	for i, bc := range bcf.Buckets {
		buckets[i] = bc.ConvertToEntity()
	}
	return buckets
}

func (bc *BucketConfig) ConvertToEntity() *entity.Bucket {
	return &entity.Bucket{
		Name:          bc.Name,
		Description:   bc.Description,
		PublishDomain: bc.PublishDomain,
	}
}

func ConvertFromEntity(bc *entity.Bucket) *BucketConfig {
	return &BucketConfig{
		Name:          bc.Name,
		Description:   bc.Description,
		PublishDomain: bc.PublishDomain,
	}
}

func ConvertFromEntities(bcs []*entity.Bucket) *BucketConfigFiledata {
	bucketConfigs := make([]BucketConfig, len(bcs))
	for i, bc := range bcs {
		bucketConfigs[i] = *ConvertFromEntity(bc)
	}
	return &BucketConfigFiledata{
		Buckets: bucketConfigs,
	}
}
