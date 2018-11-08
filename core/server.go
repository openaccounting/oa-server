package main

import (
	"encoding/json"
	"github.com/openaccounting/oa-server/core/api"
	"github.com/openaccounting/oa-server/core/auth"
	"github.com/openaccounting/oa-server/core/model"
	"github.com/openaccounting/oa-server/core/model/db"
	"github.com/openaccounting/oa-server/core/model/types"
	"github.com/openaccounting/oa-server/core/util"
	"log"
	"net/http"
	"os"
	"strconv"
	//"fmt"
)

func main() {
	//filename is the path to the json config file
	var config types.Config
	file, err := os.Open("./config.json")

	if err != nil {
		log.Fatal(err)
	}

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)

	if err != nil {
		log.Fatal(err)
	}

	connectionString := config.User + ":" + config.Password + "@/" + config.Database

	db, err := db.NewDB(connectionString)

	bc := &util.StandardBcrypt{}

	model.NewModel(db, bc, config)
	auth.NewAuthService(db, bc)

	app, err := api.Init(config.ApiPrefix)

	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(http.ListenAndServeTLS(":"+strconv.Itoa(config.Port), config.CertFile, config.KeyFile, app.MakeHandler()))
}
