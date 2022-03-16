package db

import (
	"awesomeProjectRucenter/pkg/erx"
	"database/sql"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

const (
	CREATE_TABLE_VMS = `CREATE TABLE vms(
    id            INTEGER PRIMARY KEY AUTOINCREMENT,
    name          varchar(255) not null,
    uuid          varchar(255) not null unique);
`
	CREATE_TABLE_DISKS = `CREATE TABLE disks(
	id            INTEGER PRIMARY KEY AUTOINCREMENT,
	name          varchar(255) not null,
	size          integer not null,
	vm            varchar(255) not null); 
`
	UPDATE_VMS_SEQ   = `UPDATE SQLITE_SEQUENCE SET seq = 10000 WHERE name = 'vms'`
	UPDATE_DISKS_SEQ = `UPDATE SQLITE_SEQUENCE SET seq = 15962 WHERE name = 'disks'`
	INSERT_DISKS_SEQ = `INSERT INTO SQLITE_SEQUENCE (name, seq) VALUES ('disks', 15962)`
	INSERT_VMS_SEQ   = `INSERT INTO SQLITE_SEQUENCE (name, seq) VALUES ('vms', 10000)`
)

func NewSqliteDB(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func InitVmDb(path string) error {

	pathDb := path + "/vm.db"
	//cfg, err := config.GetConfig()
	//if err != nil {
	//	tools.LogErr(err)
	//	return err
	//}
	//(*cfg).UsersDB.PathDb = pathDb

	f, err := os.Stat(pathDb)
	if err != nil && !os.IsNotExist(err) {
		return erx.New(err)
	}
	if f != nil {
		return erx.NewError(605, "Data Base is already exist")
	}
	if os.IsNotExist(err) {
		err = os.MkdirAll(path, 0755)
		if err != nil {
			return erx.New(err)
		}

		_, err = os.Create(pathDb)
		if err != nil {
			return erx.New(err)
		}

		db, err := sql.Open("sqlite3", pathDb)
		if err != nil {
			return erx.New(err)
		}

		defer db.Close()

		_, err = db.Exec(CREATE_TABLE_VMS)
		if err != nil {
			return erx.New(err)
		}

		_, err = db.Exec(INSERT_VMS_SEQ)
		if err != nil {
			return erx.New(err)
		}

		_, err = db.Exec(CREATE_TABLE_DISKS)
		if err != nil {
			return erx.New(err)
		}

		_, err = db.Exec(INSERT_DISKS_SEQ)
		if err != nil {
			return erx.New(err)
		}

	}

	return err
}

func ForceInitVmDb(path string) error {
	pathDb := path + "/vm.db"

	err := os.MkdirAll(path, 0755)
	if err != nil {
		return erx.New(err)
	}

	_, err = os.Create(pathDb)
	if err != nil {
		return erx.New(err)
	}

	db, err := sql.Open("sqlite3", pathDb)
	if err != nil {
		return erx.New(err)
	}

	defer db.Close()

	_, err = db.Exec(CREATE_TABLE_VMS)
	if err != nil {
		return erx.New(err)
	}

	_, err = db.Exec(INSERT_VMS_SEQ)
	if err != nil {
		return erx.New(err)
	}

	_, err = db.Exec(CREATE_TABLE_DISKS)
	if err != nil {
		return erx.New(err)
	}

	_, err = db.Exec(INSERT_DISKS_SEQ)
	if err != nil {
		return erx.New(err)
	}
	log.Info("DB Init Success")
	return nil
}
