package http

import (
	"banner-service/internal/models"
	"banner-service/internal/pkg/banner"
	"banner-service/internal/pkg/banner/repository"
	"banner-service/internal/utils/responser"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
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

// TODO: сделать нормальную обработку ошибок в соответсвии с api и нормальное вывод логгера

func (h *BannerHandler) GetBanner(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("get banner handler")

	tagIdStr := r.URL.Query().Get("tag_id")
	tagId, err := strconv.Atoi(tagIdStr)
	if err != nil {
		h.logger.Error("incorrect tag id ", err)
		responser.WriteError(w, http.StatusBadRequest, fmt.Errorf("incorrect tag id"))
		return
	}

	if !r.Context().Value("is_admin").(bool) && tagId != r.Context().Value("tag_id") {
		h.logger.Error("this tag id forbidden")
		responser.WriteStatus(w, http.StatusForbidden)
		return
	}

	featureIdStr := r.URL.Query().Get("feature_id")
	featureId, err := strconv.Atoi(featureIdStr)
	if err != nil {
		h.logger.Error("incorrect feature id ", err)
		responser.WriteError(w, http.StatusBadRequest, fmt.Errorf("incorrect feature id"))
		return
	}

	useLastRevisionStr := r.URL.Query().Get("use_last_revision")
	var useLastRevision bool
	if useLastRevisionStr != "" {
		useLastRevision, err = strconv.ParseBool(useLastRevisionStr)
		if err != nil {
			h.logger.Error("incorrect use last revision ", err)
			responser.WriteError(w, http.StatusBadRequest, fmt.Errorf("incorrect use last revision"))
			return
		}
	}

	banner, err := h.service.GetBanner(r.Context(), tagId, featureId, useLastRevision, r.Context().Value("is_admin").(bool))
	if err != nil {
		h.logger.Error("failed to get banner ", err)
		if errors.Is(err, repository.ErrBannerNotFound) {
			responser.WriteStatus(w, http.StatusNotFound)
			return
		}
		responser.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	responser.WriteJSON(w, http.StatusOK, banner)
}

func (h *BannerHandler) GetBannerList(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("get banners handler")

	var err error

	tagIdStr := r.URL.Query().Get("tag_id")
	tagId := 0
	if tagIdStr != "" {
		tagId, err = strconv.Atoi(tagIdStr)
		if err != nil {
			h.logger.Error("incorrect tag id ", err, tagIdStr)
			responser.WriteError(w, http.StatusBadRequest, fmt.Errorf("incorrect tag id"))
			return
		}
	}

	featureIdStr := r.URL.Query().Get("feature_id")
	featureId := 0
	if tagIdStr != "" {
		featureId, err = strconv.Atoi(featureIdStr)
		if err != nil {
			h.logger.Error("incorrect feature id ", err)
			responser.WriteError(w, http.StatusBadRequest, fmt.Errorf("incorrect feature id"))
			return
		}
	}

	limitStr := r.URL.Query().Get("limit")
	limit := 0
	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			h.logger.Error("incorrect limit ", err)
			responser.WriteError(w, http.StatusBadRequest, fmt.Errorf("incorrect limit"))
			return
		}
	}

	offsetStr := r.URL.Query().Get("offset")
	offset := 0
	if offsetStr != "" {
		offset, err = strconv.Atoi(offsetStr)
		if err != nil {
			h.logger.Error("incorrect offset ", err)
			responser.WriteError(w, http.StatusBadRequest, fmt.Errorf("incorrect offset"))
			return
		}
	}

	h.logger.Info(tagId, featureId, limit, offset)

	banners, err := h.service.GetFilterBanners(r.Context(), tagId, featureId, limit, offset)
	if err != nil {
		h.logger.Error("failed to get banners ", err)
		responser.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	bannersJSON, err := json.Marshal(banners)
	if err != nil {
		h.logger.Error("failed to get banners ", err)
		responser.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	responser.WriteJSON(w, http.StatusOK, bannersJSON)
}

func (h *BannerHandler) AddBanner(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("add banner handler")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		responser.WriteStatus(w, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	b := &models.BannerPayload{}
	err = json.Unmarshal(body, b)

	if err != nil {
		h.logger.Error("error in unmarshall")
		responser.WriteStatus(w, http.StatusBadRequest)
		return
	}

	bannerId, err := h.service.AddBanner(r.Context(), b)
	if err != nil {
		responser.WriteStatus(w, http.StatusInternalServerError)
		h.logger.Error(err)
		return
	}

	bannerJSON, err := json.Marshal(struct {
		BannerID int `json:"banner_id"`
	}{BannerID: bannerId})

	responser.WriteJSON(w, http.StatusCreated, bannerJSON)
}

func (h *BannerHandler) UpdateBanner(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok || idStr == "" {
		h.logger.Error("id is empty")
		responser.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid request"))
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.logger.Error("id is invalid", err)
		responser.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid request"))
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		responser.WriteStatus(w, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	b := &models.BannerPayload{}
	err = json.Unmarshal(body, b)

	if err != nil {
		h.logger.Error("error in unmarshall")
		responser.WriteStatus(w, http.StatusBadRequest)
		return
	}

	err = h.service.UpdateBanner(r.Context(), id, b)
	if err != nil {
		responser.WriteStatus(w, http.StatusInternalServerError)
		h.logger.Error(err)
		return
	}

	responser.WriteStatus(w, http.StatusOK)
}

func (h *BannerHandler) DeleteBanner(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok || idStr == "" {
		h.logger.Error("id is empty")
		responser.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid request"))
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.logger.Error("id is invalid", err)
		responser.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid request"))
		return
	}

	err = h.service.DeleteBanner(r.Context(), id)
	if err != nil {

	}

	responser.WriteStatus(w, http.StatusNoContent)
}
