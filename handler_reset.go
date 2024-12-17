package main

import (
	"context"
	"fmt"
)

func handlerReset(s *state, cmd command) error {

	err := s.db.ResetTable(context.Background())

	if err != nil {
		fmt.Println("Database couldn't be reset")
		return err
	}
	fmt.Println("Database reset succesfully")
	return nil
}
