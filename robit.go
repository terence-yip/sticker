package main

import (
        "fmt"
        "net/http"
        "encoding/json"
        "gobot.io/x/gobot/drivers/gpio"
        "gobot.io/x/gobot/platforms/raspi"
)

type Message struct {
        Direction string
        Pressed string
}

func update_robot(m Message, led *gpio.LedDriver) {
        if m.Pressed == "true" {
                led.On()
        } else if m.Pressed == "false" {
                led.Off()
        }
}

func api(led *gpio.LedDriver) func(w http.ResponseWriter, r *http.Request) {
        return func(w http.ResponseWriter, r *http.Request) {
                dec := json.NewDecoder(r.Body)
                var m Message
                err := dec.Decode(&m)
                if err != nil {
                        fmt.Printf("Error: %v\n", err)
                } else {
                        fmt.Printf("Received: %s, %s\n", m.Direction, m.Pressed)
                        update_robot(m, led)       
                }
        }
}

func main() {
        r := raspi.NewAdaptor()
        led := gpio.NewLedDriver(r, "36")

        http.Handle("/", http.FileServer(http.Dir("assets/html/")))
        http.Handle("/css/", http.FileServer(http.Dir("assets/")))
        http.Handle("/js/", http.FileServer(http.Dir("assets/")))
        http.HandleFunc("/api", api(led))
        err := http.ListenAndServe(":8080", nil)
        if err != nil {
            fmt.Printf("ListenAndServe: %v", err)
        }
}