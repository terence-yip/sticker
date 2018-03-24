var moveleftButton = document.getElementById("move-left");
var moverightButton = document.getElementById("move-right");
var moveupButton = document.getElementById("move-up");
var movedownButton = document.getElementById("move-down");

var lookleftButton = document.getElementById("look-left");
var lookrightButton = document.getElementById("look-right");
var lookupButton = document.getElementById("look-up");
var lookdownButton = document.getElementById("look-down");

function sendMoveCommand(direction, pressed) {
  sendCommand("move", direction, pressed)
}

function sendLookCommand(direction, pressed) {
  sendCommand("look", direction, pressed)
}

function sendCommand(command, direction, pressed) {
  var xmlHttp = new XMLHttpRequest();
  xmlHttp.open("POST", "/api", true);
  xmlHttp.setRequestHeader('Content-Type', 'application/json');
  xmlHttp.send(JSON.stringify({
    Command: command,
    Direction: direction,
    Pressed: pressed,
  }));
}

moveleftButton.addEventListener("mousedown", function() {
  sendCommand("move", "left", "true");
});
moverightButton.addEventListener("mousedown", function() {
  sendCommand("move", "right", "true");
});
moveupButton.addEventListener("mousedown", function() {
  sendCommand("move", "up", "true");
});
movedownButton.addEventListener("mousedown", function() {
  sendCommand("move", "down", "true");
});

moveleftButton.addEventListener("mouseup", function() {
  sendCommand("move", "left", "false");
});
moverightButton.addEventListener("mouseup", function() {
  sendCommand("move", "right", "false");
});
moveupButton.addEventListener("mouseup", function() {
  sendCommand("move", "up", "false");
});
movedownButton.addEventListener("mouseup", function() {
  sendCommand("move", "down", "false");
});

lookleftButton.addEventListener("mousedown", function() {
  sendCommand("look", "left", "true");
});
lookrightButton.addEventListener("mousedown", function() {
  sendCommand("look", "right", "true");
});
lookupButton.addEventListener("mousedown", function() {
  sendCommand("look", "up", "true");
});
lookdownButton.addEventListener("mousedown", function() {
  sendCommand("look", "down", "true");
});

lookleftButton.addEventListener("mouseup", function() {
  sendCommand("look", "left", "false");
});
lookrightButton.addEventListener("mouseup", function() {
  sendCommand("look", "right", "false");
});
lookupButton.addEventListener("mouseup", function() {
  sendCommand("look", "up", "false");
});
lookdownButton.addEventListener("mouseup", function() {
  sendCommand("look", "down", "false");
});

moveleftButton.addEventListener("touchstart", function() {
  sendCommand("move", "left", "true");
});
moverightButton.addEventListener("touchstart", function() {
  sendCommand("move", "right", "true");
});
moveupButton.addEventListener("touchstart", function() {
  sendCommand("move", "up", "true");
});
movedownButton.addEventListener("touchstart", function() {
  sendCommand("move", "down", "true");
});

moveleftButton.addEventListener("touchend", function() {
  sendCommand("move", "left", "false");
});
moverightButton.addEventListener("touchend", function() {
  sendCommand("move", "right", "false");
});
moveupButton.addEventListener("touchend", function() {
  sendCommand("move", "up", "false");
});
movedownButton.addEventListener("touchend", function() {
  sendCommand("move", "down", "false");
});

lookleftButton.addEventListener("touchstart", function() {
  sendCommand("look", "left", "true");
});
lookrightButton.addEventListener("touchstart", function() {
  sendCommand("look", "right", "true");
});
lookupButton.addEventListener("touchstart", function() {
  sendCommand("look", "up", "true");
});
lookdownButton.addEventListener("touchstart", function() {
  sendCommand("look", "down", "true");
});

lookleftButton.addEventListener("touchend", function() {
  sendCommand("look", "left", "false");
});
lookrightButton.addEventListener("touchend", function() {
  sendCommand("look", "right", "false");
});
lookupButton.addEventListener("touchend", function() {
  sendCommand("look", "up", "false");
});
lookdownButton.addEventListener("touchend", function() {
  sendCommand("look", "down", "false");
});

var timeoutPeriod = 500;
var imageURI = '/image.jpg';
var img = new Image();
img.onload = function() {
    var canvas = document.getElementById("current-frame");
    var context = canvas.getContext("2d");

    context.drawImage(img, 0, 0);
    setTimeout(timedRefresh,timeoutPeriod);
};

function timedRefresh() {
    // just change src attribute, will always trigger the onload callback
    img.src = imageURI + '?d=' + Date.now();
}

timedRefresh();