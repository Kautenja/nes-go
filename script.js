
ws = new WebSocket("ws://localhost:9090/screen/");

ws.onopen = function() {
    console.log("[onopen] connect ws uri.");
    ws.send(JSON.stringify({"Action" : "requireConnect"}));
}

ws.onmessage = function(e) {
    console.log("[onmessage] receive message.");
    var res = JSON.parse(e.data);
    $("#image").attr("src", "data:image/png;base64," + res["img64"]);
}

ws.onclose = function(e) {
    console.log("[onclose] connection closed (" + e.code + ")");
}

ws.onerror = function (e) {
    console.log("[onerror] error!");
}

// https://keycode.info
document.addEventListener('keydown', function(event) {
    switch (event.keyCode) {
        case "W".charCodeAt(0): alert('"W" was pressed'); break;
        case "A".charCodeAt(0): alert('"A" was pressed'); break;
        case "S".charCodeAt(0): alert('"S" was pressed'); break;
        case "D".charCodeAt(0): alert('"D" was pressed'); break;
        case "O".charCodeAt(0): alert('"O" was pressed'); break;
        case "P".charCodeAt(0): alert('"P" was pressed'); break;
        case 13:                alert('"*" was pressed'); break;
    }
});

// https://keycode.info
document.addEventListener('keyup', function(event) {
    switch (event.keyCode) {
        case "W".charCodeAt(0): alert('"W" was released'); break;
        case "A".charCodeAt(0): alert('"A" was released'); break;
        case "S".charCodeAt(0): alert('"S" was released'); break;
        case "D".charCodeAt(0): alert('"D" was released'); break;
        case "O".charCodeAt(0): alert('"O" was released'); break;
        case "P".charCodeAt(0): alert('"P" was released'); break;
        case 13:                alert('"*" was released'); break;
    }
});
