package http

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"banner-service/internal/models"
	"banner-service/internal/pkg/banner"
	"banner-service/internal/pkg/banner/repository"
	"banner-service/internal/utils/responser"
)

type BannerHandler struct {
	service banner.BannerService
	logger  *logrus.Logger
}

func NewBannerHandler(s banner.BannerService, logger *logrus.Logger) *BannerHandler {
	return &BannerHandler{s, logger}
}

func (h *BannerHandler) GetBanner(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("get banner handler")

	tagIDStr := r.URL.Query().Get("tag_id")
	tagID, err := strconv.Atoi(tagIDStr)
	if err != nil {
		h.logger.Error("incorrect tag id ", err)
		responser.WriteError(w, http.StatusBadRequest, errors.New("incorrect tag id"))
		return
	}

	if !r.Context().Value("is_admin").(bool) && tagID != r.Context().Value("tag_id") {
		h.logger.Error("this tag id forbidden")
		responser.WriteStatus(w, http.StatusForbidden)
		return
	}

	featureIDStr := r.URL.Query().Get("feature_id")
	featureID, err := strconv.Atoi(featureIDStr)
	if err != nil {
		h.logger.Error("incorrect feature id ", err)
		responser.WriteError(w, http.StatusBadRequest, errors.New("incorrect feature id"))
		return
	}

	useLastRevisionStr := r.URL.Query().Get("use_last_revision")
	var useLastRevision bool
	if useLastRevisionStr != "" {
		useLastRevision, err = strconv.ParseBool(useLastRevisionStr)
		if err != nil {
			h.logger.Error("incorrect use last revision ", err)
			responser.WriteError(w, http.StatusBadRequest, errors.New("incorrect use last revision"))
			return
		}
	}

	banner, err := h.service.GetBanner(r.Context(), tagID, featureID, useLastRevision,
		r.Context().Value("is_admin").(bool))
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

	tagIDStr := r.URL.Query().Get("tag_id")
	tagID := 0
	if tagIDStr != "" {
		tagID, err = strconv.Atoi(tagIDStr)
		if err != nil {
			h.logger.Error("incorrect tag id ", err, tagIDStr)
			responser.WriteError(w, http.StatusBadRequest, errors.New("incorrect tag id"))
			return
		}
	}

	featureIDStr := r.URL.Query().Get("feature_id")
	featureID := 0
	if tagIDStr != "" {
		featureID, err = strconv.Atoi(featureIDStr)
		if err != nil {
			h.logger.Error("incorrect feature id ", err)
			responser.WriteError(w, http.StatusBadRequest, errors.New("incorrect feature id"))
			return
		}
	}

	limitStr := r.URL.Query().Get("limit")
	limit := 0
	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			h.logger.Error("incorrect limit ", err)
			responser.WriteError(w, http.StatusBadRequest, errors.New("incorrect limit"))
			return
		}
	}

	offsetStr := r.URL.Query().Get("offset")
	offset := 0
	if offsetStr != "" {
		offset, err = strconv.Atoi(offsetStr)
		if err != nil {
			h.logger.Error("incorrect offset ", err)
			responser.WriteError(w, http.StatusBadRequest, errors.New("incorrect offset"))
			return
		}
	}

	banners, err := h.service.GetFilterBanners(r.Context(), tagID, featureID, limit, offset)
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
		responser.WriteError(w, http.StatusBadRequest, errors.New("incorrect data in body request"))
		return
	}
	defer r.Body.Close()

	b := &models.BannerPayload{}
	err = json.Unmarshal(body, b)

	if err != nil {
		h.logger.Error("error in unmarshall")
		responser.WriteError(w, http.StatusBadRequest, errors.New("invalid json in body request"))
		return
	}

	bannerID, err := h.service.AddBanner(r.Context(), b)
	if err != nil {
		responser.WriteError(w, http.StatusInternalServerError, err)
		h.logger.Error(err)
		return
	}

	bannerJSON, err := json.Marshal(struct {
		BannerID int `json:"banner_id"`
	}{BannerID: bannerID})

	responser.WriteJSON(w, http.StatusCreated, bannerJSON)
}

func (h *BannerHandler) UpdateBanner(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("update banner handler")

	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok || idStr == "" {
		h.logger.Error("id is empty")
		responser.WriteError(w, http.StatusBadRequest, errors.New("empty id in request"))
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.logger.Error("id is incorrect ", err)
		responser.WriteError(w, http.StatusBadRequest, errors.New("incorrect id in request"))
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		responser.WriteError(w, http.StatusBadRequest, errors.New("incorrect data in body request"))
		responser.WriteStatus(w, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	b := &models.BannerPayload{}
	err = json.Unmarshal(body, b)

	if err != nil {
		h.logger.Error("error in unmarshall")
		responser.WriteError(w, http.StatusBadRequest, errors.New("invalid json in body request"))
		return
	}

	err = h.service.UpdateBanner(r.Context(), id, b)
	if err != nil {
		h.logger.Error("failed to update banner ", err)
		if errors.Is(err, repository.ErrBannerNotFound) {
			responser.WriteStatus(w, http.StatusNotFound)
			return
		}
		responser.WriteError(w, http.StatusInternalServerError, errors.New("failed to update banner"))
		return
	}

	responser.WriteStatus(w, http.StatusOK)
}

func (h *BannerHandler) DeleteBanner(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("delete banner handler")

	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok || idStr == "" {
		h.logger.Error("id is empty")
		responser.WriteError(w, http.StatusBadRequest, errors.New("empty id in request"))
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.logger.Error("id is incorrect")
		responser.WriteError(w, http.StatusBadRequest, errors.New("incorrect id in  request"))
		return
	}

	err = h.service.DeleteBanner(r.Context(), id)
	if err != nil {
		h.logger.Error("failed to delete banner ", err)
		if errors.Is(err, repository.ErrBannerNotFound) {
			responser.WriteStatus(w, http.StatusNotFound)
			return
		}
		responser.WriteError(w, http.StatusInternalServerError, errors.New("failed to delete banner"))
		return
	}

	responser.WriteStatus(w, http.StatusNoContent)
}
