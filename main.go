package main

import (
	"database/sql"
	"log"
	"os"
	"strings"

	_ "github.com/lib/pq"
	"github.com/thomas-reed/gator/internal/config"
	"github.com/thomas-reed/gator/internal/database"
)

type state struct {
	db *database.Queries
	cfg *config.Config
}

func main() {
	// read config file
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("Error reading config: %v", err)
	}

	// set up db connection
	db, err := sql.Open("postgres", cfg.DbURL)
	if err != nil {
		log.Fatalf("Error opening DB connection: %v", err)
	}
	dbQueries := database.New(db)

	// save state for use in commands
	programState := &state{
		db: dbQueries,
		cfg: &cfg,
	}

	// build command registry
	cmds := commands{
		registry: make(map[string]func(s *state, cmd command) error),
	}
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerListUsers)
	cmds.register("agg", handlerAggregate)
	cmds.register("addfeed", handlerAddFeed)
	cmds.register("feeds", handlerListFeeds)

	// parse cmd line arguments
	if len(os.Args) < 2 {
		log.Fatalln("Too few arguments.  Usage: gator <command> [args...]")
	}
	cmdName := strings.ToLower(os.Args[1])
	cmdArgs := os.Args[2:]

	// run given command
	if err = cmds.run(programState, command{name: cmdName, args: cmdArgs}); err != nil {
		log.Fatalf("Error running %s command: %s\n", cmdName, err)
	}
	os.Exit(0);
}