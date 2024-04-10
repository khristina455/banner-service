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
	var unmarshalledJson bool

	err := json.Unmarshal(b, &unmarshalledJson)
	if err != nil {
		return err
	}

	nullBool.IsTrue = unmarshalledJson
	nullBool.HasValue = true

	return nil
}

type Banner struct {
	BannerId  int             `json:"banner_id"`
	TagIds    []int           `json:"tag_ids"`
	FeatureId int             `json:"feature_id"`
	Content   json.RawMessage `json:"content"`
	IsActive  bool            `json:"is_active"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

type BannerPayload struct {
	TagIds    []int           `json:"tag_ids"`
	FeatureId int             `json:"feature_id"`
	Content   json.RawMessage `json:"content"`
	IsActive  NullBool        `json:"is_active"`
}
