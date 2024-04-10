package http

import (
	"banner-service/internal/models"
	"banner-service/internal/pkg/banner"
	"banner-service/internal/pkg/banner/repository"
	"banner-service/internal/utils/responser"
	"encoding/json"
	"errors"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"strconv"
)

type BannerHandler struct {
	service banner.BannerService
	logger  *logrus.Logger
}

func NewBannerHandler(s banner.BannerService, logger *logrus.Logger) *BannerHandler {
	return &BannerHandler{s, logger}
}

// TODO: поправить обработку ошибок

func (h *BannerHandler) GetBanner(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("get banner")
	tagIdStr := r.URL.Query().Get("tag_id")
	tagId, err := strconv.Atoi(tagIdStr)
	if err != nil {
		h.logger.Error("tag id ", err)
		responser.WriteError(w, http.StatusBadRequest, err)
		return
	}

	featureIdStr := r.URL.Query().Get("feature_id")
	featureId, err := strconv.Atoi(featureIdStr)
	if err != nil {
		h.logger.Error("feature id ", err)
		responser.WriteError(w, http.StatusBadRequest, err)
		return
	}

	useLastRevisionStr := r.URL.Query().Get("use_last_revision")
	var useLastRevision bool
	if useLastRevisionStr != "" {
		useLastRevision, err = strconv.ParseBool(useLastRevisionStr)
		if err != nil {
			h.logger.Error("use_last_revision ", err)
			responser.WriteError(w, http.StatusBadRequest, err)
			return
		}
	}

	banner, err := h.service.GetUserBanner(r.Context(), tagId, featureId, useLastRevision)
	if err != nil {
		h.logger.Error("failed to get banner", err)
		if errors.Is(err, repository.ErrBannerNotFound) {
			responser.WriteStatus(w, http.StatusNotFound)
			return
		}
		responser.WriteError(w, http.StatusBadRequest, err)
		return
	}

	responser.WriteJSON(w, http.StatusOK, banner)
}

func (h *BannerHandler) GetBannerList(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("get banners")
	var err error
	tagIdStr := r.URL.Query().Get("tag_id")
	tagId := 0
	if tagIdStr != "" {
		tagId, err = strconv.Atoi(tagIdStr)
		if err != nil {
			h.logger.Error("tag id ", err)
			responser.WriteError(w, http.StatusBadRequest, err)
			return
		}
	}

	featureIdStr := r.URL.Query().Get("feature_id")
	featureId := 0
	if tagIdStr != "" {
		featureId, err = strconv.Atoi(featureIdStr)
		if err != nil {
			h.logger.Error("feature id ", err)
			responser.WriteError(w, http.StatusBadRequest, err)
			return
		}
	}

	limitStr := r.URL.Query().Get("limit")
	limit := 0
	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			h.logger.Error("limit ", err)
			responser.WriteError(w, http.StatusBadRequest, err)
			return
		}
	}

	offsetStr := r.URL.Query().Get("offset")
	offset := 0
	if tagIdStr != "" {
		offset, err = strconv.Atoi(offsetStr)
		if err != nil {
			h.logger.Error("offset ", err)
			responser.WriteError(w, http.StatusBadRequest, err)
			return
		}
	}

	banners, err := h.service.GetFilterBanners(r.Context(), tagId, featureId, limit, offset)
	if err != nil {
		h.logger.Error("failed to get banners", err)
		responser.WriteError(w, http.StatusBadRequest, err)
		return
	}

	bannersJSON, err := json.Marshal(banners)
	if err != nil {
		h.logger.Error("failed to get banners", err)
		responser.WriteError(w, http.StatusBadRequest, err)
		return
	}

	responser.WriteJSON(w, http.StatusOK, bannersJSON)
}

func (h *BannerHandler) AddBanner(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		responser.WriteStatus(w, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	b := &models.Banner{}
	err = json.Unmarshal(body, b)

	bannerId, err := h.service.AddBanner(r.Context(), b)
	if err != nil {
		responser.WriteStatus(w, http.StatusInternalServerError)
		return
	}

	bannerJSON, err := json.Marshal(struct {
		BannerID int `json:"banner_id"`
	}{BannerID: bannerId})

	responser.WriteJSON(w, http.StatusCreated, bannerJSON)
}

func (h *BannerHandler) UpdateBanner(w http.ResponseWriter, r *http.Request) {

}

func (h *BannerHandler) DeleteBanner(w http.ResponseWriter, r *http.Request) {

}
