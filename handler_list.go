package main

import (
	"context"
	"database/sql"
	"fmt"
)

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
