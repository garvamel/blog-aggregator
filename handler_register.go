package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/garvamel/blog-aggregator/internal/database"
	"github.com/google/uuid"
)

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
