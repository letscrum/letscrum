package dao

// Interface is the interface for letscrum.
type Interface interface {
	LetscrumDao() LetscrumDao
	UserDao() UserDao
	ProjectDao() ProjectDao
	ProjectMemberDao() ProjectMemberDao
	SprintDao() SprintDao
	SprintMemberDao() SprintMemberDao
	WorkItemDao() WorkItemDao
	TaskDao() TaskDao
}
