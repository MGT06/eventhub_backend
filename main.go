package main

import (
	"log"
	"os"

	"github.com/MGT06/EventHub_Backend.git/internal/config"
	"github.com/MGT06/EventHub_Backend.git/internal/router"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

//	@title			Event Hub API
//	@version		1.0
//	@description	This is a sample server celler server.

//	@host		localhost:9000
//	@BasePath	/

//	@securityDefinitions.apikey		BearerToken
//	@in								header
//	@name							Authorization
//	@description					Bearer token used as identitiy for access resource

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println(err)
	}

	pdb := config.NewPsqlDb(os.Getenv("DBUSER"), os.Getenv("DBPASS"), os.Getenv("DBHOST"), os.Getenv("DBPORT"), os.Getenv("DBNAME"))

	pool, err := pdb.Connect()
	if err != nil {
		log.Println("Cannot Connect to DB\nReason: ", err.Error())
		return
	}

	defer pool.Close()


	rc := config.NewRedisClient(os.Getenv("RDBUSER"), os.Getenv("RDBPASS"), os.Getenv("RDBHOST"), os.Getenv("RDBPORT"))
	rdb := rc.ConnectRedis()

	// if err := rdb.Ping(ctx).Err(); err != nil {

	// }

	r := gin.Default()

	router.MainRouter(r, pool, rdb)

	r.Run(":9000")
}
