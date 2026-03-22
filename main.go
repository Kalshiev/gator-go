package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/kalshiev/gator-go/internal/config"
	"github.com/kalshiev/gator-go/internal/database"
	_ "github.com/lib/pq"
)

type state struct {
	config *config.Config
	db     *database.Queries
}

func main() {

	var currentState state

	cfg, err := config.Read()
	if err != nil {
		os.Exit(1)
	}

	currentState.config = &cfg

	db, err := sql.Open("postgres", currentState.config.Db_URL)
	if err != nil {
		fmt.Println("Database connection failed")
		os.Exit(1)
	}

	dbQueries := database.New(db)

	currentState.db = dbQueries

	cmds := commands{
		cmdMap: make(map[string]func(*state, command) error),
	}
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerGetUsers)
	cmds.register("agg", handlerAgg)
	cmds.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	cmds.register("feeds", handlerFeeds)
	cmds.register("follow", middlewareLoggedIn(handlerFollow))
	cmds.register("following", middlewareLoggedIn(handlerFollowing))
	cmds.register("unfollow", middlewareLoggedIn(handlerUnfollow))
	cmds.register("browse", middlewareLoggedIn(handlerBrowse))

	if len(os.Args) < 2 {
		fmt.Println("Not enough arguments provided")
		os.Exit(1)
	}

	var currentCommand command

	currentCommand.name = os.Args[1]
	currentCommand.argv = os.Args[2:]

	err = cmds.run(&currentState, currentCommand)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
