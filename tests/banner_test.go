package tests

import (
	"banner-service/internal/models"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
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

	expectedBanners, err := db.SeedBanners(testDB)
	if err != nil {
		t.Fatalf("error seeding lists: %v", err)
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

	expectedBanners, err := db.SeedBanners(testDB)
	if err != nil {
		t.Fatalf("error seeding lists: %v", err)
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

	expectedBanners, err := db.SeedBanners(testDB)
	if err != nil {
		t.Fatalf("error seeding lists: %v", err)
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

	expectedBanners, err := db.SeedBanners(testDB)
	if err != nil {
		t.Fatalf("error seeding lists: %v", err)
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

	expectedBanners, err := db.SeedBanners(testDB)
	if err != nil {
		t.Fatalf("error seeding lists: %v", err)
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
