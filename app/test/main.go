package main

import (
	"../../sticker"
)

func main() {
	robot := sticker.NewRobot()
	sticker.RunServer(robot)
}
