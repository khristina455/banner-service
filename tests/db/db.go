package db

import (
	"banner-service/internal/models"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	// databaseUser is the user for the test database.
	databaseUser = "test_user"

	// databasePass is the password of the user for the test database.
	databasePass = "1234"

	// databaseName is the name of the test database.
	databaseName = "test_db"

	// databaseHost is the host name of the test database.
	databaseHost = "postgres"

	// databasePort is the port that the test database is listening on.
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
CREATE TABLE IF NOT EXISTS banner_tag_feature(
    banner_id INT,
    tag_id INT,
    feature_id INT,
    FOREIGN KEY (banner_id) REFERENCES banner(banner_id) ON DELETE CASCADE,
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
		err = fmt.Errorf("error happened in sql.Open: %w", err)
		return nil, err
	}

	if err = db.Ping(context.Background()); err != nil {
		return nil, err
	}

	if _, err = db.Exec(context.Background(), schema); err != nil {
		return nil, fmt.Errorf("apply database schema")
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
		return fmt.Errorf("truncate test database tables")
	}

	return nil
}

func SeedBanners(dbc *pgxpool.Pool) ([]models.Banner, error) {
	banners := []models.Banner{
		{
			Content:   []byte(`{"content": "content of banner 1"}`),
			IsActive:  true,
			TagIds:    []int{1},
			FeatureId: 1,
		},
		{
			Content:   []byte(`{"content": "content of banner 2"}`),
			IsActive:  true,
			TagIds:    []int{2},
			FeatureId: 2,
		},
		{
			Content:   []byte(`{"content": "content of banner 3"}`),
			IsActive:  true,
			TagIds:    []int{3},
			FeatureId: 3,
		},
		{
			Content:   []byte(`{"content": "content of banner 4"}`),
			IsActive:  true,
			TagIds:    []int{4},
			FeatureId: 4,
		},
		{
			Content:   []byte(`{"content": "content of banner 5"}`),
			IsActive:  true,
			TagIds:    []int{5},
			FeatureId: 5,
		},
		{
			Content:   []byte(`{"content": "content of banner 6"}`),
			IsActive:  true,
			TagIds:    []int{6},
			FeatureId: 6,
		},
	}

	for i := range banners {
		err := dbc.QueryRow(context.Background(), `INSERT INTO banner(content, is_active) VALUES ($1, $2) RETURNING banner_id;`, banners[i].Content, banners[i].IsActive).Scan(&banners[i].BannerId)
		if err != nil {
			return nil, fmt.Errorf("prepare list insertion")
		}

		_, err = dbc.Exec(context.Background(), `INSERT INTO banner_tag_feature(banner_id, tag_id, feature_id) VALUES ($1, $2, $3) RETURNING banner_id;`, banners[i].BannerId, banners[i].TagIds[0], banners[i].FeatureId)
		if err != nil {
			return nil, fmt.Errorf("prepare list insertion")
		}
	}

	return banners, nil
}
