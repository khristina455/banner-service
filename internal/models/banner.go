package models

import (
	"encoding/json"
	"time"
)

type Banner struct {
	BannerId  int             `json:"banner_id"`
	TagIds    []int           `json:"tag_ids"`
	FeatureId int             `json:"feature_id"`
	Content   json.RawMessage `json:"content"`
	IsActive  bool            `json:"is_active"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}
