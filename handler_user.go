package main

import (
	"context"
	"database/sql"
	"fmt"
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
