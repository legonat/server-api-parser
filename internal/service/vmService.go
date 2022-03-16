package service

import (
	"awesomeProjectRucenter/internal/config"
	"awesomeProjectRucenter/internal/db"
	"awesomeProjectRucenter/internal/model"
	"awesomeProjectRucenter/pkg/erx"
	"awesomeProjectRucenter/pkg/tools"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
)

var log *logrus.Logger

func init() {
	log = tools.GetLogrusInstance("")
}

type VmService struct {
	repo db.VmDb
}

func NewVmService(repo db.VmDb) *VmService {
	return &VmService{repo: repo}
}

func (s *VmService) WriteFromFile(path string) error {
	_, err := os.Stat(path)
	if err != nil {
		log.Error(erx.New(err))
		return err
	}
	b, err := os.ReadFile(path)
	if err != nil {
		log.Error(erx.New(err))
		return err
	}
	c := new(bytes.Buffer)
	err = json.Compact(c, b)
	if err != nil {
		log.Error(erx.New(err))
		return err
	}
	err = s.repo.WriteFromFile(c.Bytes())
	if err != nil {
		return erx.New(err)
	}

	return nil
}

func (s *VmService) InitDbWithData() error {
	err := db.InitVmDb("./data")
	if err != nil {
		log.Println(err)
	}

	err = s.WriteFromFile("./disks.txt")
	if err != nil {
		return erx.New(err)
	}

	err = s.WriteFromFile("./data.txt")
	if err != nil {
		return erx.New(err)
	}

	return nil
}

func (s *VmService) ReinitDbWithData() error {
	err := db.ForceInitVmDb("./data")
	if err != nil {
		log.Println(err)
	}

	err = s.WriteFromFile("./disks.txt")
	if err != nil {
		return erx.New(err)
	}

	err = s.WriteFromFile("./data.txt")
	if err != nil {
		return erx.New(err)
	}

	return nil
}

func (s *VmService) GetDisksWithLimit(limit, offset int) (res model.DiskResults, err error) {
	res.Count, err = s.repo.GetDisksCount()
	if err != nil {
		return model.DiskResults{}, erx.New(err)
	}
	res.Results, err = s.repo.GetDisksWithLimit(limit, offset)
	if err != nil {
		return model.DiskResults{}, erx.New(err)
	}

	cfg := config.GetConfigInstance()
	domain := cfg.Server.Domain
	if cfg.Server.Port != 0 {
		domain = fmt.Sprintf("http://%v:%d", domain, cfg.Server.Port)
	}

	switch {
	case offset == 0:
		res.Previous = ""
		res.Next = fmt.Sprintf("%v/disks/?limit=%d&offset=%d", domain, limit, offset+limit)
	case offset > 0 && offset+limit < res.Count:
		res.Previous = fmt.Sprintf("%v/disks/?limit=%d&offset=%d", domain, limit, offset-limit)
		res.Next = fmt.Sprintf("%v/disks/?limit=%d&offset=%d", domain, limit, offset+limit)
	case offset+limit > res.Count:
		res.Previous = fmt.Sprintf("%v/disks/?limit=%d&offset=%d", domain, limit, offset-limit)
		res.Next = ""
	}

	return res, nil
}

func (s *VmService) GetVmsWithLimit(limit, offset int) (res model.VmResults, err error) {
	res.Count, err = s.repo.GetVmsCount()
	if err != nil {
		return model.VmResults{}, erx.New(err)
	}
	res.Results, err = s.repo.GetVmsWithLimit(limit, offset)
	if err != nil {
		return model.VmResults{}, erx.New(err)
	}

	cfg := config.GetConfigInstance()
	domain := cfg.Server.Domain
	if cfg.Server.Port != 0 {
		domain = fmt.Sprintf("http://%v:%d", domain, cfg.Server.Port)
	}
	switch {
	case offset == 0:
		res.Previous = ""
		res.Next = fmt.Sprintf("%v/vms/?limit=%d&offset=%d", domain, limit, offset+limit)
	case offset > 0 && offset+limit < res.Count:
		res.Previous = fmt.Sprintf("%v/vms/?limit=%d&offset=%d", domain, limit, offset-limit)
		res.Next = fmt.Sprintf("%v/vms/?limit=%d&offset=%d", domain, limit, offset+limit)
	case offset+limit > res.Count:
		res.Previous = fmt.Sprintf("%v/vms/?limit=%d&offset=%d", domain, limit, offset-limit)
		res.Next = ""
	}

	return res, nil
}

func (s *VmService) GetDisksCount() (int, error) {
	return s.repo.GetDisksCount()
}

func (s *VmService) GetVmsCount() (int, error) {
	return s.repo.GetVmsCount()
}
