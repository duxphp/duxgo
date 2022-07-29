package database

import (
	"github.com/duxphp/duxgo/global"
	"github.com/gookit/event"
	"github.com/qiniu/qmgo"
)

func QmgoInit() {
	dbConfig := global.Config["database"].GetStringMapString("mongoDB")

	var auth = ""
	if dbConfig["username"] != "" && dbConfig["password"] != "" {
		auth = dbConfig["username"] + ":" + dbConfig["password"] + "@"
	}

	client, err := qmgo.NewClient(global.Ctx, &qmgo.Config{Uri: "mongodb://" + auth + dbConfig["host"] + ":" + dbConfig["port"]})
	if err != nil {
		panic("qmgo error :" + err.Error())
	}
	global.Mgo = client.Database(dbConfig["dbname"])

	event.On("app.close", event.ListenerFunc(func(e event.Event) error {
		return client.Close(global.Ctx)
	}))
}
