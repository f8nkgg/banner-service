package service

import (
	"banner/internal/entity"
	"banner/internal/repository"
	"banner/pkg/cache"
	"context"
	"time"
)

type BannerService struct {
	bannerRepository *repository.BannerRepository
	cache            *cache.MemoryCache
	cacheTTL         time.Duration
}

func NewBannerService(bannerRepository *repository.BannerRepository, memoryCache *cache.MemoryCache, cacheTTL time.Duration) *BannerService {
	return &BannerService{bannerRepository: bannerRepository, cache: memoryCache, cacheTTL: cacheTTL}
}
func (s *BannerService) Save(ctx context.Context, banner *entity.Banner) (int32, error) {
	return s.bannerRepository.Save(ctx, banner)
}
func (s *BannerService) GetForUser(ctx context.Context, tagID, featureID int32, isActiveParam, lastRevision bool) (map[string]interface{}, error) {
	if !lastRevision {
		if value, err := s.cache.Get(tagID, featureID); err == nil {
			return value, nil
		}
	}
	content, err := s.bannerRepository.GetBannerByTagsAndFeatureIDForUser(ctx, tagID, featureID, isActiveParam)
	if err != nil {
		return content, err
	}
	s.cache.Set(tagID, featureID, content, s.cacheTTL)
	return content, err
}
func (s *BannerService) GetBanners(ctx context.Context, featureID, tagID, limit *int32, offset int32) ([]*entity.FilteredBanner, error) {
	return s.bannerRepository.GetBannersWithOptionalFilters(ctx, featureID, tagID, limit, offset)
}
func (s *BannerService) Delete(ctx context.Context, id int32) error {
	return s.bannerRepository.DeleteByID(ctx, id)
}
func (s *BannerService) Update(ctx context.Context, banner *entity.BannerUpdate) error {
	return s.bannerRepository.UpdateBanner(ctx, banner)
}
func (s *BannerService) GetBannersHistoryByID(ctx context.Context, id int32) ([]*entity.BannerHistoryItem, error) {
	banners, err := s.bannerRepository.GetBannersHistoryByID(ctx, id)
	if err != nil {
		return nil, err
	}
	var bannerHistory []*entity.BannerHistoryItem
	for i, banner := range banners {
		bannerHistory = append(bannerHistory, &entity.BannerHistoryItem{
			Index:  i + 1,
			Banner: banner,
		})
	}
	return bannerHistory, nil
}
