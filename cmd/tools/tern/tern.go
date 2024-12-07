package main

import (
	"log/slog"
	"os/exec"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	cmd := exec.Command(
		"tern",
		"migrate",
		"--migrations",
		"./internal/store/pg/migrations",
		"--config",
		"./internal/store/pg/migrations/tern.conf",
	)

	if err := cmd.Run(); err != nil {
		slog.Error("error running tern", "error", err)
		panic(err)
	}
}
