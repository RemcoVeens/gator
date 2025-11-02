package main

import (
	"fmt"
	"os"

	"github.com/RemcoVeens/gator/internal/commands"
	"github.com/RemcoVeens/gator/internal/config"
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

	fmt.Println("we good bro")
	fmt.Println("args:", args)
	err = comms.Run(&stat, comm)
	if err != nil {
		os.Exit(1)
	}
}
