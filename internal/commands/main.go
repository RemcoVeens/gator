package commands

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/RemcoVeens/gator/internal/config"
	"github.com/RemcoVeens/gator/internal/database"
	"github.com/google/uuid"
)

type state struct {
	Config *config.Config
	DB     *database.Queries
}

func NewState(conf *config.Config) state {
	return state{Config: conf}
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
	comm.Register("register", handlerRegister)
	return comm
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("the login handler expects a single argument, the username")
	}
	username := cmd.Args[0]
	user, err := s.DB.GetUser(context.Background(), sql.NullString{
		String: username,
		Valid:  username != "",
	})
	if err != nil {
		return fmt.Errorf("could not log in %s: %w", username, err)
	}
	s.Config.SetUser(user.Name.String)
	fmt.Printf("user: %v has been login\n", user.Name.String)
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("please provide a name to register")
	}
	username := cmd.Args[0]
	user, err := s.DB.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      sql.NullString{String: username, Valid: username != ""},
	})
	if err != nil {
		return fmt.Errorf("could not insert %v: %w", username, err)
	}
	s.Config.SetUser(user.Name.String)
	fmt.Printf("user: '%s' has been created, at %v \n", user.Name.String, user.CreatedAt)
	return nil
}
