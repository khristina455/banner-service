package service

import (
	"banner-service/internal/models"
	"banner-service/internal/pkg/banner"
	"banner-service/internal/pkg/cache"
	"context"
	"errors"
	"strconv"
)

type BannerService struct {
	repo  banner.BannerRepository
	cache *cache.RedisClient
}

func NewBannerService(repo banner.BannerRepository, cache *cache.RedisClient) *BannerService {
	return &BannerService{repo: repo, cache: cache}
}

func (bs *BannerService) GetBanner(ctx context.Context, tagID, featureID int,
	useLastRevision bool, isAdmin bool) ([]byte, error) {
	var banner []byte
	ok := false
	key := strconv.Itoa(tagID) + "-" + strconv.Itoa(featureID)

	if !useLastRevision && bs.cache != nil {
		banner, ok = bs.cache.Get(ctx, key)
	}

	var err error
	if isAdmin {
		banner, err = bs.repo.ReadBanner(ctx, tagID, featureID)
		if err != nil {
			return nil, err
		}

		if bs.cache != nil {
			bs.cache.Set(key, banner)
		}
		return banner, nil
	}

	if !ok {
		banner, err = bs.repo.ReadUserBanner(ctx, tagID, featureID)
		if err != nil {
			return nil, err
		}

		if bs.cache != nil {
			bs.cache.Set(key, banner)
		}
	}
	return banner, nil
}

func (bs *BannerService) GetFilterBanners(ctx context.Context,
	tagID, featureID, limit, offset int) ([]models.Banner, error) {
	banners, err := bs.repo.ReadFilterBanners(ctx, tagID, featureID, limit, offset)
	if err != nil {
		return nil, err
	}
	return banners, nil
}

func (bs *BannerService) AddBanner(ctx context.Context, banner *models.BannerPayload) (int, error) {
	if banner.Content == nil || !banner.IsActive.HasValue || banner.TagIDs == nil || banner.FeatureID == 0 {
		return 0, errors.New("has empty fields")
	}

	bannerID, err := bs.repo.CreateBanner(ctx, banner)

	if err != nil {
		return 0, err
	}
	return bannerID, nil
}

func (bs *BannerService) UpdateBanner(ctx context.Context, id int, banner *models.BannerPayload) error {
	err := bs.repo.UpdateBanner(ctx, id, banner)
	return err
}

func (bs *BannerService) DeleteBanner(ctx context.Context, id int) error {
	err := bs.repo.DeleteBanner(ctx, id)
	return err
}

func (bs *BannerService) GetCurrentBanner(ctx context.Context, id int) (models.BannerVersion, error) {
	currentVersion, err := bs.repo.ReadCurrentBannerByID(ctx, id)
	return currentVersion, err
}

func (bs *BannerService) GetOldBanners(ctx context.Context, id int) ([]models.BannerVersion, error) {
	oldVersions, err := bs.repo.ReadOldVersions(ctx, id)
	return oldVersions, err
}

func (bs *BannerService) ChangeVersionOfBanner(ctx context.Context, id int, version int) error {
	err := bs.repo.UpdateVersionOfBanner(ctx, id, version)
	return err
}
