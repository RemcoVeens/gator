package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/RemcoVeens/gator/internal/commands"
	"github.com/RemcoVeens/gator/internal/config"
	"github.com/RemcoVeens/gator/internal/database"

	_ "github.com/lib/pq"
)

func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		fmt.Println("please provide an argument of what you want to do")
		os.Exit(1)
	}
	comm := commands.ParceInput(args)
	conf, err := config.Read()
	if err != nil {
		fmt.Println("cound not read config")
		os.Exit(1)
	}
	stat := commands.NewState(&conf)

	comms := commands.NewCommands()
	db, err := sql.Open("postgres", stat.Config.DBUrl)
	if err != nil {
		fmt.Println("could not connect to db")
		os.Exit(1)
	}
	stat.DB = database.New(db)

	err = comms.Run(&stat, comm)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
