package service

import (
	"awesomeProjectRucenter/internal/db"
	"awesomeProjectRucenter/internal/model"
)

type Vm interface {
	WriteFromFile(path string) error
	InitDbWithData() error
	ReinitDbWithData() error
	GetDisksWithLimit(limit, offset int) (model.DiskResults, error)
	GetVmsWithLimit(limit, offset int) (model.VmResults, error)
	GetDisksCount() (int, error)
	GetVmsCount() (int, error)
}

type Service struct {
	Vm
}

func NewService(repo *db.Repository) *Service {
	return &Service{Vm: NewVmService(repo.VmDb)}
}
