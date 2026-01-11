package main

import (
	"context"
	"fmt"
	"strings"

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
	testURL := "https://www.wagslane.dev/index.xml"
	rssPtr, err := fetchFeed(context.Background(), testURL)
	if err != nil {
		return fmt.Errorf("Error fetching RSS:\n%w", err)
	}
	fmt.Printf("Feed: %+v\n", rssPtr)
	return nil
}

func handlerAddFeed(s *state, cmd command) error {
	user, err := s.db.GetUserByName(context.Background(), s.cfg.CurrentUsername)
	if err != nil {
		return fmt.Errorf("Error getting user info from db:\n%w", err)
	}
	if len(cmd.args) < 2 {
		return fmt.Errorf("Not enough arguments.  Usage: gator %s <name_of_feed> <feed_url>", cmd.name)
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
	fmt.Println("Feed added:")
	fmt.Printf("ID: %s\n", feed.ID)
	fmt.Printf("Name: %s\n", feed.Name)
	fmt.Printf("URL: %s\n", feed.Url)
	fmt.Printf("UserID: %s\n", feed.UserID)
	fmt.Printf("Created At: %s\n", feed.CreatedAt)
	fmt.Printf("Updated At: %s\n", feed.UpdatedAt)
	return nil
}

func handlerListFeeds(s *state, cmd command) error {
	feeds, err := s.db.ListFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("Error getting feeds from db:\n%w", err)
	}
	for _, feed := range feeds {
		fmt.Printf("Feed name: %s\n", feed.Name)
		fmt.Printf("URL: %s\n", feed.Url)
		fmt.Printf("Added by: %s\n", feed.Name_2)
		fmt.Println()
	}
	return nil
}