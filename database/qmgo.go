package database

import (
	"github.com/duxphp/duxgo/v2/registry"
	"github.com/gookit/event"
	"github.com/qiniu/qmgo"
	"github.com/samber/do"
)

func QmgoInit() {
	dbConfig := registry.Config["database"].GetStringMapString("mongoDB")

	var auth = ""
	if dbConfig["username"] != "" && dbConfig["password"] != "" {
		auth = dbConfig["username"] + ":" + dbConfig["password"] + "@"
	}

	client, err := qmgo.NewClient(registry.Ctx, &qmgo.Config{Uri: "mongodb://" + auth + dbConfig["host"] + ":" + dbConfig["port"]})
	if err != nil {
		panic("qmgo error :" + err.Error())
	}

	do.ProvideValue[*qmgo.Database](nil, client.Database(dbConfig["dbname"]))

	event.On("app.close", event.ListenerFunc(func(e event.Event) error {
		return client.Close(registry.Ctx)
	}))
}
