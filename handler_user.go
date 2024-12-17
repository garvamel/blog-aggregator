package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/garvamel/blog-aggregator/internal/database"
	"github.com/google/uuid"
)

func handlerLogin(s *state, cmd command) error {
	if l := len(cmd.args); l < 1 {
		return fmt.Errorf("No username provided.")
	} else if l > 1 {
		return fmt.Errorf("Enter only one username.")
	}

	username := cmd.args[0]

	user, err := s.db.GetUser(context.Background(), username)
	if err == sql.ErrNoRows {
		fmt.Println(err)
		return err

	} else if err != nil {
		fmt.Println(err)
		return err
	} else {
		s.cfg.SetUser(user.Name)
	}

	fmt.Println("User has been set")
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if l := len(cmd.args); l < 1 {
		return fmt.Errorf("No username provided.")
	} else if l > 1 {
		return fmt.Errorf("Enter only one username.")
	}

	username := cmd.args[0]

	user, err := s.db.GetUser(context.Background(), username)
	if err == sql.ErrNoRows {
		user, err = s.db.CreateUser(context.Background(),
			database.CreateUserParams{
				ID:        uuid.New(),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				Name:      username,
			})
		if err != nil {
			fmt.Println(err)
			return err
		}

	} else if err != nil {
		fmt.Println(err)
		return err
	} else {
		return fmt.Errorf("User exists")
	}

	if user.Name != username {
		return fmt.Errorf("Database returned a different user")
	}

	s.cfg.SetUser(user.Name)

	fmt.Printf("User %s has been created and set as the current user\n", s.cfg.CurrentUserName)
	return nil
}

func handlerListUsers(s *state, cmd command) error {

	users, err := s.db.GetUsers(context.Background())
	if err == sql.ErrNoRows {
		fmt.Println("Database is emtpy")
		return err

	} else if err != nil {
		fmt.Println(err)
		return err
	}

	for _, user := range users {
		if user.Name == s.cfg.CurrentUserName {
			fmt.Printf("* %s (current)\n", user.Name)
		} else {
			fmt.Printf("* %s\n", user.Name)
		}
	}
	return nil
}
