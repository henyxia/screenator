package server

import (
	"github.com/henyxia/screenator/internal/config"
	"github.com/henyxia/screenator/internal/control"
	"github.com/henyxia/screenator/internal/data"
	"github.com/henyxia/screenator/internal/database"
	"log"
	"sync"
)

// Server starts the screenator server
func Server(configFile string) {
	log.Println("start screenator server")

	var wg sync.WaitGroup

	log.Println("read configuration file")
	conf := config.ReadConfigServer(configFile)

	log.Println("connect to database")
	db := database.Connect(conf.DbDriver, conf.DbConn)

	log.Println("start control server")
	wg.Add(1)
	go control.Start(conf.ControlBind, db, &wg)

	log.Println("start data")
	wg.Add(1)
	go data.Start(conf.DataBind, db, &wg)

	log.Println("server started")
	wg.Wait()
	log.Println("bye")
}
