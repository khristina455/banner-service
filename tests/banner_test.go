package tests_test

import (
	"banner-service/internal/models"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	bannerHandler "banner-service/internal/pkg/banner/http"
	bannerRepository "banner-service/internal/pkg/banner/repository"
	bannerService "banner-service/internal/pkg/banner/service"
	"banner-service/tests/db"
)

func Test_getUserBanner(t *testing.T) {
	logger := logrus.New()
	formatter := &logrus.TextFormatter{
		TimestampFormat: time.DateTime,
		FullTimestamp:   true,
	}
	logger.SetFormatter(formatter)

	testDB, err := db.Open()
	if err != nil {
		t.Fatalf("error to connect: %v", err)
	}
	defer func() {
		if err := db.Truncate(testDB); err != nil {
			t.Errorf("error truncating test database tables: %v", err)
		}
		testDB.Close()
	}()

	_, err = db.SeedFeatures(testDB)
	if err != nil {
		t.Fatalf("error seeding features: %v", err)
	}

	_, err = db.SeedTags(testDB)
	if err != nil {
		t.Fatalf("error seeding tags: %v", err)
	}

	expectedBanners, err := db.SeedBanners(testDB)
	if err != nil {
		t.Fatalf("error seeding banners: %v", err)
	}

	tests := []struct {
		Name            string
		TagID           int
		FeatureID       int
		IsAdmin         bool
		UserTagID       int
		ExpectedContent []byte
		ExpectedCode    int
	}{
		{
			Name:            "OK for user",
			TagID:           1,
			FeatureID:       1,
			IsAdmin:         false,
			UserTagID:       1,
			ExpectedContent: expectedBanners[0].Content,
			ExpectedCode:    http.StatusOK,
		},
		{
			Name:            "OK for admin",
			TagID:           2,
			FeatureID:       2,
			IsAdmin:         true,
			UserTagID:       0,
			ExpectedContent: expectedBanners[1].Content,
			ExpectedCode:    http.StatusOK,
		},
		{
			Name:            "Forbidden",
			TagID:           3,
			FeatureID:       3,
			IsAdmin:         false,
			UserTagID:       2,
			ExpectedContent: []byte{},
			ExpectedCode:    http.StatusForbidden,
		},
	}

	for _, test := range tests {
		fn := func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/user_banner?tag_id=%d&feature_id=%d",
				test.TagID, test.FeatureID), nil)
			ctx := context.WithValue(req.Context(), "is_admin", test.IsAdmin)
			ctx = context.WithValue(ctx, "tag_id", test.UserTagID)
			req = req.WithContext(ctx)
			if err != nil {
				t.Errorf("error creating request: %v", err)
			}

			w := httptest.NewRecorder()
			br := bannerRepository.NewBannerRepository(testDB)
			bs := bannerService.NewBannerService(br, nil)
			bh := bannerHandler.NewBannerHandler(bs, logger)
			bh.GetBanner(w, req)

			if e, a := test.ExpectedCode, w.Code; e != a {
				t.Errorf("expected status code: %v, got status code: %v", e, a)
			}

			if test.ExpectedContent != nil {
				resp, _ := io.ReadAll(w.Body)
				if d := cmp.Diff(test.ExpectedContent, resp); d != "" {
					t.Errorf("unexpected difference in response body:\n%v", d)
				}
			}
		}

		t.Run(test.Name, fn)
	}
}

func Test_getFilterBanner(t *testing.T) {
	logger := logrus.New()
	formatter := &logrus.TextFormatter{
		TimestampFormat: time.DateTime,
		FullTimestamp:   true,
	}
	logger.SetFormatter(formatter)

	testDB, err := db.Open()
	if err != nil {
		t.Fatalf("error to connect: %v", err)
	}
	defer func() {
		if err := db.Truncate(testDB); err != nil {
			t.Errorf("error truncating test database tables: %v", err)
		}
		testDB.Close()
	}()

	features, err := db.SeedFeatures(testDB)
	if err != nil {
		t.Fatalf("error seeding features: %v", err)
	}

	_, err = db.SeedTags(testDB)
	if err != nil {
		t.Fatalf("error seeding tags: %v", err)
	}

	expectedBanners, err := db.SeedBanners(testDB)
	if err != nil {
		t.Fatalf("error seeding banners: %v", err)
	}

	tests := []struct {
		Name         string
		TagID        int
		FeatureID    int
		Limit        int
		Offset       int
		ExpectedResp []models.Banner
		ExpectedCode int
	}{
		{
			Name:         "Get Banners by feature",
			TagID:        0,
			FeatureID:    features[0],
			Limit:        0,
			Offset:       0,
			ExpectedResp: []models.Banner{expectedBanners[0]},
			ExpectedCode: http.StatusOK,
		},
	}

	for _, test := range tests {
		fn := func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/banner?tag_id=%d&feature_id=%d",
				test.TagID, test.FeatureID), nil)

			if err != nil {
				t.Errorf("error creating request: %v", err)
			}

			w := httptest.NewRecorder()
			br := bannerRepository.NewBannerRepository(testDB)
			bs := bannerService.NewBannerService(br, nil)
			bh := bannerHandler.NewBannerHandler(bs, logger)
			bh.GetBannerList(w, req)

			if e, a := test.ExpectedCode, w.Code; e != a {
				t.Errorf("expected status code: %v, got status code: %v", e, a)
			}

			if test.ExpectedResp != nil {
				var resp []models.Banner
				if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
					t.Errorf("error decoding response body: %v", err)
				}
				if d := cmp.Diff(test.ExpectedResp, resp); d != "" {
					t.Errorf("unexpected difference in response body:\n%v", d)
				}
			}
		}

		t.Run(test.Name, fn)
	}
}

func Test_addBanner(t *testing.T) {
	logger := logrus.New()
	formatter := &logrus.TextFormatter{
		TimestampFormat: time.DateTime,
		FullTimestamp:   true,
	}
	logger.SetFormatter(formatter)

	testDB, err := db.Open()
	if err != nil {
		t.Fatalf("error to connect: %v", err)
	}
	defer func() {
		if err := db.Truncate(testDB); err != nil {
			t.Errorf("error truncating test database tables: %v", err)
		}
		testDB.Close()
	}()

	features, err := db.SeedFeatures(testDB)
	if err != nil {
		t.Fatalf("error seeding features: %v", err)
	}

	tags, err := db.SeedTags(testDB)
	if err != nil {
		t.Fatalf("error seeding tags: %v", err)
	}

	_, err = db.SeedBanners(testDB)
	if err != nil {
		t.Fatalf("error seeding banners: %v", err)
	}

	tests := []struct {
		Name        string
		RequestBody models.Banner
		Resp        struct {
			BannerID int `json:"banner_id"`
		}
		ExpectedCode int
	}{
		{
			Name: "OK for user",
			RequestBody: models.Banner{
				TagIDs:    []int{tags[5], tags[6]},
				FeatureID: features[9],
				Content:   []byte(`{"text":"awesome banner"}`),
				IsActive:  true,
			},
			ExpectedCode: http.StatusCreated,
		},
	}

	for _, test := range tests {
		fn := func(t *testing.T) {
			var b bytes.Buffer
			if err := json.NewEncoder(&b).Encode(test.RequestBody); err != nil {
				t.Errorf("error encoding request body: %v", err)
			}

			req, err := http.NewRequest(http.MethodPost, "/banner", &b)
			if err != nil {
				t.Errorf("error creating request: %v", err)
			}

			w := httptest.NewRecorder()
			br := bannerRepository.NewBannerRepository(testDB)
			bs := bannerService.NewBannerService(br, nil)
			bh := bannerHandler.NewBannerHandler(bs, logger)
			bh.AddBanner(w, req)

			if e, a := test.ExpectedCode, w.Code; e != a {
				t.Errorf("expected status code: %v, got status code: %v", e, a)
			}

			if err := json.NewDecoder(w.Body).Decode(&test.Resp); err != nil {
				t.Errorf("error decoding response body: %v", err)
			}

			var expectedBanner models.Banner

			err = testDB.QueryRow(context.Background(), `SELECT content, is_active FROM banner WHERE banner_id=$1`,
				test.Resp.BannerID).Scan(&expectedBanner.Content, &expectedBanner.IsActive)

			if err != nil {
				t.Errorf("error to get banner")
			}

			if d := cmp.Diff(test.RequestBody.Content, expectedBanner.Content); d != "" {
				t.Errorf("unexpected difference in response body:\n%v", d)
			}

			if e, a := test.RequestBody.IsActive, expectedBanner.IsActive; e != a {
				t.Errorf("expected banner activity: %v, got banner activity: %v", e, a)
			}
		}

		t.Run(test.Name, fn)
	}
}

func Test_updateBanner(t *testing.T) {
	logger := logrus.New()
	formatter := &logrus.TextFormatter{
		TimestampFormat: time.DateTime,
		FullTimestamp:   true,
	}
	logger.SetFormatter(formatter)

	testDB, err := db.Open()
	if err != nil {
		t.Fatalf("error to connect: %v", err)
	}
	defer func() {
		if err := db.Truncate(testDB); err != nil {
			t.Errorf("error truncating test database tables: %v", err)
		}
		testDB.Close()
	}()

	_, err = db.SeedFeatures(testDB)
	if err != nil {
		t.Fatalf("error seeding features: %v", err)
	}

	_, err = db.SeedTags(testDB)
	if err != nil {
		t.Fatalf("error seeding tags: %v", err)
	}

	banners, err := db.SeedBanners(testDB)
	if err != nil {
		t.Fatalf("error seeding lists: %v", err)
	}

	tests := []struct {
		Name         string
		BannerId     int
		RequestBody  models.Banner
		ExpectedCode int
	}{
		{
			Name:     "OK for user",
			BannerId: banners[0].BannerID,
			RequestBody: models.Banner{
				Content: []byte(`{"text":"new text","url":"new url"}`),
			},
			ExpectedCode: http.StatusOK,
		},
	}

	for _, test := range tests {
		fn := func(t *testing.T) {
			var b bytes.Buffer
			if err := json.NewEncoder(&b).Encode(test.RequestBody); err != nil {
				t.Errorf("error encoding request body: %v", err)
			}

			req, err := http.NewRequest(http.MethodPatch, fmt.Sprintf("/banner/%v", test.BannerId), &b)
			if err != nil {
				t.Errorf("error creating request: %v", err)
			}

			req = mux.SetURLVars(req, map[string]string{"id": strconv.Itoa(test.BannerId)})

			w := httptest.NewRecorder()
			br := bannerRepository.NewBannerRepository(testDB)
			bs := bannerService.NewBannerService(br, nil)
			bh := bannerHandler.NewBannerHandler(bs, logger)
			bh.UpdateBanner(w, req)

			if e, a := test.ExpectedCode, w.Code; e != a {
				t.Errorf("expected status code: %v, got status code: %v", e, a)
			}

			var banner models.Banner

			err = testDB.QueryRow(context.Background(), `SELECT content, is_active FROM banner WHERE banner_id=$1`,
				test.BannerId).Scan(&banner.Content, &banner.IsActive)

			if err != nil {
				t.Errorf("error to get banner")
			}

			if d := cmp.Diff(test.RequestBody.Content, banner.Content); d != "" {
				t.Errorf("unexpected difference in response body:\n%v", d)
			}
		}
		t.Run(test.Name, fn)
	}
}

func Test_deleteBanner(t *testing.T) {
	logger := logrus.New()
	formatter := &logrus.TextFormatter{
		TimestampFormat: time.DateTime,
		FullTimestamp:   true,
	}
	logger.SetFormatter(formatter)

	testDB, err := db.Open()
	if err != nil {
		t.Fatalf("error to connect: %v", err)
	}
	defer func() {
		if err := db.Truncate(testDB); err != nil {
			t.Errorf("error truncating test database tables: %v", err)
		}
		testDB.Close()
	}()

	_, err = db.SeedFeatures(testDB)
	if err != nil {
		t.Fatalf("error seeding features: %v", err)
	}

	_, err = db.SeedTags(testDB)
	if err != nil {
		t.Fatalf("error seeding tags: %v", err)
	}

	banners, err := db.SeedBanners(testDB)
	if err != nil {
		t.Fatalf("error seeding lists: %v", err)
	}

	tests := []struct {
		Name         string
		BannerId     int
		ExpectedCode int
	}{
		{
			Name:         "OK for user",
			BannerId:     banners[0].BannerID,
			ExpectedCode: http.StatusNoContent,
		},
	}

	for _, test := range tests {
		fn := func(t *testing.T) {
			req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/banner/%v", test.BannerId), nil)
			if err != nil {
				t.Errorf("error creating request: %v", err)
			}

			req = mux.SetURLVars(req, map[string]string{"id": strconv.Itoa(test.BannerId)})

			w := httptest.NewRecorder()
			br := bannerRepository.NewBannerRepository(testDB)
			bs := bannerService.NewBannerService(br, nil)
			bh := bannerHandler.NewBannerHandler(bs, logger)
			bh.DeleteBanner(w, req)

			if e, a := test.ExpectedCode, w.Code; e != a {
				t.Errorf("expected status code: %v, got status code: %v", e, a)
			}

			cmdTag, err := testDB.Exec(context.Background(), `SELECT * FROM banner WHERE banner_id=$1`, test.BannerId)
			if err != nil {
				t.Errorf("error to get banner")
			}

			if cmdTag.RowsAffected() != 0 {
				t.Errorf("banner still in table")
			}
		}

		t.Run(test.Name, fn)
	}
}
