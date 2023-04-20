module player

import dl.loader

type T_player_message = fn(voidptr, &u8)

pub struct Player {
	handle voidptr

	player_message T_player_message
}

pub fn new(h voidptr) !&Player {
	mut libs := []string{}
	libs << "dragonfly.dll"

	mut lib := loader.get_or_create_dynamic_lib_loader(loader.DynamicLibLoaderConfig{paths: libs}) or {
		panic(err)
	}

	lib.open()!

	p := &Player{
		handle: h
		player_message: T_player_message(lib.get_sym("player_message")!)
	}
	return p
}

pub fn (mut p Player)message(msg string) {
	p.player_message(p.handle, msg.str)
}