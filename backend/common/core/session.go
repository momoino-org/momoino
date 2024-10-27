package core

import (
	"net/http"
	"time"

	"github.com/alexedwards/scs/gormstore"
	"github.com/alexedwards/scs/v2"
	"github.com/alexedwards/scs/v2/memstore"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

const UIDKey = "uid"

func UseGormStore(db *gorm.DB) func(sm *scs.SessionManager) error {
	return func(sm *scs.SessionManager) error {
		// Ensure the cleanup goroutine from `memstore` is stopped before switching to `gormstore`.
		// NOTE: Awaiting fix from author: https://github.com/alexedwards/scs/issues/222
		if memStore, ok := sm.Store.(*memstore.MemStore); ok {
			// HACK: Not sure why, but it seems a delay is needed before stopping the cleanup goroutine.
			time.Sleep(time.Second)
			memStore.StopCleanup()
		}

		var err error
		if sm.Store, err = gormstore.New(db); err != nil {
			return err
		}

		return nil
	}
}

func NewSessionManager(config AppConfig, opts ...func(*scs.SessionManager) error) (*scs.SessionManager, error) {
	sessionConfig := config.GetSessionConfig()

	sessionManager := scs.New()
	sessionManager.Lifetime = sessionConfig.LifeTime
	sessionManager.IdleTimeout = sessionConfig.IdleTimeout
	sessionManager.Cookie.Name = SessionCookie
	sessionManager.Cookie.SameSite = http.SameSiteStrictMode
	sessionManager.Cookie.Domain = config.GetHost()
	sessionManager.Cookie.Path = "/"
	sessionManager.Cookie.HttpOnly = true
	sessionManager.Cookie.Persist = true
	sessionManager.Cookie.Secure = config.IsHTTPS()

	for _, opt := range opts {
		if err := opt(sessionManager); err != nil {
			return nil, err
		}
	}

	return sessionManager, nil
}

func NewLoginSessionManager(config AppConfig) (*scs.SessionManager, error) {
	const loginSessionMaxAge = 5 * time.Minute

	sessionManager := scs.New()
	sessionManager.Lifetime = loginSessionMaxAge
	sessionManager.Cookie.Name = LoginSessionCookie
	sessionManager.Cookie.SameSite = http.SameSiteStrictMode
	sessionManager.Cookie.Domain = config.GetHost()
	sessionManager.Cookie.Path = "/"
	sessionManager.Cookie.HttpOnly = true
	sessionManager.Cookie.Persist = false
	sessionManager.Cookie.Secure = config.IsHTTPS()

	return sessionManager, nil
}

func NewSessionModule() fx.Option {
	return fx.Module(
		"Session module",
		fx.Provide(
			func(config AppConfig, db *gorm.DB) (*scs.SessionManager, error) {
				return NewSessionManager(config, UseGormStore(db))
			},
			fx.Annotate(
				NewLoginSessionManager,
				fx.ResultTags(`name:"loginSessionManager"`),
			),
		),
	)
}
