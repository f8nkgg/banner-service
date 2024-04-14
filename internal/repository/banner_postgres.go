package repository

import (
	"banner/internal/entity"
	"banner/pkg/db/postgres"
	"context"
	"errors"
	"time"
)

type BannerRepository struct {
	db *postgres.DB
}

func NewBannerRepository(database *postgres.DB) *BannerRepository {
	return &BannerRepository{
		db: database,
	}
}
func (r *BannerRepository) Save(ctx context.Context, banner *entity.Banner) (int32, error) {
	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return -1, err
	}
	defer tx.Rollback(ctx)

	// Проверяем наличие записи с заданным featureId и tagId
	var count int
	sql, args, err := r.db.Builder.
		Select("COUNT(*)").
		From("banners").
		Where("feature_id = ?", banner.FeatureID).
		Where("tag_ids && ?", banner.TagIDs).
		ToSql()
	if err != nil {
		return -1, err
	}
	err = tx.QueryRow(ctx, sql, args...).Scan(&count)
	if err != nil {
		return -1, err
	}
	if count > 0 {
		return -1, errors.New("record with same featureId and tagId already exists")
	}

	// Если такой записи нет, выполняем вставку
	sql2, args2, err := r.db.Builder.
		Insert("banners").
		Columns("tag_ids", "feature_id", "content", "is_active").
		Values(banner.TagIDs, banner.FeatureID, banner.Content, banner.IsActive).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return -1, err
	}

	var id int32
	err = tx.QueryRow(ctx, sql2, args2...).Scan(&id)
	if err != nil {
		return -1, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return -1, err
	}
	return id, nil
}

func (r *BannerRepository) GetBannerByTagsAndFeatureIDForUser(ctx context.Context, tagID int32, featureID int32, isActiveParam bool) (map[string]interface{}, error) {
	sql, args, err := r.db.Builder.
		Select("content").
		From("banners").
		Where("$1 = ANY(tag_ids) AND feature_id = $2 AND (CASE WHEN $3 THEN true ELSE is_active = true END)", tagID, featureID, isActiveParam).
		Limit(1).
		ToSql()
	if err != nil {
		return nil, err
	}
	var content map[string]interface{}
	err = r.db.Pool.QueryRow(ctx, sql, args...).Scan(&content)
	//защищаем себя от того что текст ошибки может быть изменен
	if content == nil {
		return nil, errors.New("no banner found")
	}
	if err != nil {
		return nil, err
	}
	return content, nil
}
func (r *BannerRepository) GetBannersWithOptionalFilters(ctx context.Context, featureID, tagID, limit, offset *int32) ([]*entity.FilteredBanner, error) {
	sql, args, err := r.db.Builder.
		Select("*").
		From("banners").
		Where("($1::integer IS NULL OR feature_id = $1) AND ($2::integer IS NULL OR $2 = ANY(tag_ids))", featureID, tagID).
		Suffix("LIMIT $3 OFFSET COALESCE($4, 0)", limit, offset).
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var banners []*entity.FilteredBanner
	for rows.Next() {
		var banner entity.FilteredBanner
		if err := rows.Scan(&banner.ID, &banner.TagIDs, &banner.FeatureID, &banner.Content, &banner.IsActive, &banner.CreatedAt, &banner.UpdatedAt); err != nil {
			return nil, err
		}
		banners = append(banners, &banner)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return banners, nil
}
func (r *BannerRepository) DeleteByID(ctx context.Context, id int32) error {
	sql, args, err := r.db.Builder.
		Delete("banners").
		Where("id = $1", id).
		ToSql()
	if err != nil {
		return err
	}

	result, err := r.db.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}
	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("no banner found")
	}
	return nil
}
func (r *BannerRepository) UpdateBanner(ctx context.Context, banner *entity.BannerUpdate) error {
	currentTime := time.Now().UTC()
	updateBuilder := r.db.Builder.Update("banners").Where("id = ?", banner.ID)
	if banner.TagIDs != nil {
		updateBuilder = updateBuilder.Set("tag_ids", banner.TagIDs)
	}
	if banner.FeatureID != nil {
		updateBuilder = updateBuilder.Set("feature_id", banner.FeatureID)
	}
	if banner.Content != nil {
		updateBuilder = updateBuilder.Set("content", banner.Content)
	}
	if banner.IsActive != nil {
		updateBuilder = updateBuilder.Set("is_active", banner.IsActive)
	}
	updateBuilder = updateBuilder.Set("updated_at", currentTime)
	sql, args, err := updateBuilder.ToSql()
	if err != nil {
		return err
	}
	result, err := r.db.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}
	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("no banner found")
	}
	return nil
}
func (r *BannerRepository) GetBannersHistoryByID(ctx context.Context, id int32) ([]*entity.FilteredBanner, error) {
	sql, args, err := r.db.Builder.
		Select("*").
		From("banners_history").
		Where("id = $1", id).
		ToSql()
	if err != nil {
		return nil, err
	}
	rows, err := r.db.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var banners []*entity.FilteredBanner
	for rows.Next() {
		var banner entity.FilteredBanner
		if err := rows.Scan(&banner.ID, &banner.TagIDs, &banner.FeatureID, &banner.Content, &banner.IsActive, &banner.CreatedAt, &banner.UpdatedAt); err != nil {
			return nil, err
		}
		banners = append(banners, &banner)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return banners, nil
}
