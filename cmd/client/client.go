package client

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/genshen/cmds"
	cl "github.com/genshen/wssocks/client"
	tools "github.com/genshen/wssocks/client/tools"
	"github.com/genshen/wssocks/cmd/client/health"
	log "github.com/sirupsen/logrus"
)

const CommandNameClient = "client"

var clientCommand = &cmds.Command{
	Name:        CommandNameClient,
	Summary:     "run as client mode",
	Description: "run as client program.",
	CustomFlags: false,
	HasOptions:  true,
}

type listFlags []string

func (l *listFlags) String() string {
	return "my string representation"
}

func (l *listFlags) Set(value string) error {
	*l = append(*l, value)
	return nil
}

func init() {
	var client client
	fs := flag.NewFlagSet(CommandNameClient, flag.ContinueOnError)
	clientCommand.FlagSet = fs
	clientCommand.FlagSet.StringVar(&client.address, "addr", ":1080", `listen address of socks5 proxy.`)
	clientCommand.FlagSet.BoolVar(&client.http, "http", false, `enable http and https proxy.`)
	clientCommand.FlagSet.StringVar(&client.httpAddr, "http-addr", ":1086", `listen address of http proxy (if enabled).`)
	clientCommand.FlagSet.StringVar(&client.remote, "remote", "", `server address and port(e.g: ws://example.com:1088).`)
	clientCommand.FlagSet.StringVar(&client.apiKey, "key", "", `connection key.`)
	clientCommand.FlagSet.StringVar(&client.apiKeyFile, "key-file", "", `File to read the key from.`)
	clientCommand.FlagSet.Var(&client.headers, "ws-header", `list of user defined http headers in websocket request. 
(e.g: --ws-header "X-Custom-Header=some-value" --ws-header "X-Second-Header=another-value")`)
	clientCommand.FlagSet.BoolVar(&client.skipTLSVerify, "skip-tls-verify", false, `skip verification of the server's certificate chain and host name.`)

	clientCommand.FlagSet.Usage = clientCommand.Usage // use default usage provided by cmds.Command.
	clientCommand.Runner = &client

	cmds.AllCommands = append(cmds.AllCommands, clientCommand)
}

type client struct {
	address       string      // local listening address
	http          bool        // enable http and https proxy
	httpAddr      string      // listen address of http and https(if it is enabled)
	remote        string      // string usr of server
	remoteUrl     *url.URL    // url of server
	headers       listFlags   // websocket headers passed from user.
	remoteHeaders http.Header // parsed websocket headers (not presented in flag).
	apiKey        string
	apiKeyFile    string
	endpoint      string
	skipTLSVerify bool
}

func (c *client) PreRun() error {
	if c.http {
		log.Info("http(s) proxy is enabled.")
	} else {
		log.Info("http(s) proxy is disabled.")
	}

	// read key from file if needed
	if c.apiKey == "" && c.apiKeyFile != "" {
		data, err := os.ReadFile(c.apiKeyFile)
		if err != nil {
			return fmt.Errorf("error reading key file %s: %w", c.apiKeyFile, err)
		}
		c.apiKey = strings.TrimSpace(string(data))
		if c.apiKey == "" {
			return fmt.Errorf("key file %s is empty", c.apiKeyFile)
		}
		log.WithFields(log.Fields{
			"key_file": c.apiKeyFile,
		}).Info("read connection key from file.")
	}
	// read from environment variable if needed
	if c.apiKey == "" {
		if envKey := os.Getenv("WSSOCKS_API_KEY"); envKey != "" {
			c.apiKey = envKey
			log.Info("read connection key from environment variable WSSOCKS_API_KEY.")
		}
	}
	// check key
	if c.apiKey == "" {
		return fmt.Errorf("connection key is required, please provide it with `-key` or `-key-file` flag")
	}
	remote, err := tools.GetGatewayFromAPIKey(c.apiKey)
	if err != nil {
		return fmt.Errorf("error getting remote url from api key: %w", err)
	}
	if c.remote == "" {
		log.WithField("remote", remote).Info("using remote from api key.")
		c.remote = remote
	}
	endpoint, err := tools.GetEndpointFromAPIKey(c.apiKey)
	if err != nil {
		return fmt.Errorf("error getting endpoint from api key: %w", err)
	}
	if endpoint != "" {
		log.WithField("endpoint", endpoint).Info("extracted endpoint from api key")
		c.endpoint = endpoint
	}

	// check remote address
	if c.remote == "" {
		return errors.New("empty remote address")
	}
	if u, err := url.Parse(c.remote); err != nil {
		return err
	} else {
		c.remoteUrl = u
	}

	// check header format.
	c.remoteHeaders = make(http.Header)
	for _, header := range c.headers {
		index := strings.IndexByte(header, '=')
		if index == -1 || index+1 == len(header) {
			return fmt.Errorf("bad http header in websocket request: %s", header)
		}
		hKey := ([]byte(header))[:index]
		hValue := ([]byte(header))[index+1:]
		c.remoteHeaders.Add(string(hKey), string(hValue))
	}

	return nil
}

func (c *client) Run() error {
	log.WithFields(log.Fields{
		"remote": c.remoteUrl.String(),
	}).Info("connecting to wssocks server.")

	options := cl.Options{
		LocalSocks5Addr: c.address,
		HttpEnabled:     c.http,
		LocalHttpAddr:   c.httpAddr,
		RemoteUrl:       c.remoteUrl,
		RemoteHeaders:   c.remoteHeaders,
		ConnectionKey:   c.apiKey,
		SkipTLSVerify:   c.skipTLSVerify,
		Endpoint:        c.endpoint,
	}
	hdl := cl.NewClientHandles()
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute) // fixme
	defer cancel()

	wsc, err := hdl.CreateServerConn(&options, ctx)
	if err != nil {
		return err
	}
	// server connect successfully
	log.WithFields(log.Fields{
		"remote": c.remoteUrl.String(),
	}).Info("connected to wssocks server.")
	defer wsc.WSClose()

	if err := hdl.NegotiateVersion(ctx, c.remote); err != nil {
		return err
	}
	// start health server
	health.StartServer(8091)

	var once sync.Once
	hdl.StartClient(&options, &once)
	hdl.CliWait(&once)
	return nil
}
