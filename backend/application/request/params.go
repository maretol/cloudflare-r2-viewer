package request

type BucketOption struct {
	BucketName  string `json:"bucket_name"` // バケット名
	Description string `json:"description"` // 説明
	PublishURL  string `json:"publish_url"` // 公開URL
}
