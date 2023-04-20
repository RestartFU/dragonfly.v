package main

import (
	"C"
	"fmt"
	"github.com/df-mc/dragonfly/server/player"
	"os"
	"unsafe"

	"github.com/df-mc/dragonfly/server"
	"github.com/df-mc/dragonfly/server/player/chat"
	"github.com/pelletier/go-toml"
	"github.com/sirupsen/logrus"
)

/*
	PLAYER
*/

//export player_message
func player_message(p CPlayer, msg CString) {
	PlayerFromPtr(p).Message(C.GoString(msg))
}

//export player_name
func player_name(p CPlayer) CString {
	return C.CString(PlayerFromPtr(p).Name())
}

/*
	PLAYER HANDLER
*/

//export player_handle_quit
func player_handle_quit(p CPlayer, h func()) {
	hdl := PlayerFromPtr(p).Handler().(*handler)
	hdl.handleQuit = h
}

type handler struct {
	player.NopHandler

	handleQuit func()
}

func (h *handler) HandleQuit() {
	if h.handleQuit != nil {
		h.handleQuit()
	}
}

/*
	SERVER
*/

//export server_accept
func server_accept(srv CServer) (pl CPlayer) {
	ServerFromPtr(srv).Accept(func(p *player.Player) {
		p.Handle(&handler{})
		pl = uintptr(unsafe.Pointer(p))
	})
	return
}

//export server_players
func server_players(srv CServer) CArray {
	return GoArrayToCArray(ServerFromPtr(srv).Players())
}

//export server_start
func server_start() CServer {
	log := logrus.New()
	log.Formatter = &logrus.TextFormatter{ForceColors: true}
	log.Level = logrus.DebugLevel

	chat.Global.Subscribe(chat.StdoutSubscriber{})

	conf, err := readConfig(log)
	if err != nil {
		log.Fatalln(err)
	}

	srv := conf.New()
	srv.CloseOnProgramEnd()

	srv.Listen()
	return uintptr(unsafe.Pointer(srv))
}

/*
	UTILITIES
*/

// readConfig reads the configuration from the config.toml file, or creates the
// file if it does not yet exist.
func readConfig(log server.Logger) (server.Config, error) {
	c := server.DefaultConfig()
	var zero server.Config
	if _, err := os.Stat("config.toml"); os.IsNotExist(err) {
		data, err := toml.Marshal(c)
		if err != nil {
			return zero, fmt.Errorf("encode default config: %v", err)
		}
		if err := os.WriteFile("config.toml", data, 0644); err != nil {
			return zero, fmt.Errorf("create default config: %v", err)
		}
		return c.Config(log)
	}
	data, err := os.ReadFile("config.toml")
	if err != nil {
		return zero, fmt.Errorf("read config: %v", err)
	}
	if err := toml.Unmarshal(data, &c); err != nil {
		return zero, fmt.Errorf("decode config: %v", err)
	}
	return c.Config(log)
}
func main() {}
