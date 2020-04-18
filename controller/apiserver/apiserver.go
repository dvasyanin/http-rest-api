package apiserver

import (
	"context"
	"github.com/dvasyanin/http-rest-api/config"
	"github.com/dvasyanin/http-rest-api/repository/store/pgstore"
	"github.com/dvasyanin/http-rest-api/service"
	"github.com/gofiber/fiber"
	"github.com/jackc/pgx/v4"
	"log"
)

// Start service api rest
func Start(cfg *config.Config) error {
	// connect to db
	db, err := newDB(cfg.DatabaseURL())
	if err != nil {
		return err
	}
	defer db.Close(context.Background())

	// create server
	store := pgstore.New(db)
	src := service.NewService(store)
	srv := newServer(cfg.Context(), src)
	app := fiber.New()

	// chat methods ...
	app.Get("/chat/getByRespond", srv.getChatByRespond())
	app.Post("/chat/create", srv.createMessage())

	// clients methods ...
	// ...

	// run server
	if err := app.Listen(cfg.AppPort()); err != nil {
		log.Fatal("can't run rest server")
	}

	return nil
}

func newDB(databaseURL string) (*pgx.Conn, error) {
	db, err := pgx.Connect(context.Background(), databaseURL)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(context.Background()); err != nil {
		return nil, err
	}

	return db, nil
}
