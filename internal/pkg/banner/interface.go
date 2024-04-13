package banner

import (
	"banner-service/internal/models"
	"context"
)

type BannerService interface {
	GetBanner(ctx context.Context, tagId, featureId int, useLastRevision bool, isAdmin bool) ([]byte, error)
	GetFilterBanners(ctx context.Context, tagId, featureId, limit, offset int) ([]models.Banner, error)
	AddBanner(ctx context.Context, banner *models.BannerPayload) (int, error)
	UpdateBanner(ctx context.Context, id int, banner *models.BannerPayload) error
	DeleteBanner(ctx context.Context, id int) error
}

type BannerRepository interface {
	ReadBanner(ctx context.Context, tagId, featureId int) ([]byte, error)
	ReadUserBanner(ctx context.Context, tagId, featureId int) ([]byte, error)
	ReadFilterBanners(ctx context.Context, tagId, featureId, limit, offset int) ([]models.Banner, error)
	CreateBanner(ctx context.Context, banner *models.BannerPayload) (int, error)
	UpdateBanner(ctx context.Context, id int, banner *models.BannerPayload) error
	DeleteBanner(ctx context.Context, id int) error
}
