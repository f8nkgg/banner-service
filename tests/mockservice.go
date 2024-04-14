package tests

import (
	"banner/internal/entity"
	"context"
)

type MockBannerService struct {
	SaveFunc                  func(ctx context.Context, banner *entity.Banner) (int32, error)
	GetForUserFunc            func(ctx context.Context, tagID, featureID int32, isActiveParam, lastRevision bool) (map[string]interface{}, error)
	GetBannersFunc            func(ctx context.Context, featureID, tagID, limit, offset *int32) ([]*entity.FilteredBanner, error)
	DeleteFunc                func(ctx context.Context, id int32) error
	UpdateFunc                func(ctx context.Context, banner *entity.BannerUpdate) error
	GetBannersHistoryByIDFunc func(ctx context.Context, id int32) ([]*entity.BannerHistoryItem, error)
}

func (m *MockBannerService) Save(ctx context.Context, banner *entity.Banner) (int32, error) {
	return m.SaveFunc(ctx, banner)
}

func (m *MockBannerService) GetForUser(ctx context.Context, tagID, featureID int32, isActiveParam, lastRevision bool) (map[string]interface{}, error) {
	return m.GetForUserFunc(ctx, tagID, featureID, isActiveParam, lastRevision)
}

func (m *MockBannerService) GetBanners(ctx context.Context, featureID, tagID, limit, offset *int32) ([]*entity.FilteredBanner, error) {
	return m.GetBannersFunc(ctx, featureID, tagID, limit, offset)
}

func (m *MockBannerService) Delete(ctx context.Context, id int32) error {
	return m.DeleteFunc(ctx, id)
}

func (m *MockBannerService) Update(ctx context.Context, banner *entity.BannerUpdate) error {
	return m.UpdateFunc(ctx, banner)
}

func (m *MockBannerService) GetBannersHistoryByID(ctx context.Context, id int32) ([]*entity.BannerHistoryItem, error) {
	return m.GetBannersHistoryByIDFunc(ctx, id)
}
