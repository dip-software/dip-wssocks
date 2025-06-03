package wss

import (
	"context"
	"encoding/base64"
	"time"

	"github.com/segmentio/ksuid"
	log "github.com/genshen/wssocks/wss/logger"
	"nhooyr.io/websocket/wsjson"
)

const (
	TagData = iota
	TagEstOk
	TagEstErr
	TagNoMore
)

// proxy client handle one connection, send data to proxy server vai websocket.
type ProxyClient struct {
	Id       ksuid.KSUID
	onData   func(ksuid.KSUID, ServerData) // data from server todo data with  type
	onClosed func(ksuid.KSUID, bool)       // close connection, param bool: do tellClose if true
	onError  func(ksuid.KSUID, error)      // if there are error messages
}

type ServerData struct {
	Tag  int
	Data []byte
}

// tell wssocks proxy server to establish a proxy connection by sending server 
// proxy address, type, initial data.
func (p *ProxyClient) Establish(wsc *WebSocketClient, firstSendData []byte, proxyType int, addr string) error {
	estMsg := ProxyEstMessage{
		Type:     proxyType,
		Addr:     addr,
		WithData: false,
	}
	if firstSendData != nil {
		estMsg.WithData = true
		estMsg.DataBase64 = base64.StdEncoding.EncodeToString(firstSendData)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	if err := wsjson.Write(ctx, wsc.WsConn, &WebSocketMessage{
		Type: WsTpEst,
		Id:   p.Id.String(),
		Data: estMsg,
	}); err != nil {
		log.Error("json error:", err)
		return err
	}
	return nil
}
