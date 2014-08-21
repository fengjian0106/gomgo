package appcontext

import "github.com/fengjian0106/gomgo/database"

//https://gist.github.com/elithrar/5aef354a54ba71a32e23
type AppContext struct {
	// appContext contains our local context; our database pool, session store, template
	// registry and anything else our handlers need to access. We'll create an instance of it
	// in our main() function and then explicitly pass a reference to it for our handlers to access.
	Db *database.Database
}

func New() (*AppContext, error) {
	db, err := database.New()
	if err != nil {
		return nil, err
	}

	return &AppContext{Db: db}, nil
}

func (c *AppContext) FreeResource() {
	if c.Db != nil {
		c.Db.CloseSession()
	}
}
