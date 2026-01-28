package db

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/previousnext/tl-go/internal/model"
)

type Repository struct {
	dbPath string
}

func NewRepository(dbPath string) *Repository {
	return &Repository{dbPath: dbPath}
}

func (r *Repository) InitRepository() error {
	db := r.openDB()
	return db.AutoMigrate(&model.TimeEntry{})
}

func (r *Repository) CreateTimeEntry(entry *model.TimeEntry) error {
	db := r.openDB()
	if err := db.Create(&entry).Error; err != nil {
		return err
	}
	return nil
}

func (r *Repository) FindTimeEntry(id uint) (*model.TimeEntry, error) {
	db := r.openDB()
	var entry model.TimeEntry
	if err := db.First(&entry, id).Error; err != nil {
		return nil, err
	}
	return &entry, nil
}

func (r *Repository) FindAllTimeEntries() ([]model.TimeEntry, error) {
	db := r.openDB()
	var entries []model.TimeEntry
	if err := db.Find(&entries).Error; err != nil {
		return nil, err
	}
	return entries, nil
}

func (r *Repository) UpdateTimeEntry(entry *model.TimeEntry) error {
	db := r.openDB()
	return db.Save(entry).Error
}

func (r *Repository) DeleteTimeEntry(id uint) error {
	db := r.openDB()
	return db.Delete(&model.TimeEntry{}, id).Error
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
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
	return db
}
