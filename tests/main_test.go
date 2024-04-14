package tests

import (
	v1 "banner/internal/controller/http/v1"
	"banner/internal/repository"
	"banner/internal/service"
	"banner/pkg/cache"
	"banner/pkg/db/postgres"
	"banner/pkg/logger"
	"context"
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
	"time"
)

var pgURL string

func init() {
	pgURL = os.Getenv("TEST_DB_URL")
	time.Sleep(5 * time.Second)
}

type APITestSuite struct {
	suite.Suite

	db      *postgres.DB
	handler *v1.BannerController
	service *service.BannerService
	repo    *repository.BannerRepository
	logger  logger.Logger
}

func TestAPISuite(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	suite.Run(t, new(APITestSuite))
}

func (s *APITestSuite) SetupSuite() {
	s.logger = logger.New("debug")
	db, err := postgres.New(pgURL, 5, 5, 100)
	if err != nil {
		s.FailNow("Failed to connect to PostgreSQL", err)
	}
	s.db = db
	s.initialize()
	if err := s.createTable(); err != nil {
		s.FailNow("Failed to create table", err)
	}
}
func (s *APITestSuite) TearDownSuite() {
	_, err := s.db.Pool.Exec(context.Background(), `
	      DROP TABLE banners;
DROP TABLE banners_history;`)
	if err != nil {
		s.FailNow("Failed to drop table", err)
	}
	s.db.Close()
}
func (s *APITestSuite) initialize() {
	repo := repository.NewBannerRepository(s.db)
	memCache := cache.NewMemoryCache(1000, 20)
	serv := service.NewBannerService(repo, memCache, 5*time.Minute)
	contr := v1.NewBannerController(serv, s.logger)
	s.repo = repo
	s.service = serv
	s.handler = contr
}
func TestMain(m *testing.M) {
	rc := m.Run()
	os.Exit(rc)
}
func (s *APITestSuite) createTestBanner() {
	sqlQuery, args, err := s.db.Builder.Insert("banners").SetMap(map[string]interface{}{
		"id":         1,
		"tag_ids":    []int32{4, 5, 6},
		"feature_id": 123,
		"content": map[string]interface{}{
			"text":  "some_text3",
			"title": "some_title",
			"url":   "some_url2",
		},
		"is_active": true,
	}).ToSql()
	s.NoError(err)
	_, err = s.db.Pool.Exec(context.Background(), sqlQuery, args...)
	s.NoError(err)
}
func (s *APITestSuite) deleteTestBanner() {
	sqlQuery, args, err := s.db.Builder.Delete("banners").Where("id = ?", 1).ToSql()
	s.NoError(err)
	_, err = s.db.Pool.Exec(context.Background(), sqlQuery, args...)
	s.NoError(err)

}
func (s *APITestSuite) createTable() error {
	_, err := s.db.Pool.Exec(context.Background(), `
        CREATE TABLE IF NOT EXISTS banners (
                         id SERIAL PRIMARY KEY,
                         tag_ids integer[],
                         feature_id integer,
                         content jsonb,
                         is_active boolean,
                         created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                         updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_tag_ids ON banners USING GIN (tag_ids);
CREATE INDEX IF NOT EXISTS idx_feature_id ON banners (feature_id);
CREATE INDEX IF NOT EXISTS idx_is_active ON banners (id) WHERE is_active = true;

CREATE TABLE IF NOT EXISTS banners_history (
                                               id integer,
                                               tag_ids integer[],
                                               feature_id integer,
                                               content jsonb,
                                               is_active boolean,
                                               created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                                               updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                                               PRIMARY KEY (id, updated_at)
);
CREATE OR REPLACE FUNCTION save_banner_history()
    RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO banners_history (id, tag_ids, feature_id, content, is_active, created_at)
    VALUES (OLD.id, OLD.tag_ids, OLD.feature_id, OLD.content, OLD.is_active, OLD.created_at);

    DELETE FROM banners_history
    WHERE (id, updated_at) NOT IN (
        SELECT id, updated_at
        FROM banners_history
        ORDER BY updated_at DESC
        LIMIT 3
    );

    RETURN OLD;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER banners_history_trigger
    AFTER UPDATE ON banners
    FOR EACH ROW EXECUTE FUNCTION save_banner_history();
    `)
	if err != nil {
		return err
	}
	return nil
}
