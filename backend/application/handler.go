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
	SetBucketOptions(ctx context.Context, request request.BucketOption) error
}

func NewHandler(viewerUsecase usecase.Viewer) Handler {
	return &handler{
		viewerUsecase: viewerUsecase,
	}
}

func (h *handler) GetBucketList(ctx context.Context) (response.BucketList, error) {
	bucketList, err := h.viewerUsecase.GetBucketList(ctx)
	if err != nil {
		return response.BucketList{}, err
	}

	bucketNames := make([]string, len(bucketList))
	for i, bucket := range bucketList {
		bucketNames[i] = bucket.Name
	}
	return response.BucketList{
		Buckets: bucketNames,
	}, nil
}

func (h *handler) GetObjectList(ctx context.Context, bucketName string) (response.ObjectList, error) {
	objectList, err := h.viewerUsecase.GetObjectList(ctx, bucketName)
	if err != nil {
		return response.ObjectList{}, err
	}
	bucketInfo, err := h.viewerUsecase.GetBucketInfo(ctx, bucketName)
	if err != nil {
		return response.ObjectList{}, err
	}

	publishURL := ""
	if bucketInfo != nil && bucketInfo.PublishDomain != "" {
		publishURL = "https://" + bucketInfo.PublishDomain
	}

	objects := make([]response.Obj, len(objectList))
	for i, object := range objectList {
		objects[i] = response.Obj{
			Name:       object.Name,
			Path:       object.Path,
			Size:       object.Size,
			PublishURL: publishURL + "/" + object.Path,
		}
	}
	return response.ObjectList{
		Objects: objects,
	}, nil
}

func (h *handler) GetObjectThumbnail(ctx context.Context, bucketName, objectPath string) ([]byte, error) {
	thumbnail, err := h.viewerUsecase.FetchThumbnail(ctx, bucketName, objectPath)
	if err != nil {
		return nil, err
	}
	return *thumbnail, nil
}

func (h *handler) GetObjectImage(ctx context.Context, bucketName, objectPath string) ([]byte, error) {
	image, err := h.viewerUsecase.FetchImage(ctx, bucketName, objectPath)
	if err != nil {
		return nil, err
	}
	return *image, nil
}

func (h *handler) GetBucketOptions(ctx context.Context, bucketName string) (response.BucketOption, error) {
	bucketInfo, err := h.viewerUsecase.GetBucketInfo(ctx, bucketName)
	if err != nil {
		return response.BucketOption{}, err
	}
	if bucketInfo == nil {
		return response.BucketOption{}, nil
	}
	return response.BucketOption{
		BucketName:  bucketInfo.Name,
		Description: bucketInfo.Description,
		PublishURL:  bucketInfo.PublishDomain,
	}, nil
}

func (h *handler) SetBucketOptions(ctx context.Context, request request.BucketOption) error {
	err := h.viewerUsecase.SetBucketInfo(ctx, request.BucketName, request.Description, request.PublishURL)
	if err != nil {
		return err
	}
	return nil
}
