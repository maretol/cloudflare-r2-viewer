package response

type BucketList struct {
	Buckets []string `json:"buckets"` // バケット名
}

type ObjectList struct {
	Objects []Object `json:"objects"` // オブジェクトの配列
}

type Object struct {
	Name       string `json:"name"`        // ファイル名
	Path       string `json:"path"`        // パス
	Size       int64  `json:"size"`        // サイズ
	PublishURL string `json:"publish_url"` // 公開URL
}

type BucketOption struct {
	BucketName  string `json:"bucket_name"` // バケット名
	Description string `json:"description"` // 説明
	PublishURL  string `json:"publish_url"` // 公開URL
}
