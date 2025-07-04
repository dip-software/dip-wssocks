package wss

import (
	"context"
	log "github.com/genshen/wssocks/wss/logger"
	"io"
	"net/http"
	"nhooyr.io/websocket"
	"github.com/loafoe/caddy-token/keys"
)

type WebsocksServerConfig struct {
	EnableHttp       bool
	EnableConnKey    bool   // bale connection key
	ConnKey          string // connection key
	EnableStatusPage bool   // enable/disable status page
}

type ServerWS struct {
	config WebsocksServerConfig
	hc     *HubCollection
}

// return a a function handling websocket requests from the peer.
func NewServeWS(hc *HubCollection, config WebsocksServerConfig) *ServerWS {
	return &ServerWS{config: config, hc: hc}
}

func (s *ServerWS) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// check connection key
	apiKey := r.Header.Get("Key")
	verified, k, err :=  keys.VerifyAPIKey(apiKey, s.config.ConnKey)
	if err != nil {
		log.Error("error verifying API key:", err)
		w.WriteHeader(401)
		w.Write([]byte("Access denied!\n"))
		return
	}
	if !verified {
		log.Error("invalid API key")
		w.WriteHeader(401)
		w.Write([]byte("Access denied!\n"))
		return
	}
	if k.Project != "dip" {
		log.Error("invalid API key project:", k.Project)
		w.WriteHeader(401)
		w.Write([]byte("Access denied!\n"))
		return
	}
	// TODO: process key Scopes to limit upstream access

	wc, err := websocket.Accept(w, r, nil)
	if err != nil {
		log.Error(err)
		return
	}
	defer wc.Close(websocket.StatusNormalClosure, "the sky is falling")

	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	// negotiate version with client.
	if err := NegVersionServer(ctx, wc, s.config.EnableStatusPage); err != nil {
		return
	}

	hub := s.hc.NewHub(wc)
	defer s.hc.RemoveProxy(hub.id)
	defer hub.Close()
	// read messages from webSocket
	wc.SetReadLimit(1 << 23) // 8 MiB
	for {
		msgType, p, err := wc.Read(ctx) // fixme context
		// if WebSocket is closed by some reason, then this func will return,
		// and 'done' channel will be set, the outer func will reach to the end.
		if err != nil && err != io.EOF {
			log.Error("error reading webSocket message:", err)
			break
		}
		if err = dispatchMessage(hub, msgType, p, s.config); err != nil {
			log.Error("error proxy:", err)
			// break skip error
		}
	}
}
