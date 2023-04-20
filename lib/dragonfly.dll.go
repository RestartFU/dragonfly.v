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
func player_message(pl uintptr, msg *C.char) {
	p := (*player.Player)(unsafe.Pointer(pl))
	p.Message(C.GoString(msg))
}

//export player_name
func player_name(pl uintptr) *C.char {
	p := (*player.Player)(unsafe.Pointer(pl))
	return C.CString(p.Name())
}

/*
	SERVER
*/

//export server_accept
func server_accept(srv uintptr) uintptr {
	s := (*server.Server)(unsafe.Pointer(srv))
	pl := uintptr(0)
	s.Accept(func(p *player.Player) { pl = uintptr(unsafe.Pointer(p)) })
	return pl
}

//export server_players
func server_players(srv uintptr) []uintptr {
	s := (*server.Server)(unsafe.Pointer(srv))
	p := make([]uintptr, 0)
	for _, pl := range s.Players() {
		p = append(p, uintptr(unsafe.Pointer(pl)))
	}
	return p
}

//export server_start
func server_start() uintptr {
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
