package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/giapoldo/blog-aggregator/internal/config"
	"github.com/giapoldo/blog-aggregator/internal/database"
	_ "github.com/lib/pq"
)

type state struct {
	cfg *config.Config
	db  *database.Queries
}

type command struct {
	name string
	args []string
}

type commands struct {
	commands map[string]func(*state, command) error
}

// This method registers a new handler function for a command name.
func (c *commands) register(name string, f func(*state, command) error) {

	c.commands[name] = f
}

// This method runs a given command with the provided state if it exists.
func (c *commands) run(s *state, cmd command) error {

	err := c.commands[cmd.name](s, cmd)
	if err != nil {
		return err
	}
	return nil

}

func main() {

	programState := state{cfg: &config.Config{}}
	configFile := config.Read()
	*programState.cfg = configFile
	db, err := sql.Open("postgres", programState.cfg.DBUrl)
	programState.db = database.New(db)

	cmds := commands{commands: map[string]func(*state, command) error{}}

	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerListUsers)
	cmds.register("agg", handlerAgg)
	cmds.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	cmds.register("feeds", handlerFeeds)
	cmds.register("follow", middlewareLoggedIn(handlerFollow))
	cmds.register("following", middlewareLoggedIn(handlerFollowing))
	cmds.register("unfollow", middlewareLoggedIn(handlerUnfollow))

	if len(os.Args) < 2 {

		fmt.Println("Not enough arguments")
		os.Exit(1)
	}

	cmd := command{
		name: os.Args[1],
		args: os.Args[2:],
	}

	err = cmds.run(&programState, cmd)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
