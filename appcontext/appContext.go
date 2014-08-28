package appcontext

import (
	"github.com/fengjian0106/gomgo/database"
	"github.com/fengjian0106/gomgo/taskqueue"
)

//https://gist.github.com/elithrar/5aef354a54ba71a32e23
type AppContext struct {
	// appContext contains our local context; our database pool, session store, template
	// registry and anything else our handlers need to access. We'll create an instance of it
	// in our main() function and then explicitly pass a reference to it for our handlers to access.
	Db *database.Database

	ReqRepTaskQueue taskqueue.ReqRepTaskQueue
	PubTaskQueue    taskqueue.PubTaskQueue
}

func New() (*AppContext, error) {
	//<1> mongodb
	db, err := database.New()
	if err != nil {
		return nil, err
	}

	reqRepTaskQueue := taskqueue.NewReqRepTaskQueue(10)
	pubTaskQueue := taskqueue.NewPubTaskQueue(5)

	ctx := AppContext{
		Db:              db,
		ReqRepTaskQueue: reqRepTaskQueue,
		PubTaskQueue:    pubTaskQueue,
	}

	return &ctx, nil
}

func (c *AppContext) FreeResource() {
	if c.Db != nil {
		c.Db.CloseSession()
	}
}
