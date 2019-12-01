package database

import (
	"log"
	// pq lib
	"github.com/gobuffalo/packr/v2"
	"github.com/henyxia/screenator/internal/model"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/rubenv/sql-migrate"
)

// Database handler
type Database struct {
	conn *sqlx.DB
}

// GetDisplay returns displays for a given mac
func (db Database) GetDisplay(mac string) string {
	query := db.conn.Rebind(`SELECT c.id, url
    FROM display di, device de, content c
    WHERE
        de.id = di.device AND
        c.id = di.content AND
        mac = ?
    ORDER BY start_time DESC
    LIMIT 1`)

	row := db.conn.QueryRowx(query, mac)
	var url string
	var contentID int
	err := row.Scan(&contentID, &url)
	if err != nil {
		log.Println("database select failed:", err)
	}

	return url
}

// GetContent returns a content from ID
func (db Database) GetContent(id int) model.Content {
	query := db.conn.Rebind("SELECT * FROM content WHERE id = ?")
	content := model.Content{}
	err := db.conn.Get(&content, query, id)
	if err != nil {
		log.Println("cannot get content:", err)
	}

	return content
}

// GetDeviceFromMac returns a device from a mac
func (db Database) GetDeviceFromMac(mac string) model.Device {
	query := db.conn.Rebind("SELECT * FROM device WHERE mac = ?")
	device := model.Device{}
	err := db.conn.Get(&device, query, mac)
	if err != nil {
		log.Println("cannot get device:", err)
	}

	return device
}

// GetDeviceDisplays returns displays for a device ID
func (db Database) GetDeviceDisplays(deviceID int) []model.Display {
	query := db.conn.Rebind("SELECT * FROM display WHERE device = ?")
	displays := []model.Display{}
	err := db.conn.Select(&displays, query, deviceID)
	if err != nil {
		log.Println("cannot get displays:", err)
	}

	return displays
}

// GetHeadersOfContent returns headers for a content ID
func (db Database) GetHeadersOfContent(contentID int) []model.Header {
	query := db.conn.Rebind("SELECT * FROM header WHERE content = ?")
	headers := []model.Header{}
	err := db.conn.Select(&headers, query, contentID)
	if err != nil {
		log.Println("cannot get headers:", err)
	}

	return headers
}

// Connect establish the first connection to the DB and check migrations
func Connect(driver string, conn string) Database {
	db, err := sqlx.Connect(driver, conn)
	if err != nil {
		log.Fatalln("unable to connect to the database!", err)
	}

	var dbHandler = Database{db}

	log.Println("apply migration if needed")
	box := packr.New("test", "../../sql")
	migrations := &migrate.PackrMigrationSource{
		Box: box,
	}
	log.Println("items in the box:", box.List())

	n, err := migrate.Exec(db.DB, driver, migrations, migrate.Up)
	if err != nil {
		log.Fatalln("cannot apply migrations", err)
	}
	log.Printf("applied %d migrations!\n", n)

	return dbHandler
}
