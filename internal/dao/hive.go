package dao

//go:generate mockgen -destination mock/mock_hive.go --build_flags=--mod=mod github.com/daocloud/skoala/app/hive/internal/dao HiveDao,BookDao,RegistryDao

// HiveDao is the interface for hive.
type HiveDao interface {
	RegistryDao() RegistryDao
	BookDao() BookDao
}
