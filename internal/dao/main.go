package dao

// Interface is the interface for letscrum.
type Interface interface {
	LetscrumDao() LetscrumDao
	UserDao() UserDao
	OrgDao() OrgDao
	ProjectDao() ProjectDao
	SprintDao() SprintDao
	WorkItemDao() WorkItemDao
	TaskDao() TaskDao
}
