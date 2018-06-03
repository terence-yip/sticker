package sticker

import (
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type RobotMessage struct {
	Command   string
	Direction string
	Pressed   string
}

type LoginMessage struct {
	Username string
	Password string
}

func home_page() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if logged_in(r) {
			display_control_page(w)
		} else {
			display_login_page(w)
		}
	}
}

func logged_in(r *http.Request) bool {
	cookie, err := r.Cookie("session")
	if err == nil {
		// err = bcrypt.CompareHashAndPassword([]byte("world"))
		return cookie.Value != ""
	}
	return false
}

func display_control_page(w http.ResponseWriter) {
	w.WriteHeader(http.StatusFound)
	w.Write([]byte(`
<html>
<head>
    <link href='/css/sticker.css' rel='stylesheet'>
</head>
<body>
  <h1>Sticker</h1>
  <h2>Movement</h2>
  <div>
      <div class="move-button" id="move-left">Left</div>
      <div class="move-button" id="move-right">Right</div>
      <div class="move-button" id="move-up">Up</div>
      <div class="move-button" id="move-down">Down</div>
  </div>
  <h2>Camera</h2>
  <div>
      <div class="look-button" id="look-left">Left</div>
      <div class="look-button" id="look-right">Right</div>
      <div class="look-button" id="look-up">Up</div>
      <div class="look-button" id="look-down">Down</div>
  </div>
  <canvas id="current-frame" width="400" height="300" />
  <script src="/js/app.js"></script>
</body>
</html>
`))
}

func display_login_page(w http.ResponseWriter) {
	w.WriteHeader(http.StatusFound)
	w.Write([]byte(`
<html>
<head>
    <link href='/css/sticker.css' rel='stylesheet'>
</head>
<body>
  <h1>Login</h1>
  <div>
  	<form action="/login" method="post">
  		<h2>Username</h2>
    	<input type="text" placeholder="Enter Username" name="username" required />
  		<h2>Password</h2>
    	<input type="password" placeholder="Enter Password" name="password" required />
    	<button type="submit" value="Login">Login</button>
  	</form>
  </div>
  <script src="/js/login.js"></script>
</body>
</html>
`))
}

func robot_api(robot *Robot) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		var m RobotMessage
		err := decoder.Decode(&m)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		} else {
			update_robot(m, robot)
		}
	}
}

func login_api() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		login(r.Form["username"][0], r.Form["password"][0], w)
		http.Redirect(w, r, "/", http.StatusPermanentRedirect)
	}
}

func login(username string, password string, w http.ResponseWriter) {
	if username == "hello" && password == "world" {
		cookie, err := session_cookie(username, password)
		if err == nil {
			http.SetCookie(w, cookie)
		} else {
			fmt.Printf("Failed to generate password")
		}
	}
}

func session_cookie(username string, password string) (*http.Cookie, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return empty_cookie(), err
	}
	cookie := &http.Cookie{
		Name:  "session",
		Value: string(hash),
		Path:  "/",
	}
	return cookie, nil
}

func empty_cookie() *http.Cookie {
	cookie := &http.Cookie{
		Name:  "session",
		Value: "",
		Path:  "/",
	}
	return cookie
}

func logout_api() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		logout(w)
	}
}

func logout(w http.ResponseWriter) {

}

func update_robot(m RobotMessage, robot *Robot) {
	switch m.Command {
	case "move":
		update_motor(m, robot)
	case "look":
		update_servo_mode(m, robot)
	}
}

func update_motor(m RobotMessage, robot *Robot) {
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

func update_servo_mode(m RobotMessage, robot *Robot) {
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
	http.HandleFunc("/", home_page())
	http.Handle("/css/", http.FileServer(http.Dir("assets/")))
	http.Handle("/js/", http.FileServer(http.Dir("assets/")))
	http.HandleFunc("/api", robot_api(robot))
	http.HandleFunc("/login", login_api())
	http.HandleFunc("/logout", logout_api())
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("ListenAndServe: %v", err)
	}
}
