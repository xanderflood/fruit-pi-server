package main

import (
	"context"
	"database/sql"
	"log"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	flag "github.com/jessevdk/go-flags"

	"github.com/xanderflood/fruit-pi-server/cmd/api/server"
	"github.com/xanderflood/fruit-pi-server/cmd/api/server/auth"
	"github.com/xanderflood/fruit-pi-server/internal/pkg/db"
	"github.com/xanderflood/fruit-pi-server/lib/tools"

	//postgres driver for db/sql
	_ "github.com/lib/pq"
)

var options struct {
	PostgresConnectionString string `long:"postgres-connection-string" env:"POSTGRES_CONNECTION_STRING"  required:"true"`
	JWTSigningSecret         string `long:"jwt-signing-secret"         env:"JWT_SIGNING_SECRET"          required:"true"`

	Port  string `long:"port"          env:"PORT" default:"8000"`
	Debug bool   `long:"debug"         env:"DEBUG"`
}

func main() {
	_, err := flag.Parse(&options)
	if err != nil {
		log.Fatal(err)
	}

	sqlDB, err := sql.Open("postgres", options.PostgresConnectionString)
	if err != nil {
		log.Fatalf("couldn't initialize database connection: %s", err.Error())
	}

	dbClient := db.NewDBAgent(sqlDB)
	if err = db.EnsureDatabase(context.Background(), dbClient); err != nil {
		log.Fatalf("couldn't initialize accounts table: %s", err.Error())
	}

	logger := tools.NewStdoutLogger()

	authMgr := auth.NewAuthorizationManager(
		logger,
		options.JWTSigningSecret,
		&jwt.Parser{ValidMethods: []string{"HS256"}},
		dbClient,
	)

	srv := server.NewServer(
		logger,
		authMgr,
		auth.GetAuthorizationFromContext,
		dbClient,
		options.JWTSigningSecret,
	)

	//build the gin server
	r := gin.Default()

	r.LoadHTMLFiles(
		"templates/index.tmpl",
		"templates/not_registered.tmpl",
		"templates/error_code.tmpl",
	)
	r.Static("/static", "./static")
	server.AddRoutes(r, srv)

	log.Fatal(r.Run(":" + options.Port))
}
