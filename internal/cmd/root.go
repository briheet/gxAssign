package cmd

import (
	"context"

	"github.com/joho/godotenv"
)

func Execute(context context.Context) int {
	err := godotenv.Load()
	if err != nil {
		return 1
	}

	return 0
}
