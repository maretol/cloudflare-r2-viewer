package repository

import "errors"

var ErrBucketConfigNotFound = errors.New("bucket config not found")
var ErrPublishURLNotSet = errors.New("publish url not set")
