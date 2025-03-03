package service

import (
	"context"
	"encoding/json"
	"skylight/utility"
	"sync"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gctx"
)

type EventData struct {
	Type    string    `json:"type"`
	Title   string    `json:"title"`
	Message any       `json:"message"`
	Time    time.Time `json:"time"`
}

type sseService struct {
	channels map[string]chan EventData
	mux      *sync.Mutex
}

func (s *sseService) getChannel(sessionId string) chan EventData {
	defer s.mux.Unlock()
	s.mux.Lock()
	return s.channels[sessionId]
}
func (s *sseService) setChannel(sessionId string, channel chan EventData) {
	defer s.mux.Unlock()
	s.mux.Lock()
	s.channels[sessionId] = channel
}
func (s *sseService) closeChannel(sessionId string) {
	defer s.mux.Unlock()
	s.mux.Lock()

	close(s.channels[sessionId])
	delete(s.channels, sessionId)
}

func (s *sseService) RegisterAndWait(sessionId string, req *ghttp.Request) {
	g.Log().Infof(req.GetCtx(), "register sse channel for session: %s", sessionId)

	if c := s.getChannel(sessionId); c != nil {
		s.closeChannel(sessionId)
	}
	s.setChannel(sessionId, make(chan EventData))
	req.Response.Header().Set("Content-Type", "text/event-stream")
	req.Response.Header().Set("Cache-Control", "no-cache")
	req.Response.Header().Set("Connection", "keep-alive")
	go func() {
		s.Send(sessionId, "info", "SSE连接成功", "")
	}()

	for event := range s.channels[sessionId] {
		data, _ := json.Marshal(event)
		g.Log().Info(req.GetCtx(), "send event to session", sessionId, string(data))
		req.Response.Writef("data: %s\n\n", data)
		req.Response.Flush()
	}
	g.Log().Infof(req.GetCtx(), "close sse request for session: %s", sessionId)
}
func (s *sseService) Unregister(sessionId string) {
	if _, ok := s.channels[sessionId]; !ok {
		return
	}
	g.Log().Infof(context.TODO(), "close sse channel for sessoin %s", sessionId)
	s.closeChannel(sessionId)
}

func (s *sseService) Send(sessionId string, eType, eTitle string, eMsg any) {
	defer func() {
		if r := recover(); r != nil {
			g.Log().Error(context.TODO(), "send event to %s failed, because: %s", sessionId, r)
		}
	}()

	channel := s.getChannel(sessionId)
	if channel == nil {
		g.Log().Error(context.TODO(), "SSE channel for session %s not exists", sessionId)
		return
	}
	channel <- EventData{eType, eTitle, eMsg, time.Now()}
}
func (s *sseService) countChannels() {
	g.Log().Infof(gctx.GetInitCtx(), "SSE channels count is %d", len(s.channels))
}

var SseService *sseService

func init() {
	SseService = &sseService{
		channels: map[string]chan EventData{},
		mux:      &sync.Mutex{},
	}
	go utility.StartPeriodTask(600, func() { SseService.countChannels() })
}
