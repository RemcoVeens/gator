package commands

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/RemcoVeens/gator/internal/config"
	"github.com/RemcoVeens/gator/internal/database"
	"github.com/RemcoVeens/gator/internal/feed"
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
	comm.Register("reset", handlerReset)
	comm.Register("users", handlerUsers)
	comm.Register("agg", handlerAgg)
	comm.Register("addfeed", handlerAddFeed)
	comm.Register("feeds", handlerFeeds)
	comm.Register("follow", handlerFollow)
	comm.Register("following", handlerFollowing)

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
func handlerReset(s *state, cmd command) error {
	err := s.DB.DeleteAllUsers(context.Background())
	if err != nil {
		return fmt.Errorf("could not reset db: %w", err)
	}
	return nil
}
func handlerUsers(s *state, cmd command) error {
	users, err := s.DB.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("could not get users:\n%w", err)
	}
	for _, user := range users {
		user_name := user.Name.String
		if user_name == s.Config.CurentUserName {
			fmt.Println(user_name, "(current)")
		} else {
			fmt.Println(user_name)
		}
	}

	return nil
}
func handlerAgg(s *state, cmd command) error {
	url := "https://www.wagslane.dev/index.xml"
	Feed, err := feed.FetchFeed(context.Background(), url)
	if err != nil {
		return err
	}
	fmt.Println(*Feed)

	return nil
}
func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.Args) != 2 {
		return fmt.Errorf("please provide a name and url to add a feed")
	}
	name := cmd.Args[0]
	url := cmd.Args[1]
	fmt.Printf("current user: %s\n", s.Config.CurentUserName)
	user, err := s.DB.GetUser(context.Background(), sql.NullString{
		String: s.Config.CurentUserName,
		Valid:  true,
	})
	if err != nil {
		return err
	}
	fmt.Println("ID:", user.ID)
	feed, err := s.DB.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
		Url:       url,
		UserID:    user.ID,
	})
	if err != nil {
		return fmt.Errorf("could not insert feed %v: \n\r%w", name, err)
	}
	fmt.Printf("feed: '%s' has been created, at %v \n", feed.Name, feed.CreatedAt)
	_, err = s.DB.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("could not follow: %w", err)
	}
	fmt.Printf("feed: %s is now followd by you (%s)\n", feed.Name, user.Name.String)
	return nil
}
func handlerFeeds(s *state, cmd command) error {
	feeds, err := s.DB.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("could not get feeds:\n%w", err)
	}
	for _, feed := range feeds {
		user, err := s.DB.GetUserByID(context.Background(), feed.UserID)
		if err != nil {
			return err
		}
		fmt.Printf("feed: '%s' (%s) created by %s\n", feed.Name, feed.Url, user.Name.String)
	}
	return nil
}
func handlerFollow(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("please provide a url to follow a feed")
	}
	url := cmd.Args[0]
	feed, err := s.DB.GetFeedFromUrl(context.Background(), url)
	if err != nil {
		return fmt.Errorf("feed not found")
	}
	user, err := s.DB.GetUser(context.Background(), sql.NullString{
		String: s.Config.CurentUserName,
		Valid:  true,
	})
	if err != nil {
		return err
	}
	ffr, err := s.DB.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("could not follow feed %v: \n\r%w", feed.Name, err)
	}
	fmt.Printf("you (%s) now follow feed: '%s' (%s) has been followed\n", ffr.UserName.String, ffr.FeedName, ffr.FeedUrl)
	return nil
}
func handlerFollowing(s *state, cmd command) error {
	user, err := s.DB.GetUser(context.Background(), sql.NullString{
		String: s.Config.CurentUserName,
		Valid:  true,
	})
	if err != nil {
		return err
	}
	following, err := s.DB.GetFeedFollowsByUser(context.Background(), user.ID)
	if len(following) == 0 {
		return nil
	}
	if err != nil {
		return fmt.Errorf("could get feeds from %v: \n\r%w", user.Name, err)
	}
	for _, feed := range following {
		fmt.Println("-", feed.Name)
	}
	return nil
}
