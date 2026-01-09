package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
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
	if len(cmd.args) == 0 {
		return fmt.Errorf("Username required.  Usage: gator %s <name>", cmd.name)
	}
	username := strings.ToLower(cmd.args[0])
	user, err := s.db.GetUserByName(context.Background(), username)
	if err != nil {
		return fmt.Errorf("User '%s' not registered", username)
	}

	if err := s.cfg.SetUser(user.Name); err != nil {
		return fmt.Errorf("Couldn't set user: %w", err)
	}
	fmt.Printf("User %s set in config\n", username)
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("Username required.  Usage: gator %s <name>", cmd.name)
	}
	name := strings.ToLower(cmd.args[0])
	
	user, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID: uuid.New(),
		Name: name,
	})
	if err != nil {
		return fmt.Errorf("Couldn't create user: %w", err)
	}
	if err := s.cfg.SetUser(user.Name); err != nil {
		return fmt.Errorf("Couldn't set user: %w", err)
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
		return fmt.Errorf("Error reseting user table: %w", err)
	}
	fmt.Println("User table reset successfully")
	return nil
}

func handlerListUsers(s *state, cmd command) error {
	currentUser := s.cfg.CurrentUsername
	allUsers, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("Error retrieving users from db: %w", err)
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