package customerror

import "errors"

var ErrBucketConfigNotFound = errors.New("bucket config not found")
var ErrPublishURLNotSet = errors.New("publish url not set")
var ErrR2ClientNotReady = errors.New("r2 client not ready")
