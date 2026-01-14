package main

import (
	"database/sql"
	"log"
	"os"
	"strings"
	"context"
	"fmt"
	"time"

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
	cmds.register("addfeed", loggedIn(handlerAddFeed))
	cmds.register("feeds", handlerListFeeds)
	cmds.register("follow", loggedIn(handlerFollow))
	cmds.register("following", loggedIn(handlerFollowing))
	cmds.register("unfollow", loggedIn(handlerUnfollow))
	cmds.register("browse", loggedIn(handlerBrowse))

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

func loggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		user, err := s.db.GetUserByName(context.Background(), s.cfg.CurrentUsername)
		if err != nil {
			return fmt.Errorf("Error getting user info from db:\n%w", err)
		}
		return handler(s, cmd, user)
	}
}

func scrapeFeeds(s *state) error {
	feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("Error getting next feed to fetch:\n%w", err)
	}
	rssFeed, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {
		return fmt.Errorf("Error fetching feed:\n%w", err)
	}
	if _, err = s.db.MarkFeedFetched(context.Background(), feed.ID); err != nil {
		return fmt.Errorf("Error marking feed as fetched:\n%w", err)
	}
	fmt.Printf("Latests Posts from %s:\n", feed.Name)
	for _, item := range rssFeed.Channel.Item {
		_, err := s.db.GetPostByURL(context.Background(), item.Link)
		if err == nil {
			// post already exists
			continue
		}
		parsedTime, err := parseTime(item.PubDate)
		if err != nil {
			fmt.Println("Published At time couldn't be parsed - value set to NULL")
		}
		pubTime := sql.NullTime{
			Time: parsedTime,
			Valid: err == nil,
		}
		desc := sql.NullString{
			String: item.Description,
			Valid: item.Description != "",
		}
		post, err := s.db.CreatePost(context.Background(), database.CreatePostParams{
			Title: item.Title,
			Url: item.Link,
			PublishedAt: pubTime,
			Description: desc,
			FeedID: feed.ID,
		})
		if err != nil {
			return fmt.Errorf("Error adding post to db:\n%w", err)
		}
		fmt.Printf("Post downloaded: %s\n", post.Title)
	}
	return nil
}

func parseTime(timeStr string) (time.Time, error) {
    formats := []string{
        time.RFC1123Z,
        time.RFC1123,
        time.RFC3339,
        // add more if needed
    }
    
    for _, format := range formats {
        t, err := time.Parse(format, timeStr)
        if err == nil {
            return t, nil
        }
    }
    
    return time.Time{}, fmt.Errorf("could not parse time: %s", timeStr)
}
