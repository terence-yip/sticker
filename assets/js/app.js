var moveleftButton = document.getElementById("move-left");
var moverightButton = document.getElementById("move-right");
var moveupButton = document.getElementById("move-up");
var movedownButton = document.getElementById("move-down");

function sendMoveCommand(direction, pressed) {
  var xmlHttp = new XMLHttpRequest();
  xmlHttp.open("POST", "/api", true);
  xmlHttp.setRequestHeader('Content-Type', 'application/json');
  xmlHttp.send(JSON.stringify({
    Direction: direction,
    Pressed: pressed,
  }));
}

moveleftButton.addEventListener("mousedown", function() {
  sendMoveCommand("left", "true");
});
moverightButton.addEventListener("mousedown", function() {
  sendMoveCommand("right", "true");
});
moveupButton.addEventListener("mousedown", function() {
  sendMoveCommand("up", "true");
});
movedownButton.addEventListener("mousedown", function() {
  sendMoveCommand("down", "true");
});

moveleftButton.addEventListener("mouseup", function() {
  sendMoveCommand("left", "false");
});
moverightButton.addEventListener("mouseup", function() {
  sendMoveCommand("right", "false");
});
moveupButton.addEventListener("mouseup", function() {
  sendMoveCommand("up", "false");
});
movedownButton.addEventListener("mouseup", function() {
  sendMoveCommand("down", "false");
});