package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/openaccounting/oa-server/core/api"
	"github.com/openaccounting/oa-server/core/auth"
	"github.com/openaccounting/oa-server/core/model"
	"github.com/openaccounting/oa-server/core/model/db"
	"github.com/openaccounting/oa-server/core/model/types"
	"github.com/openaccounting/oa-server/core/util"
)

func main() {
	//filename is the path to the json config file
	var config types.Config
	file, err := os.Open("./config.json")

	if err != nil {
		log.Fatal(fmt.Errorf("failed to open ./config.json with: %s", err.Error()))
	}

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)

	if err != nil {
		log.Fatal(fmt.Errorf("failed to decode ./config.json with: %s", err.Error()))
	}

	connectionString := config.User + ":" + config.Password + "@" + config.DatabaseAddress + "/" + config.Database

	db, err := db.NewDB(connectionString)

	bc := &util.StandardBcrypt{}

	model.NewModel(db, bc, config)
	auth.NewAuthService(db, bc)

	app, err := api.Init(config.ApiPrefix)

	if err != nil {
		log.Fatal(fmt.Errorf("failed to create api instance with: %s", err.Error()))
	}

	if config.CertFile == "" || config.KeyFile == "" {
		err = http.ListenAndServe(config.Address+":"+strconv.Itoa(config.Port), app.MakeHandler())
	} else {
		err = http.ListenAndServeTLS(config.Address+":"+strconv.Itoa(config.Port), config.CertFile, config.KeyFile, app.MakeHandler())
	}
	log.Fatal(fmt.Errorf("failed to start server with: %s", err.Error()))
}
