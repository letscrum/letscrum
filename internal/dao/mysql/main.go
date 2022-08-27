package mysql

import (
	"fmt"
	"github.com/letscrum/letscrum/internal/dao"
	"github.com/letscrum/letscrum/pkg/db"
	"gorm.io/gorm"
	"sync"
)

type Dao struct {
	Db *gorm.DB
}

func (d *Dao) LetscrumDao() dao.LetscrumDao {
	return NewLetscrumDao(d.Db)
}

func (d *Dao) ProjectDao() dao.ProjectDao {
	return NewProjectDao(d.Db)
}

func GetDao(opts *db.Options) (dao.Interface, error) {
	var daoInterface dao.Interface
	var once sync.Once

	if opts == nil && daoInterface == nil {
		return nil, fmt.Errorf("failed to get mysql hive dao")
	}

	var err error
	var dbIns *gorm.DB
	once.Do(func() {
		options := &db.Options{
			Host:                  opts.Host,
			Port:                  opts.Port,
			Username:              opts.Username,
			Password:              opts.Password,
			Database:              opts.Database,
			MaxIdleConnections:    opts.MaxIdleConnections,
			MaxOpenConnections:    opts.MaxOpenConnections,
			MaxConnectionLifeTime: opts.MaxConnectionLifeTime,
			Logger:                opts.Logger,
		}
		dbIns, err = db.NewGORM(options)
		daoInterface = &Dao{dbIns}
	})

	if daoInterface == nil || err != nil {
		return nil, fmt.Errorf("failed to get mysql hive dao, : %+v, error: %w", daoInterface, err)
	}

	return daoInterface, nil
}
