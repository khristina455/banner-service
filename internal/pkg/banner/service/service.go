package service

import (
	"banner-service/internal/models"
	"banner-service/internal/pkg/banner"
	"banner-service/internal/pkg/cache"
	"context"
	"fmt"
	"strconv"
)

type BannerService struct {
	repo  banner.BannerRepository
	cache *cache.RedisClient
}

func NewBannerService(repo banner.BannerRepository, cache *cache.RedisClient) *BannerService {
	return &BannerService{repo: repo, cache: cache}
}

// TODO:решить вопрос с флагом активности

func (bs *BannerService) GetBanner(ctx context.Context, tagId, featureId int, useLastRevision bool, isAdmin bool) ([]byte, error) {
	var banner []byte
	ok := false
	key := strconv.Itoa(tagId) + "-" + strconv.Itoa(featureId)

	if !useLastRevision && bs.cache != nil {
		banner, ok = bs.cache.Get(ctx, key)
	}

	if isAdmin {
		banner, err := bs.repo.ReadBanner(ctx, tagId, featureId)
		if err != nil {
			return nil, err
		}

		if bs.cache != nil {
			bs.cache.Set(key, banner)
		}
		return banner, nil
	}

	if !ok {
		banner, err := bs.repo.ReadUserBanner(ctx, tagId, featureId)
		if err != nil {
			return nil, err
		}

		if bs.cache != nil {
			bs.cache.Set(key, banner)
		}
	}
	return banner, nil
}

func (bs *BannerService) GetFilterBanners(ctx context.Context, tagId, featureId, limit, offset int) ([]models.Banner, error) {
	banners, err := bs.repo.ReadFilterBanners(ctx, tagId, featureId, limit, offset)
	if err != nil {
		return nil, err
	}
	return banners, nil
}

func (bs *BannerService) AddBanner(ctx context.Context, banner *models.BannerPayload) (int, error) {
	if banner.Content == nil || !banner.IsActive.HasValue || banner.TagIds == nil || banner.FeatureId == 0 {
		return 0, fmt.Errorf("incorrect data")
	}

	bannerId, err := bs.repo.CreateBanner(ctx, banner)

	if err != nil {
		return 0, err
	}
	return bannerId, nil
}

func (bs *BannerService) UpdateBanner(ctx context.Context, id int, banner *models.BannerPayload) error {
	err := bs.repo.UpdateBanner(ctx, id, banner)
	return err
}

func (bs *BannerService) DeleteBanner(ctx context.Context, id int) error {
	err := bs.repo.DeleteBanner(ctx, id)
	if err != nil {
		return err
	}
	return nil
}
