package data

import (
	"fmt"
	"github.com/henyxia/screenator/internal/database"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"sync"
)

var db database.Database

type dataHandler struct {
}

type spaHandler struct {
	staticPath string
	indexPath  string
}

func (m *dataHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	/* retrieve host from request */

	// FORMAT:
	// PROTO REMOTE - HOST - URI
	fmt.Printf(
		"DATA: %-7s %-21s - %-5s - %s\n",
		req.Method,
		req.RemoteAddr,
		req.Host,
		req.RequestURI,
	)

	contentID, err := strconv.Atoi(req.Host)
	if err != nil {
		log.Println("cannot convert host to content id")
	}
	content := db.GetContent(contentID)
	headers := db.GetHeadersOfContent(contentID)
	url, _ := url.Parse(content.URL)

	/* redirect once if necesary */
	if (req.URL.Path == "/" || req.URL.Path == "") && url.Path != "/" {
		http.Redirect(res, req, url.RequestURI(), 301)
	} else {
		url.Path = "/"
	}

	// create the reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(url)

	// Update the headers to allow for SSL redirection
	req.URL.Host = url.Host
	req.URL.Scheme = url.Scheme
	req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))

	for _, header := range headers {
		req.Header.Set(header.Key, header.Value)
	}
	req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
	req.Host = url.Host

	// Note that ServeHttp is non blocking and uses a go routine under the hood
	proxy.ServeHTTP(res, req)
}

// Start starts the data server
func Start(conn string, database database.Database, wg *sync.WaitGroup) {
	defer wg.Done()

	db = database

	log.Println("start listen data plane")
	err := http.ListenAndServe(conn, &dataHandler{})
	if err != nil {
		log.Fatalln("cannot listen:", err)
	}
	log.Println("stop data plane")
}
