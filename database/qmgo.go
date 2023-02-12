package database

import (
	"github.com/duxphp/duxgo/v2/core"
	"github.com/gookit/event"
	"github.com/qiniu/qmgo"
)

func QmgoInit() {
	dbConfig := core.Config["database"].GetStringMapString("mongoDB")

	var auth = ""
	if dbConfig["username"] != "" && dbConfig["password"] != "" {
		auth = dbConfig["username"] + ":" + dbConfig["password"] + "@"
	}

	client, err := qmgo.NewClient(core.Ctx, &qmgo.Config{Uri: "mongodb://" + auth + dbConfig["host"] + ":" + dbConfig["port"]})
	if err != nil {
		panic("qmgo error :" + err.Error())
	}
	core.Mgo = client.Database(dbConfig["dbname"])

	event.On("app.close", event.ListenerFunc(func(e event.Event) error {
		return client.Close(core.Ctx)
	}))
}
