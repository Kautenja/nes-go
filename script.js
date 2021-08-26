
// Create a web socket connection to the back-end game engine.
ws = new WebSocket("ws://localhost:9090/screen/");

// Create a placeholder for the controller bitmap
var controller = 0

ws.onopen = function() {
    console.log("[onopen] connect WebSocket URI.");
    ws.send(JSON.stringify({"Action" : "requireConnect"}));
}

ws.onmessage = function(message) {
    // console.log("[onmessage] received message.");
    var res = JSON.parse(message.data);
    $("#image").attr("src", "data:image/png;base64," + res["img64"]);
}

ws.onclose = function(message) {
    console.log("[onclose] connection closed (" + message.code + ")");
}

ws.onerror = function (message) {
    console.log("[onerror] error!");
}

// Add an event listener for detecting key-down events.
document.addEventListener('keydown', function(event) {
    switch (event.keyCode) {
        case "D".charCodeAt(0): controller |= 0b10000000; break;
        case "A".charCodeAt(0): controller |= 0b01000000; break;
        case "S".charCodeAt(0): controller |= 0b00100000; break;
        case "W".charCodeAt(0): controller |= 0b00010000; break;
        case 13:                controller |= 0b00001000; break;
        case 32:                controller |= 0b00000100; break;
        case "P".charCodeAt(0): controller |= 0b00000010; break;
        case "O".charCodeAt(0): controller |= 0b00000001; break;
    }
});

// Add an event listener for detecting key-up events.
document.addEventListener('keyup', function(event) {
    switch (event.keyCode) {
        case "D".charCodeAt(0): controller &= 0xff - 0b10000000; break;
        case "A".charCodeAt(0): controller &= 0xff - 0b01000000; break;
        case "S".charCodeAt(0): controller &= 0xff - 0b00100000; break;
        case "W".charCodeAt(0): controller &= 0xff - 0b00010000; break;
        case 13:                controller &= 0xff - 0b00001000; break;
        case 32:                controller &= 0xff - 0b00000100; break;
        case "P".charCodeAt(0): controller &= 0xff - 0b00000010; break;
        case "O".charCodeAt(0): controller &= 0xff - 0b00000001; break;
    }
});
