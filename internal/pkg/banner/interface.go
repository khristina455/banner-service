package banner

import (
	"banner-service/internal/models"
	"context"
)

type BannerService interface {
	GetUserBanner(ctx context.Context, tagId, featureId int, useLastRevision bool) ([]byte, error)
	GetFilterBanners(ctx context.Context, tagId, featureId, limit, offset int) ([]models.Banner, error)
	AddBanner(ctx context.Context, banner *models.Banner) (int, error)
}

type BannerRepository interface {
	ReadBanner(ctx context.Context, tagId, featureId int) ([]byte, error)
	ReadFilterBanners(ctx context.Context, tagId, featureId, limit, offset int) ([]models.Banner, error)
	AddBanner(ctx context.Context, banner *models.Banner) (int, error)
}
