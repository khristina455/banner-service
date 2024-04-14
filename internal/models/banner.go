package models

import (
	"encoding/json"
	"time"
)

type NullBool struct {
	IsTrue   bool
	HasValue bool
}

func (nullBool *NullBool) UnmarshalJSON(b []byte) error {
	var unmarshalledJSON bool

	err := json.Unmarshal(b, &unmarshalledJSON)
	if err != nil {
		return err
	}

	nullBool.IsTrue = unmarshalledJSON
	nullBool.HasValue = true

	return nil
}

type Banner struct {
	BannerID  int             `json:"banner_id"`
	TagIDs    []int           `json:"tag_ids"`
	FeatureID int             `json:"feature_id"`
	Content   json.RawMessage `json:"content"`
	IsActive  bool            `json:"is_active"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

type BannerPayload struct {
	TagIDs    []int           `json:"tag_ids"`
	FeatureID int             `json:"feature_id"`
	Content   json.RawMessage `json:"content"`
	IsActive  NullBool        `json:"is_active"`
}
