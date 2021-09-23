package database

import (
	"errors"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type Config struct {
	ContainerID string `yaml:"container_id"`
	Address     string `yaml:"address"`
	User        string `yaml:"user"`
	Password    string `yaml:"password"`
	Database    string `yaml:"database"`
}

type DB struct {
	*sqlx.DB
}

// Seed insert rows into database tabel
func (db *DB) Seed(seedSql string) error {
	result, err := db.Exec(seedSql)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("insert 0 rows")
	}

	return nil
}

// Open get mysql db connection
func Open(cfg *Config) (*DB, error) {
	mysqlCfg := mysql.Config{
		Addr:                 cfg.Address,
		User:                 cfg.User,
		Passwd:               cfg.Password,
		DBName:               cfg.Database,
		AllowNativePasswords: true,
		Collation:            "utf8mb4_unicode_ci",
		ParseTime:            true,
	}

	db, err := sqlx.Open("mysql", mysqlCfg.FormatDSN())
	if err != nil {
		return nil, err
	}

	return &DB{db.Unsafe()}, nil
}
