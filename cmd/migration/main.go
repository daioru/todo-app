package main

import (
	"context"
	"time"

	"github.com/daioru/todo-app/internal/config"
	"github.com/daioru/todo-app/migrations"
	"github.com/daioru/todo-app/internal/pkg/db"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/rs/zerolog/log"
)

func main() {
	if err := config.ReadConfigYML("config.yml"); err != nil {
		log.Fatal().Err(err).Msg("Failed init configuration")
	}
	cfg := config.GetConfigInstance()

	conn, err := db.ConnectDB(&cfg.DB)
	if err != nil {
		log.Fatal().Err(err).Msg("sql.Open() error")
	}
	defer conn.Close()

	goose.SetBaseFS(migrations.EmbedFS)

	const cmd = "up"

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	err = goose.RunContext(ctx, cmd, conn.DB, ".")
	if err != nil {
		log.Fatal().Err(err).Msg("goose.Status() error")
	}
}
