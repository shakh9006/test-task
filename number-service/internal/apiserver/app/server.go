package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-pg/migrations/v8"
	"github.com/go-pg/pg/v10"
	"github.com/rs/zerolog/log"
	"github.com/shakh9006/numbers-store/internal/apiserver/delivery/gin"
	"github.com/shakh9006/numbers-store/internal/apiserver/models"
	"github.com/shakh9006/numbers-store/internal/apiserver/services"
	"github.com/shakh9006/numbers-store/utils"
)

func DBConnection() (*pg.DB, error) {
	opts := &pg.Options{
		Addr:     fmt.Sprintf("%s:%s", utils.GetEnvVar("PG_HOST"), utils.GetEnvVar("PG_PORT")),
		User:     utils.GetEnvVar("PG_USER"),
		Password: utils.GetEnvVar("PG_USER"),
	}

	db := pg.Connect(opts)

	collection := migrations.NewCollection()
	err := collection.DiscoverSQLMigrations("migrations")
	if err != nil {
		return nil, err
	}

	_, _, err = collection.Run(db, "init")
	if err != nil {
		return nil, err
	}

	oldVersion, newVersion, err := collection.Run(db, "up")
	if err != nil {
		return nil, err
	}

	if newVersion != oldVersion {
		log.Printf("migrated from version %d to %d\n", oldVersion, newVersion)
	} else {
		log.Printf("version is %d\n", oldVersion)
	}

	return db, err
}

func SetupApp() *gin.Engine {
	log.Info().Msg("Initializing service")

	pgdb, err := DBConnection()
	if err != nil {
		log.Fatal().Err(err).Msg(err.Error())
	}

	numRepo := models.NewNumberRepository(pgdb)
	numberService := services.NewNumberRepository(numRepo)
	numberCtrl := delivery.NewNumberCtrl(numberService)
	router := delivery.NewGinRouter(numberCtrl)
	return router.GetRouter()
}
