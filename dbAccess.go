package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type dbAccess struct {
	conString string
}

//CommandList matching structure to the table CommandList in our database
type CommandList struct {
	ID      int    `json:"id"`
	Command string `json:"command"`
}

//AmharicWords matching structure to the table AmharicWords in our database
type AmharicWords struct {
	ID        int    `json:"id"`
	CommandID int    `json:"commandId"`
	Word      string `json:"word"`
}

func (con dbAccess) getListofCommands() {
	// open the database connection
	db, err := sql.Open("mysql", con.conString)

	if err != nil {
		panic(err.Error())
	}

	//make sure to always close the connection
	defer db.Close()

	var listOfCommands []CommandList
	//select query
	results, err := db.Query("SELECT id, command from CommandList")
	if err != nil {
		panic(err.Error())
	}

	for results.Next() {
		var command CommandList
		err = results.Scan(&command.ID, &command.Command)
		if err != nil {
			panic(err.Error())
		}
		listOfCommands = append(listOfCommands, command)
	}

	fmt.Println(listOfCommands)

}
