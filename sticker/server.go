package sticker

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Message struct {
	Command   string
	Direction string
	Pressed   string
}

func api(robot *Robot) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		var m Message
		err := decoder.Decode(&m)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		} else {
			update_robot(m, robot)
		}
	}
}

func update_robot(m Message, robot *Robot) {
	switch m.Command {
	case "move":
		update_motor(m, robot)
	case "look":
		update_servo_mode(m, robot)
	}
}

func update_motor(m Message, robot *Robot) {
	if m.Pressed == "true" {
		switch m.Direction {
		case "left":
			robot.turnLeft()
		case "right":
			robot.turnRight()
		case "up":
			robot.moveForward()
		case "down":
			robot.moveBackward()
		default:
			robot.stopMove()
		}
	} else if m.Pressed == "false" {
		robot.stopMove()
	}
}

func update_servo_mode(m Message, robot *Robot) {
	if m.Pressed == "true" {
		switch m.Direction {
		case "left":
			robot.lookLeft()
		case "right":
			robot.lookRight()
		case "up":
			robot.lookUp()
		case "down":
			robot.lookDown()
		default:
			robot.stopLook()
		}
	} else if m.Pressed == "false" {
		robot.stopLook()
	}
}

func RunServer(robot *Robot) {
	http.Handle("/", http.FileServer(http.Dir("assets/html/")))
	http.Handle("/css/", http.FileServer(http.Dir("assets/")))
	http.Handle("/js/", http.FileServer(http.Dir("assets/")))
	http.HandleFunc("/api", api(robot))
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("ListenAndServe: %v", err)
	}
}
