package application

import (
	"cloudflare-r2-viewer/backend/application/request"
	"cloudflare-r2-viewer/backend/application/response"
	"cloudflare-r2-viewer/backend/usecase"
	"context"
)

type handler struct {
	viewerUsecase usecase.Viewer
}

type Handler interface {
	GetBucketList(ctx context.Context) (response.BucketList, error)
	GetObjectList(ctx context.Context, bucketName string) (response.ObjectList, error)
	GetObjectThumbnail(ctx context.Context, bucketName, objectPath string) ([]byte, error)
	GetObjectImage(ctx context.Context, bucketName, objectPath string) ([]byte, error)
	GetBucketOptions(ctx context.Context, bucketName string) (response.BucketOption, error)
	SetBucketOptions(ctx context.Context, request request.BucketOption)
}

func NewHandler(viewerUsecase usecase.Viewer) Handler {
	return &handler{
		viewerUsecase: viewerUsecase,
	}
}

func (h *handler) GetBucketList(ctx context.Context) (response.BucketList, error) {
	return response.BucketList{}, nil
}

func (h *handler) GetObjectList(ctx context.Context, bucketName string) (response.ObjectList, error) {
	return response.ObjectList{}, nil
}

func (h *handler) GetObjectThumbnail(ctx context.Context, bucketName, objectPath string) ([]byte, error) {
	return nil, nil
}

func (h *handler) GetObjectImage(ctx context.Context, bucketName, objectPath string) ([]byte, error) {
	return nil, nil
}

func (h *handler) GetBucketOptions(ctx context.Context, bucketName string) (response.BucketOption, error) {
	return response.BucketOption{}, nil
}

func (h *handler) SetBucketOptions(ctx context.Context, request request.BucketOption) {

}
