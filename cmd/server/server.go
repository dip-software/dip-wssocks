package server

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/genshen/cmds"
	_ "github.com/genshen/wssocks/cmd/server/statik"
	"github.com/genshen/wssocks/wss"
	"github.com/genshen/wssocks/wss/status"
	"github.com/rakyll/statik/fs"
	log "github.com/sirupsen/logrus"
)

var serverCommand = &cmds.Command{
	Name:        "server",
	Summary:     "run as server mode",
	Description: "run as server program.",
	CustomFlags: false,
	HasOptions:  true,
}

func init() {
	var s server
	fs := flag.NewFlagSet("server", flag.ContinueOnError)
	serverCommand.FlagSet = fs
	serverCommand.FlagSet.StringVar(&s.address, "addr", ":1088", `listen address.`)
	serverCommand.FlagSet.StringVar(&s.wsBasePath, "ws_base_path", "/", "base path for serving websocket.")
	serverCommand.FlagSet.BoolVar(&s.http, "http", true, `enable http and https proxy.`)
	serverCommand.FlagSet.StringVar(&s.signKey, "sign-key", "", "signing key of API tokens.")
	serverCommand.FlagSet.StringVar(&s.signKeyFile, "sign-key-file", "", "file to read the signing key from.")
	serverCommand.FlagSet.BoolVar(&s.tls, "tls", false, "enable/disable HTTPS/TLS support of server.")
	serverCommand.FlagSet.StringVar(&s.tlsCertFile, "tls-cert-file", "", "path of certificate file if HTTPS/tls is enabled.")
	serverCommand.FlagSet.StringVar(&s.tlsKeyFile, "tls-key-file", "", "path of private key file if HTTPS/tls is enabled.")
	serverCommand.FlagSet.BoolVar(&s.status, "status", false, `enable/disable service status page.`)
	serverCommand.FlagSet.Usage = serverCommand.Usage // use default usage provided by cmds.Command.

	serverCommand.Runner = &s
	cmds.AllCommands = append(cmds.AllCommands, serverCommand)
}

type server struct {
	address     string
	wsBasePath  string // base path for serving websocket and status page
	http        bool   // enable http and https proxy
	signKey     string // the connection key if authentication is enabled
	tls         bool   // enable/disable HTTPS/tls support of server.
	tlsCertFile string // path of certificate file if HTTPS/tls is enabled.
	tlsKeyFile  string // path of private key file if HTTPS/tls is enabled.
	status      bool   // enable service status page
	signKeyFile string // file to read the signing key from
}

func (s *server) PreRun() error {
	// read signing key file if provided
	if s.signKey == "" && s.signKeyFile != "" {
		// read signing key from file
		data, err := os.ReadFile(s.signKeyFile)
		if err != nil {
			log.WithFields(log.Fields{
				"file": s.signKeyFile,
			}).Error("error reading signing key file:", err)
			return fmt.Errorf("error reading signing key file: %w", err)
		}
		s.signKey = strings.TrimSpace(string(data))
		log.WithFields(log.Fields{
			"file": s.signKeyFile,
		}).Info("signing key read from file.")
	}
	// check if signing key is provided
	if s.signKey == "" {
		log.Trace("empty singing key provided.")
		return fmt.Errorf("signing key is required, please provide it with `-sign_key` flag")
	}
	// set base url
	if s.wsBasePath == "" {
		s.wsBasePath = "/"
	}
	// complete prefix and suffix
	if !strings.HasPrefix(s.wsBasePath, "/") {
		s.wsBasePath = "/" + s.wsBasePath
	}
	if !strings.HasSuffix(s.wsBasePath, "/") {
		s.wsBasePath = s.wsBasePath + "/"
	}
	return nil
}

func (s *server) Run() error {
	config := wss.WebsocksServerConfig{EnableHttp: s.http, EnableConnKey: true, ConnKey: s.signKey, EnableStatusPage: s.status}
	hc := wss.NewHubCollection()

	http.Handle(s.wsBasePath, wss.NewServeWS(hc, config))
	if s.status {
		statikFS, err := fs.New()
		if err != nil {
			log.Fatal(err)
		}
		http.Handle("/status/", http.StripPrefix("/status", http.FileServer(statikFS)))
		http.Handle("/api/status/", status.NewStatusHandle(hc, s.http, true, s.wsBasePath))
	}

	if s.status {
		log.Info("service status page is enabled at `/status` endpoint")
	}

	listenAddrToLog := s.address + s.wsBasePath
	if s.wsBasePath == "/" {
		listenAddrToLog = s.address
	}
	log.WithFields(log.Fields{
		"listen address": listenAddrToLog,
	}).Info("listening for incoming messages.")

	if s.tls {
		log.Fatal(http.ListenAndServeTLS(s.address, s.tlsCertFile, s.tlsKeyFile, nil))
	} else {
		log.Fatal(http.ListenAndServe(s.address, nil))
	}
	return nil
}
