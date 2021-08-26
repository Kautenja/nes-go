
// Create a web socket connection to the back-end game engine.
ws = new WebSocket("ws://localhost:9090/screen/");

// Create a placeholder for the 8-bit controller bitmap. This value will be
// updated by the key-down and key-up event handlers when keys are pressed and
// released, respectively.
var controller = 0

ws.onopen = function() {
    console.log("[onopen] connect WebSocket URI.");
    ws.send(JSON.stringify({"Action" : "requireConnect"}));
}

ws.onmessage = function(message) {
    // Parse the base64 image from the JSON message.
    var data = JSON.parse(message.data);
    $("#image").attr("src", "data:image/png;base64," + data["img64"]);
    // Send the controller state to the engine.
    ws.send(JSON.stringify({"controller" : controller}));
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
