package service

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/glog"
)

type sseService struct {
	channels map[string]chan EventData
	mux      *sync.Mutex
}

func (s *sseService) Register(sessionId string, req *ghttp.Request) {
	glog.Infof(req.GetCtx(), "register sse channel for session: %s", sessionId)
	s.mux.Lock()
	s.channels[sessionId] = make(chan EventData)
	s.mux.Unlock()
	req.Response.Header().Set("Content-Type", "text/event-stream")
	req.Response.Header().Set("Cache-Control", "no-cache")
	req.Response.Header().Set("Connection", "keep-alive")
	go func() {
		s.Send(sessionId, "info", "SSE连接成功", "")
	}()

	for {
		event := <-s.channels[sessionId]
		data, _ := json.Marshal(event)

		glog.Info(context.TODO(), "send event to session", sessionId, string(data))
		req.Response.Writef("data: %s\n\n", data)
		req.Response.Flush()
	}
}
func (s *sseService) Unregister(req *ghttp.Request) {
	if _, ok := s.channels[req.GetSessionId()]; !ok {
		return
	}
	delete(s.channels, req.GetSessionId())
}

type EventData struct {
	Type    string    `json:"type"`
	Title   string    `json:"title"`
	Message any       `json:"message"`
	Time    time.Time `json:"time"`
}

func (s *sseService) Send(sessionId string, eType, eTitle string, eMsg any) {
	channel, ok := s.channels[sessionId]
	if !ok {
		glog.Error(context.TODO(), "sse channel for session %s not exists", sessionId)
		return
	}
	channel <- EventData{eType, eTitle, eMsg, time.Now()}
}

var SseService *sseService

func init() {
	SseService = &sseService{
		channels: map[string]chan EventData{},
		mux:      &sync.Mutex{}}
}
