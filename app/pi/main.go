package main

import (
	"../../sticker"
)

func main() {
	robot := sticker.NewRobot()
	robot.Start()
	sticker.RunServer(robot)
}
