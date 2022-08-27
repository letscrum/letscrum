package dao

// Interface is the interface for letscrum.
type Interface interface {
	LetscrumDao() LetscrumDao
	ProjectDao() ProjectDao
	DemoDbDao() DemoDbDao
}
