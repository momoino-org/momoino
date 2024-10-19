package core

import (
	"fmt"
	"net/http"
	"time"

	"github.com/alexedwards/scs/gormstore"
	"github.com/alexedwards/scs/v2"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

const UIDKey = "uid"

func newSessionManager(config AppConfig, db *gorm.DB) (*scs.SessionManager, error) {
	sessionConfig := config.GetSessionConfig()

	sessionManager := scs.New()
	sessionManager.Lifetime = sessionConfig.LifeTime
	sessionManager.IdleTimeout = sessionConfig.IdleTimeout
	sessionManager.Cookie.Name = SessionCookie
	sessionManager.Cookie.SameSite = http.SameSiteStrictMode
	sessionManager.Cookie.Domain = fmt.Sprintf(".%v", config.GetHost())
	sessionManager.Cookie.Path = "/"
	sessionManager.Cookie.HttpOnly = true
	sessionManager.Cookie.Persist = true
	sessionManager.Cookie.Secure = config.IsHTTPS()

	var err error
	if sessionManager.Store, err = gormstore.New(db); err != nil {
		return nil, err
	}

	return sessionManager, nil
}

func newShortLivedSessionManager(config AppConfig, db *gorm.DB) (*scs.SessionManager, error) {
	const shortLivedSessionMaxAge = 5 * time.Minute

	sessionManager := scs.New()
	sessionManager.Lifetime = shortLivedSessionMaxAge
	sessionManager.Cookie.Name = LoginSessionCookie
	sessionManager.Cookie.SameSite = http.SameSiteStrictMode
	sessionManager.Cookie.Domain = fmt.Sprintf(".%v", config.GetHost())
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
			newSessionManager,
			fx.Annotate(
				newShortLivedSessionManager,
				fx.ResultTags(`name:"shortLivedSessionManager"`),
			),
		),
	)
}
