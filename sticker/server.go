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

func homePage() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if loggedIn(r) {
			displayControlPage(w)
		} else {
			displayLoginPage(w)
		}
	}
}

func loggedIn(r *http.Request) bool {
	cookie, err := r.Cookie("session")
	if err == nil {
		// err = bcrypt.CompareHashAndPassword([]byte("world"))
		return cookie.Value != ""
	}
	return false
}

func displayControlPage(w http.ResponseWriter) {
	w.WriteHeader(http.StatusFound)
	w.Write([]byte(`
<html>
<head>
    <link href='/static/sticker.css' rel='stylesheet'>
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
  <script src="/static/app.js"></script>
</body>
</html>
`))
}

func displayLoginPage(w http.ResponseWriter) {
	w.WriteHeader(http.StatusFound)
	w.Write([]byte(`
<html>
<head>
    <link href='/static/sticker.css' rel='stylesheet'>
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
  <script src="/static/login.js"></script>
</body>
</html>
`))
}

func robotApi(robot *Robot) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		var m RobotMessage
		err := decoder.Decode(&m)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		} else {
			updateRobot(m, robot)
		}
	}
}

func loginApi() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		login(r.Form["username"][0], r.Form["password"][0], w)
		http.Redirect(w, r, "/", http.StatusPermanentRedirect)
	}
}

func login(username string, password string, w http.ResponseWriter) {
	if username == "hello" && password == "world" {
		cookie, err := sessionCookie(username, password)
		if err == nil {
			http.SetCookie(w, cookie)
		} else {
			fmt.Printf("Failed to generate password")
		}
	}
}

func sessionCookie(username string, password string) (*http.Cookie, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return emptyCookie(), err
	}
	cookie := &http.Cookie{
		Name:  "session",
		Value: string(hash),
		Path:  "/",
	}
	return cookie, nil
}

func emptyCookie() *http.Cookie {
	cookie := &http.Cookie{
		Name:  "session",
		Value: "",
		Path:  "/",
	}
	return cookie
}

func logoutApi() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		logout(w)
	}
}

func logout(w http.ResponseWriter) {

}

func updateRobot(m RobotMessage, robot *Robot) {
	switch m.Command {
	case "move":
		updateMotor(m, robot)
	case "look":
		updateServoMode(m, robot)
	}
}

func updateMotor(m RobotMessage, robot *Robot) {
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

func updateServoMode(m RobotMessage, robot *Robot) {
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
	http.HandleFunc("/", homePage())
	http.Handle("/static/", http.FileServer(http.Dir("./")))
	http.Handle("/images/", http.FileServer(http.Dir("./")))
	http.HandleFunc("/api", robotApi(robot))
	http.HandleFunc("/login", loginApi())
	http.HandleFunc("/logout", logoutApi())
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("ListenAndServe: %v", err)
	}
}
