package main

import (
	"encoding/json"
	"github.com/openaccounting/oa-server/core/model/db"
	"github.com/openaccounting/oa-server/core/model/types"
	"log"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Usage: migrate3.go <upgrade/downgrade>")
	}

	command := os.Args[1]

	if command != "upgrade" && command != "downgrade" {
		log.Fatal("Usage: migrate3.go <upgrade/downgrade>")
	}

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
	db, err := db.NewDB(connectionString, "mysql")

	if command == "upgrade" {
		err = upgrade(db)
	} else {
		err = downgrade(db)
	}

	if err != nil {
		log.Fatal(err)
	}

	log.Println("done")
}

func upgrade(db *db.DB) (err error) {
	tx, err := db.Begin()

	if err != nil {
		return
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p) // re-throw panic after Rollback
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	query1 := "CREATE TABLE budgetitem (id INT UNSIGNED NOT NULL AUTO_INCREMENT, orgId BINARY(16) NOT NULL, accountId BINARY(16) NOT NULL, inserted BIGINT UNSIGNED NOT NULL, amount BIGINT NOT NULL, PRIMARY KEY(id)) ENGINE=InnoDB;"

	if _, err = tx.Exec(query1); err != nil {
		return
	}

	return
}

func downgrade(db *db.DB) (err error) {
	tx, err := db.Begin()

	if err != nil {
		return
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p) // re-throw panic after Rollback
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	query1 := "DROP TABLE budgetitem"

	if _, err = tx.Exec(query1); err != nil {
		return
	}

	return
}
