package main

import (
	"context"
	"fmt"
	"strings"
	"time"
	"strconv"

	"github.com/thomas-reed/gator/internal/database"
)

type command struct {
	name string
	args []string
}

type commands struct {
	registry map[string]func(s *state, cmd command) error
}

func (c *commands) run(s *state, cmd command) error {
	cmdFunc, found := c.registry[cmd.name]
	if !found {
		return fmt.Errorf("%s not found in list of commands", cmd.name)
	}
	return cmdFunc(s, cmd)
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.registry[name] = f
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("Username required.  Usage: gator %s <name>", cmd.name)
	}
	username := strings.ToLower(cmd.args[0])
	user, err := s.db.GetUserByName(context.Background(), username)
	if err != nil {
		return fmt.Errorf("User '%s' not registered", username)
	}

	if err := s.cfg.SetUser(user.Name); err != nil {
		return fmt.Errorf("Couldn't set user:\n%w", err)
	}
	fmt.Printf("User %s set in config\n", username)
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("Username required.  Usage: gator %s <name>", cmd.name)
	}
	name := strings.ToLower(cmd.args[0])
	
	user, err := s.db.CreateUser(context.Background(), name)
	if err != nil {
		return fmt.Errorf("Couldn't create user:\n%w", err)
	}
	if err := s.cfg.SetUser(user.Name); err != nil {
		return fmt.Errorf("Couldn't set user:\n%w", err)
	}
	fmt.Println("User created:")
	fmt.Printf("Name: %s\n", user.Name)
	fmt.Printf("ID: %s\n", user.ID)
	fmt.Printf("Created At: %s\n", user.CreatedAt)
	fmt.Printf("Updated At: %s\n", user.UpdatedAt)
	return nil
}

func handlerReset(s *state, cmd command) error {
	if err := s.db.ResetUsers(context.Background()); err != nil {
		return fmt.Errorf("Error reseting user table:\n%w", err)
	}
	fmt.Println("User table reset successfully")
	return nil
}

func handlerListUsers(s *state, cmd command) error {
	currentUser := s.cfg.CurrentUsername
	allUsers, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("Error retrieving users from db:\n%w", err)
	}
	for _, user := range allUsers {
		if user == currentUser {
			fmt.Printf("* %s (current)\n", user)
		} else {
			fmt.Printf("* %s\n", user)
		}
	}
	return nil
}

func handlerAggregate(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("Poll interval required (e.g. 10s, 5m, 1h, etc.).  Usage: gator %s <poll_interval>", cmd.name)
	}
	timeBetweenRequests, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return fmt.Errorf("Error parsing time period:\n%w", err)
	}
	fmt.Printf("Collecting feeds every %s...\n", timeBetweenRequests)
	ticker := time.NewTicker(timeBetweenRequests)
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 2 {
		return fmt.Errorf("Feed name and URL required. Usage: gator %s <name_of_feed> <feed_url>", cmd.name)
	}
	name := cmd.args[0]
	url := cmd.args[1]
	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		Name: name,
		Url: url,
		UserID: user.ID,
	})
	if err != nil {
		return fmt.Errorf("Error adding feed to db:\n%w", err)
	}
	feedFollow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		return fmt.Errorf("Error creating feed follow:\n%w", err)
	}
	fmt.Println("Feed added:")
	fmt.Printf("FeedID: %s\n", feed.ID)
	fmt.Printf("FeedFollowID: %s\n", feedFollow.ID)
	fmt.Printf("Name: %s\n", feed.Name)
	fmt.Printf("URL: %s\n", feed.Url)
	fmt.Printf("UserID: %s\n", feed.UserID)
	fmt.Printf("User: %s\n", feedFollow.UserName)
	fmt.Printf("Created At: %s\n", feed.CreatedAt)
	fmt.Printf("Updated At: %s\n", feed.UpdatedAt)
	return nil
}

func handlerListFeeds(s *state, cmd command) error {
	feeds, err := s.db.ListFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("Error getting feeds from db:\n%w", err)
	}

	if len(feeds) == 0 {
		fmt.Println("No feeds found.")
		return nil
	}
	
	for _, feed := range feeds {
		fmt.Printf("Feed name: %s\n", feed.FeedName)
		fmt.Printf("URL: %s\n", feed.FeedUrl)
		fmt.Printf("Added by: %s\n", feed.UserName)
		fmt.Println()
	}
	return nil
}

func handlerFollow (s *state, cmd command, user database.User) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("Feed URL required.  Usage: gator %s <feed_url>", cmd.name)
	}
	url := cmd.args[0]

	feed, err := s.db.GetFeedByURL(context.Background(), url)
	if err != nil {
		return fmt.Errorf("Error getting user info from db:\n%w", err)
	}

	feedFollow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		return fmt.Errorf("Error creating feed follow:\n%w", err)
	}

	fmt.Println("Feed followed:")
	fmt.Printf("ID: %s\n", feedFollow.ID)
	fmt.Printf("Name: %s\n", feedFollow.FeedName)
	fmt.Printf("Username: %s\n", feedFollow.UserName)
	fmt.Printf("Created At: %s\n", feedFollow.CreatedAt)
	fmt.Printf("Updated At: %s\n", feedFollow.UpdatedAt)
	return nil
}

func handlerFollowing (s *state, cmd command, user database.User) error {
	follows, err := s.db.GetFeedFollowsByUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("Error getting user follows:\n%w", err)
	}

	if len(follows) == 0 {
		fmt.Println("You are not following any feeds.")
		return nil
	}

	fmt.Println("You follow:")
	for _, follows := range follows {
		fmt.Printf(" * %s\n", follows.FeedName)
	}
	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("Feed URL required.  Usage: gator %s <feed_url>", cmd.name)
	}
	feed, err := s.db.GetFeedByURL(context.Background(), cmd.args[0])
	if err != nil {
		return fmt.Errorf("Error getting feed with the given URL:\n%w", err)
	}
	if err = s.db.DeleteFeedFollowByUserAndName(context.Background(), database.DeleteFeedFollowByUserAndNameParams{
		UserID: user.ID,
		FeedID: feed.ID,
	}); err != nil {
		return fmt.Errorf("Error deleting feed follow:\n%w", err)
	}
	fmt.Printf("You are no longer following %s", feed.Name)
	return nil
}

func handlerBrowse(s *state, cmd command, user database.User) error {
	postLimit := 2
	if len(cmd.args) >= 1 {
		limit, err := strconv.Atoi(cmd.args[0])
		if err != nil || limit < 1 {
			fmt.Printf("Invalid post limit - defaulting to %d.  Usage:  gator %s <post_limit>\n", postLimit, cmd.name)
		} else {
			postLimit = limit
		}
	}

	posts, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit: int32(postLimit),
	})
	if err != nil {
		return fmt.Errorf("Error getting posts from db:\n%w", err)
	}

	fmt.Printf("Displaying most recent %d posts:\n", postLimit)
	for _, post := range posts {
		fmt.Printf("%s\n", post.Title)
		fmt.Printf("Link: %s\n", post.Url)
		if post.PublishedAt.Valid {
			fmt.Printf("Published: %v\n", post.PublishedAt.Time)
		} else {
			fmt.Println("Published: unknown")
		}
		fmt.Printf("By: %s\n", post.FeedName)
		fmt.Println("----------------------------------------")
		if post.Description.Valid {
			fmt.Printf("%s\n", post.Description.String)
		} else {
			fmt.Println("No description available")
		}
		fmt.Println("========================================")
		fmt.Println()
	}
	return nil
}