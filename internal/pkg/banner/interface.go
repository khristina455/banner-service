package banner

import (
	"banner-service/internal/models"
	"context"
)

type BannerService interface {
	GetBanner(ctx context.Context, tagID, featureID int, useLastRevision bool, isAdmin bool) ([]byte, error)
	GetFilterBanners(ctx context.Context, tagID, featureID, limit, offset int) ([]models.Banner, error)
	AddBanner(ctx context.Context, banner *models.BannerPayload) (int, error)
	UpdateBanner(ctx context.Context, id int, banner *models.BannerPayload) error
	DeleteBanner(ctx context.Context, id int) error
	GetCurrentBanner(ctx context.Context, id int) (models.BannerVersion, error)
	GetOldBanners(ctx context.Context, id int) ([]models.BannerVersion, error)
	ChangeVersionOfBanner(ctx context.Context, id int, version int) error
}

type BannerRepository interface {
	ReadBanner(ctx context.Context, tagID, featureID int) ([]byte, error)
	ReadUserBanner(ctx context.Context, tagID, featureID int) ([]byte, error)
	ReadFilterBanners(ctx context.Context, tagID, featureID, limit, offset int) ([]models.Banner, error)
	CreateBanner(ctx context.Context, banner *models.BannerPayload) (int, error)
	UpdateBanner(ctx context.Context, id int, banner *models.BannerPayload) error
	DeleteBanner(ctx context.Context, id int) error
	ReadCurrentBannerByID(ctx context.Context, id int) (models.BannerVersion, error)
	ReadOldVersions(ctx context.Context, id int) ([]models.BannerVersion, error)
	UpdateVersionOfBanner(ctx context.Context, id int, version int) error
}
