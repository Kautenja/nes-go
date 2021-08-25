package main

import (
    "github.com/gorilla/websocket"
    "net/http"
    "log"
    "fmt"
    "encoding/base64"
    "time"
    "bytes"
    "image"
    "image/png"
);

func main() {
    http.HandleFunc("/connws/", ConnWs)
    err := http.ListenAndServe(":9090", nil)
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}

func ConnWs(w http.ResponseWriter, r *http.Request) {
    ws, err := websocket.Upgrade(w, r, nil, 1024, 1024)
    _, ok := err.(websocket.HandshakeError)
    if ok {
        http.Error(w, "Not a websocket handshake", 400)
        return
    } else if err != nil {
        log.Println(err)
        return
    }

    emulator := load("smb1.nes")
    reset(emulator)

    res := map[string]interface{}{}
    	err = ws.ReadJSON(&res)
        if err != nil {
            if err.Error() == "EOF" {
                return
            }
            // ErrShortWrite means a write accepted fewer bytes than requested then failed to return an explicit error.
            if err.Error() == "unexpected EOF" {
                return
            }
            fmt.Println("Read : " + err.Error())
            return
        }

        res["a"] = "a"
        log.Println(res)

    for {
	    step(emulator)

	    img := image.NewRGBA(image.Rect(0, 0, screen_width(), screen_height()))
		img.Pix = pixels(emulator)
		buf := new(bytes.Buffer)
		png.Encode(buf, img)

        str := base64.StdEncoding.EncodeToString(buf.Bytes())
        res["img64"] = str

        err = ws.WriteJSON(&res)
        if err != nil {
            fmt.Println("watch dir - Write : " + err.Error())
        }
        time.Sleep(50 * time.Millisecond);
    }
}

