package frame

import (
	"net/http"
	"time"

	gfcache "github.com/greatfocus/gf-cache"
	gfcron "github.com/greatfocus/gf-cron"
	gfdispatcher "github.com/greatfocus/gf-dispatcher"
	"github.com/greatfocus/gf-sframe/config"
	"github.com/greatfocus/gf-sframe/database"
	"github.com/greatfocus/gf-sframe/server"
	gfvalidator "github.com/greatfocus/gf-validator"
)

// Frame struct
type Frame struct {
	env    string
	server *server.Meta
}

// NewFrame get new instance of frame
func NewFrame(impl *config.Impl) *Frame {
	var f = &Frame{env: impl.Env}
	f.server = f.init(impl)
	return f
}

// Init provides a way to initialize the frame
func (f *Frame) init(impl *config.Impl) *server.Meta {

	// read the config file and prepare object
	config := f.initConfig(impl)

	// initCron creates instance of cron
	cron := f.initCron()

	// initCache creates instance of cache
	cache := f.initCache(config.Cache.DefaultExpiration, config.Cache.CleanupInterval)

	// initDB create database connection
	db := f.initDB(config, impl)

	// initCron creates instance of cron
	jwt := f.initJWT(config)

	dispatcher := f.initDispatcher(config)

	// Initiate validator
	gfvalidator.SetFieldsRequiredByDefault(true)

	return &server.Meta{
		Env:        impl.Env,
		Config:     config,
		Cron:       cron,
		Cache:      cache,
		DB:         db,
		JWT:        jwt,
		Dispatcher: dispatcher,
	}
}

// Start spins up the service
func (f *Frame) Start(mux *http.ServeMux) {
	f.server.Mux = mux
	f.server.Start()
}

// initConfig read the configuration file
func (f *Frame) initConfig(impl *config.Impl) *config.Config {
	var config = impl.GetConfig()
	return &config
}

// initCron creates instance of cron
func (f *Frame) initCron() *gfcron.Cron {
	return gfcron.New()
}

// initCache creates instance of cache
func (f *Frame) initCache(defaultExpiration, cleanupInterval int64) *gfcache.Cache {
	// Create a cache with a default expiration time of 5 minutes, and which
	// purges expired items every 10 minutes
	return gfcache.New(time.Duration(defaultExpiration), time.Duration(cleanupInterval))
}

// initDB read the configuration file
func (f *Frame) initDB(config *config.Config, impl *config.Impl) *database.Conn {
	// create database connection
	var db = database.Conn{}
	db.Init(config, impl)
	return &db
}

// initJWT creates instance of auth
func (f *Frame) initJWT(config *config.Config) *server.JWT {
	var jwt = server.JWT{}
	jwt.Init(config)
	return &jwt
}

// initDispatcher creates instance of dispatcher
func (f *Frame) initDispatcher(config *config.Config) *gfdispatcher.Disp {
	d := gfdispatcher.NewDispatcher(int(config.Server.Workers)).Start()
	return d
}
