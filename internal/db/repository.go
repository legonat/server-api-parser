package db

import (
	"awesomeProjectRucenter/internal/model"
	"database/sql"
)

type VmDb interface {
	WriteFromFile(data []byte) error
	GetDisksWithLimit(limit, offset int) ([]model.Disk, error)
	GetVmsWithLimit(limit, offset int) ([]model.Vm, error)
	GetDisksCount() (int, error)
	GetVmsCount() (int, error)
}

type Repository struct {
	VmDb
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		NewDbSqlite(db),
	}
}
