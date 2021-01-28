package session

import (
	"go-app-template/src/apperror"
	"go-app-template/src/config"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

var _sessionKey = config.GetConfig()["session_key"].(string)

type SessionManager struct {
}

func NewSessionManager() SessionManager {
	return SessionManager{}
}

func (s SessionManager) Get(c echo.Context) *sessions.Session {
	sess, _ := session.Get(_sessionKey, c)
	return sess
}

func (s SessionManager) Set(c echo.Context, key string, value string) error {
	sess := s.Get(c)
	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400,
		HttpOnly: true,
	}
	sess.Values[key] = value
	if err := sess.Save(c.Request(), c.Response()); err != nil {
		return apperror.NewAppErrorWithStatus(err, http.StatusInternalServerError)
	}
	return nil
}

func (s SessionManager) Invalidate(c echo.Context) error {
	sess := s.Get(c)
	sess.Options.MaxAge = -1
	if err := sess.Save(c.Request(), c.Response()); err != nil {
		return apperror.NewAppErrorWithStatus(err, http.StatusInternalServerError)
	}
	return nil
}
