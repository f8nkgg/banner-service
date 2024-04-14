package entity

import "time"

type Banner struct {
	ID        int32                  `json:"id"`
	TagIDs    []int32                `json:"tag_ids"`
	FeatureID int32                  `json:"feature_id"`
	Content   map[string]interface{} `json:"content"`
	IsActive  bool                   `json:"is_active"`
}
type BannerUpdate struct {
	ID        *int32                  `json:"id,omitempty"`
	TagIDs    *[]int32                `json:"tag_ids,omitempty"`
	FeatureID *int32                  `json:"feature_id,omitempty"`
	Content   *map[string]interface{} `json:"content,omitempty"`
	IsActive  *bool                   `json:"is_active,omitempty"`
}
type FilteredBanner struct {
	ID        int32                  `json:"id"`
	TagIDs    []int32                `json:"tag_ids"`
	FeatureID int32                  `json:"feature_id"`
	Content   map[string]interface{} `json:"content"`
	IsActive  bool                   `json:"is_active"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
}
type BannerHistoryItem struct {
	Index  int
	Banner *FilteredBanner
}
