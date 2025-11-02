package commands

import (
	"fmt"

	"github.com/RemcoVeens/gator/internal/config"
)

type state struct {
	config *config.Config
}

func NewState(conf *config.Config) state {
	return state{config: conf}
}

type commands struct {
	command map[string]func(*state, command) error
}

type command struct {
	Name string
	Args []string
}

func ParceInput(args []string) command {
	return command{
		Name: args[0],
		Args: args[1:],
	}
}

func (c *commands) Run(s *state, cmd command) (err error) {
	err = c.command[cmd.Name](s, cmd)
	if err != nil {
		return err
	}
	return

}
func (c *commands) Register(name string, f func(*state, command) error) {
	c.command[name] = f
}
func NewCommands() commands {
	comm := commands{}
	comm.command = make(map[string]func(*state, command) error)
	comm.Register("login", handlerLogin)
	return comm
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("the login handler expects a single argument, the username")
	}
	username := cmd.Args[0]
	(s).config.SetUser(username)
	fmt.Printf("user: %v has been set", username)
	return nil
}
