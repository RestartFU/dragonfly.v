# Example
```v
module main

import server

fn main() {
	mut s := server.new()!
	for {
		mut p, ok := s.accept()
		if !ok {
			return
		}
		p.message("welcome, ${p.name()}")
		mut player_list := []string{}
		for mut pl in s.players() {
			player_list << pl.name()
		}
		println("players online: ${player_list.join(", ")}")
	}
}
```

# Features
- Server methods:
  - `accept() (&player.Player, bool)` - accept incoming players.
  - `players() []&player.Player`      - returns an array of players that are on the server.
- Player methods:
  - `message(msg string)` - sends the player the given message.
  - `name() string`       - returns the player username
