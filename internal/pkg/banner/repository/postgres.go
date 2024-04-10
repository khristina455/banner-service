package repository

import (
	"banner-service/internal/models"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"strconv"
)

const (
	getBannerIdByTagFeature = `SELECT banner_id FROM banner_tag_feature WHERE tag_id=$1 AND feature_id=$2;`
	getBannerIdsByTag       = `SELECT banner_id FROM banner_tag_feature WHERE tag_id=$1`
	getBannerIdsByFeature   = `SELECT banner_id FROM banner_tag_feature WHERE feature_id=$1`
	getAllBannerIds         = `SELECT banner_id FROM banner_tag_feature`
	getBannerContentById    = `SELECT content FROM banner WHERE banner_id=$1 AND is_active=TRUE;`
	getBannerById           = `SELECT content, is_active, created_at, updated_at FROM banner WHERE banner_id=$1;`
	getFeatureForBanner     = `SELECT feature_id FROM banner_tag_feature WHERE banner_id=$1;`
	getTagsForBanner        = `SELECT tag_id FROM banner_tag_feature WHERE banner_id=$1;`
	createBanner            = `INSERT INTO banner(content, is_active) VALUES ($1, $2) RETURNING banner_id;`
	createFeatureAndTag     = `INSERT INTO banner_tag_feature(banner_id, tag_id, feature_id) VALUES ($1, $2, $3);`
	updateBanner            = `UPDATE banner SET content = COALESCE($1, content),
                			   is_active= COALESCE($2, is_active), 
                               updated_at = now() 
					           WHERE banner_id = $3;`
	deleteBanner = `DELETE FROM banner WHERE banner_id=$1;`
)

var (
	ErrBannerNotFound = errors.New("banner not found")
)

type BannerRepository struct {
	db *pgxpool.Pool
}

func NewBannerRepository(db *pgxpool.Pool) *BannerRepository {
	return &BannerRepository{db: db}
}

// TODO:вынести взятие тэгов в отдельную функцию

func (br *BannerRepository) ReadBanner(ctx context.Context, tagId, featureId int) ([]byte, error) {
	var bannerId int
	if err := br.db.QueryRow(ctx, getBannerIdByTagFeature, tagId, featureId).
		Scan(&bannerId); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrBannerNotFound
		}
		return []byte{}, err
	}

	var b []byte
	if err := br.db.QueryRow(ctx, getBannerContentById, bannerId).
		Scan(&b); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrBannerNotFound
		}
		return []byte{}, err
	}

	return b, nil
}

func (br *BannerRepository) ReadFilterBanners(ctx context.Context, tagId, featureId, limit, offset int) ([]models.Banner, error) {
	endOfExp := ""

	if limit != 0 {
		endOfExp = endOfExp + " LIMIT " + strconv.Itoa(limit)
	}

	if offset != 0 {
		endOfExp = endOfExp + " OFFSET " + strconv.Itoa(offset)
	}

	endOfExp = endOfExp + ";"

	var rows pgx.Rows
	var err error
	if tagId != 0 && featureId != 0 {
		if offset != 0 {
			return make([]models.Banner, 0), nil
		}

		rows, err = br.db.Query(ctx, getBannerIdByTagFeature)
	} else if tagId != 0 {
		rows, err = br.db.Query(ctx, getBannerIdsByTag+endOfExp)
	} else if featureId != 0 {
		rows, err = br.db.Query(ctx, getBannerIdsByFeature+endOfExp)
	} else {
		rows, err = br.db.Query(ctx, getAllBannerIds+endOfExp)
	}

	if err != nil {
		return make([]models.Banner, 0), err
	}

	var bannerId int
	banners := make([]models.Banner, 0)
	for rows.Next() {
		err = rows.Scan(&bannerId)
		if err != nil {
			err = fmt.Errorf("error happened in rows.Scan: %w", err)

			return make([]models.Banner, 0), err
		}

		var banner models.Banner
		banner.BannerId = bannerId

		if err := br.db.QueryRow(ctx, getBannerById, bannerId).
			Scan(&banner.Content, &banner.IsActive, &banner.CreatedAt, &banner.UpdatedAt); err != nil {
			fmt.Println("error here", getBannerById, bannerId)
			continue
		}

		if err := br.db.QueryRow(ctx, getFeatureForBanner, bannerId).
			Scan(&banner.FeatureId); err != nil {
			return make([]models.Banner, 0), err
		}

		tagsRows, err := br.db.Query(ctx, getTagsForBanner, bannerId)

		if err != nil {
			return make([]models.Banner, 0), err
		}

		tags := make([]int, 0)
		for tagsRows.Next() {
			var tagId int
			err = tagsRows.Scan(&tagId)

			if err != nil {
				err = fmt.Errorf("error happened in rows.Scan: %w", err)
				return make([]models.Banner, 0), err
			}

			tags = append(tags, tagId)
		}

		banner.TagIds = tags

		banners = append(banners, banner)
	}

	return banners, nil
}

func (br *BannerRepository) CreateBanner(ctx context.Context, banner *models.BannerPayload) (int, error) {
	bannerId := 0
	err := br.db.QueryRow(ctx, createBanner, banner.Content, banner.IsActive).Scan(&bannerId)
	if err != nil {
		return 0, err
	}

	for _, val := range banner.TagIds {
		fmt.Println(val, banner.FeatureId)
		_, err = br.db.Exec(ctx, createFeatureAndTag, bannerId, val, banner.FeatureId)
		if err != nil {
			return 0, err
		}
	}

	return bannerId, nil
}

// TODO:добавить обновление тэгов и фичей

func (br *BannerRepository) UpdateBanner(ctx context.Context, id int, banner *models.BannerPayload) (err error) {
	if !banner.IsActive.HasValue {
		_, err = br.db.Exec(ctx, updateBanner, banner.Content, sql.NullBool{}, id)
	} else {
		_, err = br.db.Exec(ctx, updateBanner, banner.Content, banner.IsActive.IsTrue, id)
	}

	if banner.FeatureId != 0 && banner.TagIds != nil {

	} else if banner.FeatureId != 0 {

	} else if banner.TagIds != nil {

	}

	return err
}

func (br *BannerRepository) DeleteBanner(ctx context.Context, id int) error {
	_, err := br.db.Exec(ctx, deleteBanner, id)
	return err
}
