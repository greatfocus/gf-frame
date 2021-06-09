package server

import (
	"crypto/tls"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"time"

	gfbus "github.com/greatfocus/gf-bus"
	gfcache "github.com/greatfocus/gf-cache"
	gfcron "github.com/greatfocus/gf-cron"
	gfdispatcher "github.com/greatfocus/gf-dispatcher"
	"github.com/greatfocus/gf-frame/config"
	"github.com/greatfocus/gf-frame/database"
)

// HandlerFunc custom server handler
type HandlerFunc func(http.ResponseWriter, *http.Request)

// Meta struct
type Meta struct {
	Env        string
	Mux        *http.ServeMux
	Config     *config.Config
	DB         *database.Conn
	Cache      *gfcache.Cache
	Cron       *gfcron.Cron
	JWT        *JWT
	Dispatcher *gfdispatcher.Disp
	Bus        *gfbus.Bus
}

// Start the server
func (m *Meta) Start() {
	// setUploadPath creates an upload path
	m.setUploadPath()

	// serve creates server instance
	m.serve()
}

// setUploadPath creates an upload path
func (m *Meta) setUploadPath() {
	if m.Config.Server.UploadPath != "" {
		fs := http.FileServer(http.Dir(m.Config.Server.UploadPath + "/"))
		m.Mux.Handle("/file/", http.StripPrefix("/file/", fs))
	}
}

// serve creates server instance
func (m *Meta) serve() {
	addr := ":" + m.Config.Server.Port
	srv := &http.Server{
		Addr:           addr,
		ReadTimeout:    time.Duration(m.Config.Server.Timeout) * time.Second,
		WriteTimeout:   time.Duration(m.Config.Server.Timeout) * time.Second,
		MaxHeaderBytes: 1 << 20,
		Handler:        m.Mux,
		TLSConfig:      tlsConfig(m),
	}

	// create server connection
	if m.Config.Env == "prod" {
		log.Println("Listening to port secure HTTPS", addr)
		log.Fatal(srv.ListenAndServeTLS("", ""))
	} else {
		log.Println("Listening to port HTTP", addr)
		log.Fatal(srv.ListenAndServe())
	}
}

// tlsConfig update cert and key
func tlsConfig(m *Meta) *tls.Config {
	// load certificate file
	crt, err := ioutil.ReadFile(filepath.Clean(m.Config.Server.Secure.TLS.Cert))
	if err != nil {
		log.Fatalf("error reading CA certificate: %v", err)
	}

	// load certificate key
	ky, err := ioutil.ReadFile(filepath.Clean(m.Config.Server.Secure.TLS.Key))
	if err != nil {
		log.Fatalf("error reading CA certificate: %v", err)
	}

	cert, err := tls.X509KeyPair(crt, ky)
	if err != nil {
		log.Fatal(err)
	}

	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		ServerName:   m.Config.Server.Secure.TLS.Domain,
	}
}
