package server

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"

	jsonitor "github.com/json-iterator/go"
	"github.com/labstack/echo"
	"gopkg.in/olahol/melody.v1"
)

const playerIdKey = "playerId"

func (srv *Server) handleUpgrade(c echo.Context) error {
	playerId := c.Param("id")
	if playerId == "" {
		c.String(http.StatusBadRequest, "invalid identity")
		return nil
	}

	return srv.m.HandleRequestWithKeys(c.Response(), c.Request(), map[string]interface{}{
		playerIdKey: playerId,
	})
}

func (srv *Server) handleMessage(s *melody.Session, msg []byte) {
	defer func() {
		if r := recover(); r != nil {
			err, ok := r.(error)
			if !ok {
				err = fmt.Errorf("non-error panic: %v", r)
			}
			srv.e.Logger.Error("recovered panic in message handler", err.Error())
		}
	}()

	playerId := s.MustGet(playerIdKey).(string)

	parts := strings.SplitN(string(msg), "|", 2)
	action := parts[0]
	argsJson := ""
	if len(parts) > 1 {
		argsJson = parts[1]
	}
	srv.e.Logger.Debugf("playerId=%s, action=%s, argsJson=%s", playerId, action, argsJson)

	if strings.ToLower(action[:1]) == action[:1] {
		srv.e.Logger.Errorf("forbidden action received %s", action)
		return
	}

	method := reflect.ValueOf(srv.gm).MethodByName(action)
	if !method.IsValid() {
		srv.e.Logger.Errorf("invalid action received %s", action)
		return
	}

	args := []reflect.Value{
		reflect.ValueOf(playerId),
	}
	methodType := method.Type()
	if methodType.NumIn() > 1 {
		argsVal := reflect.New(methodType.In(1))
		err := jsonitor.ConfigFastest.Unmarshal([]byte(argsJson), argsVal.Interface())
		if err != nil {
			srv.e.Logger.Errorf("invalid argsJson %s", argsJson)
			return
		}
		args = append(args, argsVal.Elem())
	}

	out := method.Call(args)
	if len(out) > 0 {
		err := out[0].Interface().(error)
		if err != nil {
			s.Write([]byte(fmt.Sprintf("error|%s", err.Error())))
			srv.e.Logger.Error("message handler returned error", err)
		}
	}
}
