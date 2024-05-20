package service

import (
	"banner/internal/entity"
	"context"
)

type Service interface {
	Save(ctx context.Context, banner *entity.Banner) (int32, error)
	GetForUser(ctx context.Context, tagID, featureID int32, isActiveParam, lastRevision bool) (map[string]interface{}, error)
	GetBanners(ctx context.Context, featureID, tagID, limit *int32, offset int32) ([]*entity.FilteredBanner, error)
	Delete(ctx context.Context, id int32) error
	Update(ctx context.Context, banner *entity.BannerUpdate) error
	GetBannersHistoryByID(ctx context.Context, i int32) ([]*entity.BannerHistoryItem, error)
}
