package service

import (
	"fmt"
	"skylight/internal/model/entity"

	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/glog"
)

type publishService struct {
	websockets map[string]*ghttp.WebSocket
}

// cluster
func (s publishService) RegisterPublisher(req *ghttp.Request) *ghttp.WebSocket {
	LOG := glog.Expose().Clone()
	LOG.SetPrefix(fmt.Sprintf("[%s]", req.GetSessionId()))
	ws, err := req.WebSocket()
	if err != nil {
		LOG.Errorf(req.GetCtx(), "get websocket of session '%s' failed", req.GetSessionId())
		req.Exit()
		return nil
	}
	if ws != nil {
		LOG.Infof(req.GetCtx(), "get websocket %s -> %s", ws.RemoteAddr(), ws.LocalAddr())
	} else {
		return nil
	}
	LOG.Infof(req.GetCtx(), "register websocket")
	s.websockets[req.GetSessionId()] = ws
	return ws
}

func (s publishService) Publish(req *ghttp.Request, message entity.Message) {
	defer func() {
		if r := recover(); r != nil {
			glog.Errorf(req.GetCtx(), "got panic: %v", r)
		}
	}()
	sessionId := req.GetSessionId()
	ws, ok := s.websockets[sessionId]
	if !ok || ws == nil {
		glog.Warningf(req.GetCtx(), "websocket not exists")
		// ws = s.RegisterPublisher(req)
	}
	if ws == nil {
		glog.Errorf(req.GetCtx(), "session websocket is none")
		return
	}
	glog.Infof(req.GetCtx(), "send message to session %s: %s", sessionId, message)
	if err := ws.WriteJSON(message); err != nil {
		glog.Errorf(req.GetCtx(), "write message failed")
		return
	}
}

var PublishService *publishService

func init() {
	PublishService = &publishService{
		websockets: map[string]*ghttp.WebSocket{},
	}
}
