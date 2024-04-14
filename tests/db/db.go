package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"banner-service/internal/models"
)

const (
	databaseUser = "test_user"
	databasePass = "1234"
	databaseName = "test_db"
	databaseHost = "postgres"
	databasePort = 5432
)

const schema = `
CREATE TABLE IF NOT EXISTS banner(
    banner_id  SERIAL PRIMARY KEY,
    content    BYTEA NOT NULL,
    is_active  BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS tag(
    tag_id  INT PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS feature(
    feature_id  INT PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS banner_tag_feature(
    banner_id INT,
    tag_id INT,
    feature_id INT,
    FOREIGN KEY (banner_id) REFERENCES banner(banner_id) ON DELETE CASCADE,
    FOREIGN KEY (tag_id) REFERENCES tag(tag_id) ON DELETE CASCADE,
    FOREIGN KEY (feature_id) REFERENCES feature(feature_id) ON DELETE CASCADE,
    CONSTRAINT PK_TagFeature PRIMARY KEY (tag_id, feature_id)
);
`

type Config struct {
	User string
	Pass string
	Name string
	Host string
	Port int
}

// NewConnection returns a new database connection with the schema applied, if not already
// applied.
func NewConnection(cfg Config) (*pgxpool.Pool, error) {
	db, err := pgxpool.New(context.Background(), fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable",
		cfg.User,
		cfg.Pass,
		cfg.Host,
		cfg.Port,
		cfg.Name))

	if err != nil {
		err = errors.New("error happened in sql.Open: " + err.Error())
		return nil, err
	}

	if err = db.Ping(context.Background()); err != nil {
		return nil, err
	}

	if _, err = db.Exec(context.Background(), schema); err != nil {
		return nil, errors.New("apply database schema")
	}

	return db, nil
}

// Open returns a new database connection for the test database.
func Open() (*pgxpool.Pool, error) {
	return NewConnection(Config{
		User: databaseUser,
		Pass: databasePass,
		Name: databaseName,
		Host: databaseHost,
		Port: databasePort,
	})
}

func Truncate(dbc *pgxpool.Pool) error {
	stmt := `TRUNCATE TABLE banner_tag_feature, banner;`

	if _, err := dbc.Exec(context.Background(), stmt); err != nil {
		return errors.New("truncate test database tables")
	}

	return nil
}

func SeedTags(dbc *pgxpool.Pool) ([]int, error) {
	tags := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	for _, i := range tags {
		_, err := dbc.Exec(context.Background(), `INSERT INTO tag(tag_id) VALUES ($1);`, i)
		if err != nil {
			return nil, errors.New("prepare list insertion")
		}
	}

	return tags, nil
}

func SeedFeatures(dbc *pgxpool.Pool) ([]int, error) {
	features := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	for _, i := range features {
		_, err := dbc.Exec(context.Background(), `INSERT INTO features(features_id) VALUES ($1);`, i)
		if err != nil {
			return nil, errors.New("prepare list insertion")
		}
	}

	return features, nil
}

func SeedBanners(dbc *pgxpool.Pool) ([]models.Banner, error) {
	banners := []models.Banner{
		{
			Content:   []byte(`{"content": "content of banner 1"}`),
			IsActive:  true,
			TagIDs:    []int{1},
			FeatureID: 1,
		},
		{
			Content:   []byte(`{"content": "content of banner 2"}`),
			IsActive:  true,
			TagIDs:    []int{2},
			FeatureID: 2,
		},
		{
			Content:   []byte(`{"content": "content of banner 3"}`),
			IsActive:  true,
			TagIDs:    []int{3},
			FeatureID: 3,
		},
		{
			Content:   []byte(`{"content": "content of banner 4"}`),
			IsActive:  true,
			TagIDs:    []int{4},
			FeatureID: 4,
		},
		{
			Content:   []byte(`{"content": "content of banner 5"}`),
			IsActive:  true,
			TagIDs:    []int{5},
			FeatureID: 5,
		},
		{
			Content:   []byte(`{"content": "content of banner 6"}`),
			IsActive:  true,
			TagIDs:    []int{6},
			FeatureID: 6,
		},
	}

	for i := range banners {
		err := dbc.QueryRow(context.Background(),
			`INSERT INTO banner(content, is_active) VALUES ($1, $2) RETURNING banner_id;`,
			banners[i].Content, banners[i].IsActive).Scan(&banners[i].BannerID)
		if err != nil {
			return nil, errors.New("prepare list insertion")
		}

		_, err = dbc.Exec(context.Background(),
			`INSERT INTO banner_tag_feature(banner_id, tag_id, feature_id) VALUES ($1, $2, $3) RETURNING banner_id;`,
			banners[i].BannerID, banners[i].TagIDs[0], banners[i].FeatureID)
		if err != nil {
			return nil, errors.New("prepare list insertion")
		}
	}

	return banners, nil
}
