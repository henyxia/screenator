package proxylocal

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

// LocalProxy is the local proxy definition
type LocalProxy struct {
	Server              *http.Server
	Host                string
	Bind                string
	Proxy               string
	ProxyScheme         string
	ProxyTLSCertificate string
	ProxyTLSKey         string
}

func (p *LocalProxy) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	/* create the reverse proxy */
	url, _ := url.Parse(p.Proxy)
	proxy := httputil.NewSingleHostReverseProxy(url)

	if p.ProxyScheme == "https" {
		cert, err := tls.LoadX509KeyPair(p.ProxyTLSCertificate, p.ProxyTLSKey)
		if err != nil {
			log.Fatalf("server: loadkeys: %s", err)
		}
		config := tls.Config{Certificates: []tls.Certificate{cert}, InsecureSkipVerify: true}
		transport := &http.Transport{TLSClientConfig: &config}
		proxy.Transport = transport
	}

	/* log time! */
	fmt.Printf("DATA: %-7s %-21s - %-5s - %s\n", req.Method, req.RemoteAddr, p.Host, req.RequestURI)

	/* Update the headers to allow for SSL redirection */
	req.URL.Host = url.Host
	req.URL.Scheme = url.Scheme
	req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
	req.Host = p.Host

	/* gogo proxy */
	proxy.ServeHTTP(res, req)
}

// Start handle new connections for the local client
func (p *LocalProxy) Start() {
	log.Println("begin listen local proxy on " + p.Bind + " for host " + p.Host)
	p.Server = &http.Server{
		Addr:    p.Bind,
		Handler: p,
	}
	go func() {
		err := p.Server.ListenAndServe()
		if err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe(): %s", err)
		}
		log.Println("stop local proxy")
	}()
}
