package server

import (
	"github.com/gwwfps/assembly-lines/db"
	"github.com/gwwfps/assembly-lines/manager"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
	"gopkg.in/olahol/melody.v1"
)

type Server struct {
	db *db.DB
	m  *melody.Melody
	e  *echo.Echo
	gm *manager.GameManager
}

func NewServer(db *db.DB) *Server {
	m := melody.New()
	e := echo.New()
	return &Server{
		db: db,
		m:  m,
		gm: manager.NewGameManager(m, e.Logger),
		e:  e,
	}
}

func (srv *Server) Start() error {
	srv.e.Logger.SetLevel(log.DEBUG)

	srv.e.Use(middleware.Logger())
	srv.e.Use(middleware.Recover())

	srv.e.GET("/ws/:id", srv.handleUpgrade)

	srv.m.HandleMessage(srv.handleMessage)

	go manager.Watch(srv.gm)

	return srv.e.Start(":5555")
}
