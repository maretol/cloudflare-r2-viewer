package backend

import (
	"cloudflare-r2-viewer/backend/application"
	"cloudflare-r2-viewer/backend/controller"
	"cloudflare-r2-viewer/backend/repository"
	"cloudflare-r2-viewer/backend/usecase"
)

func NewViewerHandler() application.Handler {
	r2Client := controller.NewR2Client()
	bucketRepository := repository.NewBucketRepository()
	objectRepository := repository.NewObjectRepository()

	usecase := usecase.NewViewer(bucketRepository, objectRepository, r2Client)

	return application.NewHandler(usecase)
}
