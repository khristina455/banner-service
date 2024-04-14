package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"banner-service/internal/models"
)

// TODO:сделать транзакции

const (
	getBannerIDByTagFeature    = `SELECT banner_id FROM banner_tag_feature WHERE tag_id=$1 AND feature_id=$2;`
	getBannerIDsByTag          = `SELECT DISTINCT banner_id FROM banner_tag_feature WHERE tag_id=$1`
	getBannerIDsByFeature      = `SELECT banner_id FROM banner_tag_feature WHERE feature_id=$1`
	getAllBannerIDs            = `SELECT banner_id FROM banner_tag_feature`
	getActiveBannerContentByID = `SELECT content FROM banner WHERE banner_id=$1 AND is_active=TRUE;`
	getBannerContentByID       = `SELECT content FROM banner WHERE banner_id=$1;`
	getBannerByID              = `SELECT content, is_active, created_at, updated_at FROM banner WHERE banner_id=$1;`
	getFeatureForBanner        = `SELECT feature_id FROM banner_tag_feature WHERE banner_id=$1;`
	getTagsForBanner           = `SELECT tag_id FROM banner_tag_feature WHERE banner_id=$1;`
	getFeatureTagsForBanner    = `SELECT tag_id, feature_id FROM banner_tag_feature WHERE banner_id=$1;`
	createBanner               = `INSERT INTO banner(content, is_active) VALUES ($1, $2) RETURNING banner_id;`
	createFeatureAndTag        = `INSERT INTO banner_tag_feature(banner_id, tag_id, feature_id) VALUES ($1, $2, $3);`
	updateBanner               = `UPDATE banner SET content = COALESCE($1, content),
                			         is_active= COALESCE($2, is_active), 
                                     updated_at = now() 
					                 WHERE banner_id = $3;`
	updateTagFeatureForBanner = `UPDATE banner_tag_feature SET tag_id = $1,
                			         feature_id = $2 
					                 WHERE tag_id=$3 AND feature_id=$4;`
	updateFeatureForBanner = `UPDATE banner_tag_feature SET
                			         feature_id = $1 
					                 WHERE banner_id=$2;`
	deleteTagFeatureForBanner = `DELETE FROM banner_tag_feature WHERE tag_id=$1 AND feature_id=$2;`
	deleteBanner              = `DELETE FROM banner WHERE banner_id=$1;`
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

func (br *BannerRepository) ReadUserBanner(ctx context.Context, tagID, featureID int) ([]byte, error) {
	var bannerID int
	if err := br.db.QueryRow(ctx, getBannerIDByTagFeature, tagID, featureID).
		Scan(&bannerID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrBannerNotFound
		}
		return []byte{}, err
	}

	var b []byte
	if err := br.db.QueryRow(ctx, getActiveBannerContentByID, bannerID).
		Scan(&b); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrBannerNotFound
		}
		return []byte{}, err
	}

	return b, nil
}

func (br *BannerRepository) ReadBanner(ctx context.Context, tagID, featureID int) ([]byte, error) {
	var bannerID int
	if err := br.db.QueryRow(ctx, getBannerIDByTagFeature, tagID, featureID).
		Scan(&bannerID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrBannerNotFound
		}
		return []byte{}, err
	}

	var b []byte
	if err := br.db.QueryRow(ctx, getBannerContentByID, bannerID).
		Scan(&b); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrBannerNotFound
		}
		return []byte{}, err
	}

	return b, nil
}

func (br *BannerRepository) ReadFilterBanners(ctx context.Context,
	tagID, featureID, limit, offset int) ([]models.Banner, error) {
	endOfExp := ""

	if limit != 0 {
		endOfExp = endOfExp + " LIMIT " + strconv.Itoa(limit)
	}

	if offset != 0 {
		endOfExp = endOfExp + " OFFSET " + strconv.Itoa(offset)
	}

	endOfExp += ";"

	var rows pgx.Rows
	var err error
	if tagID != 0 && featureID != 0 {
		if offset != 0 {
			return make([]models.Banner, 0), nil
		}

		rows, err = br.db.Query(ctx, getBannerIDByTagFeature, tagID, featureID)
	} else if tagID != 0 {
		rows, err = br.db.Query(ctx, getBannerIDsByTag+endOfExp, tagID)
	} else if featureID != 0 {
		rows, err = br.db.Query(ctx, getBannerIDsByFeature+endOfExp, featureID)
	} else {
		rows, err = br.db.Query(ctx, getAllBannerIDs+endOfExp)
	}

	if err != nil {
		return make([]models.Banner, 0), err
	}

	var bannerID int
	banners := make([]models.Banner, 0)
	for rows.Next() {
		err = rows.Scan(&bannerID)
		if err != nil {
			err = fmt.Errorf("error happened in rows.Scan: %w", err)

			return make([]models.Banner, 0), err
		}

		var banner models.Banner
		banner.BannerID = bannerID

		if err = br.db.QueryRow(ctx, getBannerByID, bannerID).
			Scan(&banner.Content, &banner.IsActive, &banner.CreatedAt, &banner.UpdatedAt); err != nil {
			return make([]models.Banner, 0), err
		}

		if err = br.db.QueryRow(ctx, getFeatureForBanner, bannerID).
			Scan(&banner.FeatureID); err != nil {
			return make([]models.Banner, 0), err
		}

		tagsRows, err := br.db.Query(ctx, getTagsForBanner, bannerID)

		if err != nil {
			return make([]models.Banner, 0), err
		}

		tags := make([]int, 0)
		for tagsRows.Next() {
			var tagIDOfRow int
			err = tagsRows.Scan(&tagIDOfRow)

			if err != nil {
				err = fmt.Errorf("error happened in rows.Scan: %w", err)
				return make([]models.Banner, 0), err
			}

			tags = append(tags, tagIDOfRow)
		}

		banner.TagIDs = tags

		banners = append(banners, banner)
	}

	return banners, nil
}

func (br *BannerRepository) CreateBanner(ctx context.Context, banner *models.BannerPayload) (int, error) {
	bannerID := 0
	err := br.db.QueryRow(ctx, createBanner, banner.Content, banner.IsActive.IsTrue).Scan(&bannerID)
	if err != nil {
		return 0, err
	}

	for _, val := range banner.TagIDs {
		_, err = br.db.Exec(ctx, createFeatureAndTag, bannerID, val, banner.FeatureID)
		if err != nil {
			return 0, err
		}
	}

	return bannerID, nil
}

func (br *BannerRepository) UpdateBanner(ctx context.Context, id int, banner *models.BannerPayload) (err error) {
	var cmdTag pgconn.CommandTag
	if !banner.IsActive.HasValue {
		cmdTag, err = br.db.Exec(ctx, updateBanner, banner.Content, sql.NullBool{}, id)
	} else {
		cmdTag, err = br.db.Exec(ctx, updateBanner, banner.Content, banner.IsActive.IsTrue, id)
	}

	if cmdTag.RowsAffected() == 0 {
		return ErrBannerNotFound
	}

	if err != nil {
		return
	}

	if banner.TagIDs != nil {
		rows, err := br.db.Query(ctx, getFeatureTagsForBanner, id)
		if err != nil {
			return
		}

		var tagID, featureID, cnt int
		for rows.Next() {
			err = rows.Scan(&tagID, &featureID)
			if err != nil {
				return
			}
			if banner.FeatureID == 0 {
				banner.FeatureID = featureID
			}
			if cnt < len(banner.TagIDs) {
				_, err = br.db.Exec(ctx, updateTagFeatureForBanner, banner.TagIDs[cnt], banner.FeatureID, tagID, featureID)
				if err != nil {
					return
				}
			} else {
				_, err = br.db.Exec(ctx, deleteTagFeatureForBanner, tagID, featureID)
				if err != nil {
					return
				}
			}
			cnt++
		}

		for cnt < len(banner.TagIDs) {
			_, err = br.db.Exec(ctx, createFeatureAndTag, id, banner.TagIDs[cnt], banner.FeatureID)
			if err != nil {
				return
			}
			cnt++
		}
	} else if banner.FeatureID != 0 {
		_, err = br.db.Exec(ctx, updateFeatureForBanner, banner.FeatureID, id)
	}
	return
}

func (br *BannerRepository) DeleteBanner(ctx context.Context, id int) error {
	cmdTag, err := br.db.Exec(ctx, deleteBanner, id)
	if cmdTag.RowsAffected() == 0 {
		return ErrBannerNotFound
	}

	return err
}
