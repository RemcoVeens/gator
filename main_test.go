package main_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/RemcoVeens/gator/internal/commands"
	"github.com/RemcoVeens/gator/internal/config"
)

func TestRead(t *testing.T) {
	// --- Test Case 1: Successful Read ---
	t.Run("BaseInit", func(t *testing.T) {
		test_args := make([]string, 2)
		test_args = append(test_args, "login", "remco")
		_ = commands.ParceInput(test_args)
		conf, err := config.Read()
		if err != nil {
			fmt.Println("cound not read config")
			os.Exit(1)
		}
		_ = commands.NewState(&conf)

		_ = commands.NewCommands()

		fmt.Println("we good bro")
	})

	// t.Run("FileNotFound", func(t *testing.T) {

	// })
}
