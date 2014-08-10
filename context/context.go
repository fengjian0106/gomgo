package context

import "github.com/fengjian0106/gomgo/database"

type Context struct {
	Db *database.Database
}

func New() (*Context, error) {
	db, err := database.New()
	if err != nil {
		return nil, err
	}

	return &Context{Db: db}, nil
}

func (c *Context) FreeResource() {
	if c.Db != nil {
		c.Db.CloseSession()
	}
}
