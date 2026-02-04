package db

import (
	"fmt"
	"log"
	"os"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/previousnext/tl-go/internal/model"
)

type RepositoryInterface interface {
	AutoMigrate() error
}

type Repository struct {
	dbPath string
}

func NewRepository(dbPath string) *Repository {
	return &Repository{dbPath: dbPath}
}

func (r *Repository) AutoMigrate() error {
	db := r.openDB()
	return db.AutoMigrate(
		&model.TimeEntry{},
		&model.Issue{},
		&model.Project{},
	)
}

func (r *Repository) openDB() *gorm.DB {
	l := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			IgnoreRecordNotFoundError: true,
		},
	)

	db, err := gorm.Open(sqlite.Open(r.dbPath), &gorm.Config{
		Logger: l,
	})

	if err != nil {
		// golangci-lint-ignore errcheck
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
	return db
}
