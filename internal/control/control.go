package control

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/henyxia/screenator/internal/database"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

var db database.Database

func getDeviceFromMac(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	log.Println("search for device with mac:", vars["mac"])

	device := db.GetDeviceFromMac(vars["mac"])
	deviceJSON, err := json.Marshal(device)
	if err != nil {
		log.Println("cannot marshal content")
	}

	res.WriteHeader(http.StatusOK)
	fmt.Fprintf(res, "%s", deviceJSON)
}

func getDeviceDisplays(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	log.Println("search for displays for device:", vars["deviceID"])

	deviceID, err := strconv.Atoi(vars["deviceID"])
	if err != nil {
		log.Println("cannot convert deviceID to int")
	}

	displays := db.GetDeviceDisplays(deviceID)
	displaysJSON, err := json.Marshal(displays)
	if err != nil {
		log.Println("cannot marshal displays")
	}

	res.WriteHeader(http.StatusOK)
	fmt.Fprintf(res, "%s", displaysJSON)
}

func getContent(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	log.Println("search for content:", vars["contentID"])

	contentID, err := strconv.Atoi(vars["contentID"])
	if err != nil {
		log.Println("cannot convert contentID to int")
	}

	content := db.GetContent(contentID)
	contentJSON, err := json.Marshal(content)
	if err != nil {
		log.Println("cannot marshal content")
	}

	res.WriteHeader(http.StatusOK)
	fmt.Fprintf(res, "%s", contentJSON)
}

// Start starts the control endpoint
func Start(conn string, database database.Database, wg *sync.WaitGroup) {
	defer wg.Done()

	/* store the global db handler */
	db = database

	r := mux.NewRouter()
	r.HandleFunc("/device/searchByMac/{mac:(?:[0-9a-f]{2}:){5}[0-9a-f]{2}}", getDeviceFromMac)
	r.HandleFunc("/device/{deviceID:[0-9]+}/display", getDeviceDisplays)
	r.HandleFunc("/content/{contentID:[0-9]+}", getContent)
	http.Handle("/", r)

	srv := http.Server{
		Handler:      r,
		Addr:         conn,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Println("start control plane")
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatalln("cannot listen:", err)
	}
	log.Println("stop control plane")
}
