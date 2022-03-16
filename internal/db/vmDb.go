package db

import (
	"awesomeProjectRucenter/internal/model"
	"awesomeProjectRucenter/pkg/erx"
	"awesomeProjectRucenter/pkg/tools"
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"strings"
)

const (
	INSERT_VMS        = `INSERT INTO vms (name, uuid) VALUES ($1, $2);`
	INSERT_DISKS      = `INSERT INTO disks (name, size, uuid) VALUES ($1, $2, $3);`
	INSERT_MANY_DISKS = `INSERT INTO disks (name, size, vm) VALUES (?, ?, ?);`
	INSERT_MANY_VMS   = `INSERT INTO vms (name, uuid) VALUES (?, ?);`
	GET_DISKS         = `SELECT * FROM disks LIMIT $1 OFFSET $2`
	COUNT_DISKS       = `SELECT COUNT (*) FROM disks`
	GET_VMS           = `SELECT * FROM vms LIMIT $1 OFFSET $2`
	COUNT_VMS         = `SELECT COUNT (*) FROM vms`
)

var log *logrus.Logger

func init() {
	log = tools.GetLogrusInstance("")
}

type arrVms struct {
	Arr []model.VmResults `json:"arr"`
}

type arrDisks struct {
	Arr []model.DiskResults `json:"arr"`
}

type VmDbSqlite struct {
	db *sql.DB
}

func NewDbSqlite(db *sql.DB) *VmDbSqlite {
	return &VmDbSqlite{db: db}
}

func (r *VmDbSqlite) WriteFromFile(data []byte) error {
	switch {
	case bytes.Contains(data, []byte("vms")):
		req, valArgs, err := prepareInsertVmsRequest(data)
		if err != nil {
			return erx.New(err)
		}
		_, err = r.db.Exec(req, valArgs...)

	case bytes.Contains(data, []byte("disks")):
		req, valArgs, err := prepareInsertDisksRequest(data)
		if err != nil {
			return erx.New(err)
		}
		_, err = r.db.Exec(req, valArgs...)
	default:
		return fmt.Errorf("unable to parse data\n")
	}

	return nil
}

func (r *VmDbSqlite) GetDisksCount() (int, error) {
	var count int

	row := r.db.QueryRow(COUNT_DISKS)

	err := row.Scan(&count)
	if err != nil {
		return 0, erx.New(err)
	}

	return count, nil
}

func (r *VmDbSqlite) GetVmsCount() (int, error) {
	var count int

	row := r.db.QueryRow(COUNT_VMS)

	err := row.Scan(&count)
	if err != nil {
		return 0, erx.New(err)
	}

	return count, nil
}

func (r VmDbSqlite) GetDisksWithLimit(limit, offset int) ([]model.Disk, error) {
	var disk model.Disk
	disks := make([]model.Disk, 0, limit)

	rows, err := r.db.Query(GET_DISKS, limit, offset)
	defer rows.Close()
	if err != nil {
		return nil, erx.New(err)
	}
	for rows.Next() {
		err = rows.Scan(&disk.Id, &disk.Name, &disk.Size, &disk.Vm)
		if err != nil {
			return nil, erx.New(err)
		}
		disks = append(disks, disk)
	}

	return disks, nil
}

func (r VmDbSqlite) GetVmsWithLimit(limit, offset int) ([]model.Vm, error) {
	var vm model.Vm
	vms := make([]model.Vm, 0, limit)

	rows, err := r.db.Query(GET_VMS, limit, offset)
	defer rows.Close()
	if err != nil {
		return nil, erx.New(err)
	}
	for rows.Next() {
		err = rows.Scan(&vm.Id, &vm.Name, &vm.Uuid)
		if err != nil {
			return nil, erx.New(err)
		}
		vms = append(vms, vm)
	}

	return vms, nil
}

func prepareInsertVmsRequest(bytes []byte) (req string, valArgs []interface{}, err error) {
	d := new(arrVms)
	err = json.Unmarshal(bytes, d)
	if err != nil {
		log.Error(erx.New(err))
		return
	}
	var rowCount int
	if d != nil {
		rowCount = d.Arr[0].Count
	}
	valArgs = make([]interface{}, 0, rowCount*2)
	valStrings := make([]string, 0, rowCount)
	for _, v := range d.Arr {
		for _, val := range v.Results {
			valArgs = append(valArgs, val.Name)
			valArgs = append(valArgs, val.Uuid)
			if err != nil {
				log.Error(erx.New(err))
				return
			}
			valStrings = append(valStrings, INSERT_MANY_VMS)
		}
	}
	req = strings.Join(valStrings, "")

	return
}

func prepareInsertDisksRequest(bytes []byte) (req string, valArgs []interface{}, err error) {
	d := new(arrDisks)
	err = json.Unmarshal(bytes, d)
	if err != nil {
		log.Error(erx.New(err))
		return
	}
	var rowCount int
	if d != nil {
		rowCount = d.Arr[0].Count
	}
	valArgs = make([]interface{}, 0, rowCount*3)
	valStrings := make([]string, 0, rowCount)
	for _, v := range d.Arr {
		for _, val := range v.Results {
			valArgs = append(valArgs, val.Name)
			valArgs = append(valArgs, val.Size)
			valArgs = append(valArgs, val.Vm)
			if err != nil {
				log.Error(erx.New(err))
				return
			}
			valStrings = append(valStrings, INSERT_MANY_DISKS)
		}
	}
	req = strings.Join(valStrings, "")

	return
}
