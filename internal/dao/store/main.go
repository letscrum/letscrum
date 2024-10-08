package store

import (
	"fmt"
	"sync"

	"github.com/letscrum/letscrum/internal/dao"
	"github.com/letscrum/letscrum/internal/model"
	"github.com/letscrum/letscrum/pkg/db"
	"gorm.io/gorm"
)

type Dao struct {
	DB *gorm.DB
}

func (d *Dao) WorkItemDao() dao.WorkItemDao {
	return NewWorkItemDao(d.DB)
}

func (d *Dao) TaskDao() dao.TaskDao {
	return NewTaskDao(d.DB)
}

func (d *Dao) LetscrumDao() dao.LetscrumDao {
	return NewLetscrumDao(d.DB)
}

func (d *Dao) UserDao() dao.UserDao {
	return NewUserDao(d.DB)
}

func (d *Dao) ProjectDao() dao.ProjectDao {
	return NewProjectDao(d.DB)
}

func (d *Dao) OrgDao() dao.OrgDao {
	return NewOrgDao(d.DB)
}

func (d *Dao) SprintDao() dao.SprintDao {
	return NewSprintDao(d.DB)
}

type ColumnType interface {
	DatabaseTypeName() string // varchar
}

func GetDao(opts *db.Options) (dao.Interface, error) {
	var daoInterface dao.Interface
	var once sync.Once

	if opts == nil {
		return nil, fmt.Errorf("failed to get database options")
	}

	var err error
	var dbIns *gorm.DB
	once.Do(func() {
		options := &db.Options{
			Driver:                opts.Driver,
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
	})

	initErr := dbIns.AutoMigrate(
		&model.User{},
		&model.Org{},
		&model.OrgUser{},
		&model.Project{},
		&model.Sprint{},
		&model.Epic{},
		&model.Feature{},
		&model.WorkItem{},
		&model.Task{},
		&model.ItemLog{},
		&model.SprintBurndown{},
	)

	if initErr != nil {
		return nil, fmt.Errorf("failed to init letscrum database: %w", err)
	}

	daoInterface = &Dao{dbIns}

	if err != nil {
		return nil, fmt.Errorf("failed to get letscrum dao: %+v, error: %w", daoInterface, err)
	}

	return daoInterface, nil
}
