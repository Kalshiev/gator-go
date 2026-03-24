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
		cmdMap:   make(map[string]func(*state, command) error),
		cmdUsage: make(map[string]usage),
	}
	cmds.register("login", handlerLogin)
	cmds.registerUsage("login", "gator login [username]; [password]", "Logs in an existing user")

	cmds.register("register", handlerRegister)
	cmds.registerUsage("register", "gator register [username]; [password]", "Creates a new user with username [username] and [password]")

	cmds.register("reset", handlerReset)
	cmds.registerUsage("reset", "gator reset", "Drops all tables from database. ONLY IN DEV MODE")

	cmds.register("users", handlerGetUsers)
	cmds.registerUsage("users", "gator users", "Lists all registered users")

	cmds.register("agg", handlerAgg)
	cmds.registerUsage("agg", "gator agg [time interval]", "Fetches posts from feeds every [time interval] e.g. 5s, 5m, 5h")

	cmds.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	cmds.registerUsage("addfeed", "gator addfeed [name] [url]", "Adds a new feed with [name] and [url]")

	cmds.register("feeds", handlerFeeds)
	cmds.registerUsage("feeds", "gator feeds", "Lists all saved feeds with [name], [url] and [user]")

	cmds.register("follow", middlewareLoggedIn(handlerFollow))
	cmds.registerUsage("follow", "gator follow [url]", "Add another users feed to current user feed")

	cmds.register("following", middlewareLoggedIn(handlerFollowing))
	cmds.registerUsage("following", "gator following", "Lists all feeds being followed by the current user")

	cmds.register("unfollow", middlewareLoggedIn(handlerUnfollow))
	cmds.registerUsage("unfollow", "gator unfollow [url]", "Removes a feed from feeds followed by the current user")

	cmds.register("browse", middlewareLoggedIn(handlerBrowse))
	cmds.registerUsage("browse", "gator browse [limit: optional]", "Displays the posts from followed feeds by the current user")

	cmds.register("help", handlerHelp)
	cmds.registerUsage("help", "gator help", "Displays a list for all registered commands, usage and description")

	if len(os.Args) < 2 {
		fmt.Println("Not enough arguments provided")
		os.Exit(1)
	}

	var currentCommand command

	currentCommand.name = os.Args[1]
	currentCommand.argv = os.Args[2:]
	currentCommand.help = cmds.cmdUsage

	err = cmds.run(&currentState, currentCommand)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
