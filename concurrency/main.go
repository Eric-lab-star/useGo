package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"text/template"

	"golang.org/x/net/websocket"
)

const listenAddr = "localhost:4000"

func main() {
	http.HandleFunc("/", rootHandler)
	http.Handle("/socket", websocket.Handler(socketHandler))
	err := http.ListenAndServe(listenAddr, nil)
	if err != nil {
		log.Fatal(err)
	}

}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	rootTemplate.Execute(w, listenAddr)
}

var rootTemplate = template.Must(template.New("root").Parse(`
<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8" />
<body>
<form id="form">
<input id="msg" type="text" placeholder="your message"/>
</form>
</body>
<script>
	var form = document.getElementById("form")
	var websocket = new WebSocket("ws://{{.}}/socket");
	const input = document.querySelector("input")
	function onSubmit(event){
		event.preventDefault()
		const msg = input.value
		const h1 = document.createElement("h1");
		document.body.appendChild(h1);
		h1.innerText = "Me: " + msg
		websocket.send("Partner:" + msg)
		input.value = ""
	}
	form.addEventListener("submit", onSubmit)
    websocket.onmessage = function(m){
		var h1 = document.createElement("h1");
		document.body.appendChild(h1);
		h1.innerText = m.data;
	}
    websocket.onclose = onClose;
	function onClose(){
		var h1 = document.createElement("h1");
		document.body.appendChild(h1);
		h1.innerText = "user left";
	}
</script>
</html>
`))

type socket struct {
	io.ReadWriter
	done chan bool
}

func (s socket) Close() error {
	s.done <- true
	return nil
}

func socketHandler(ws *websocket.Conn) {
	s := socket{ws, make(chan bool)}
	go match(s)
	<-s.done
}

var partner = make(chan io.ReadWriteCloser)

func match(c io.ReadWriteCloser) {
	fmt.Fprint(c, "Waiting for a partner...")
	select {
	case partner <- c:
		// now handled by the other goroutine
	case p := <-partner:
		chat(p, c)
	}
}

func chat(a, b io.ReadWriteCloser) {
	fmt.Fprintln(a, "Found one! Say hi.")
	fmt.Fprintln(b, "Found one! Say hi.")
	errc := make(chan error, 1)
	go cp(a, b, errc)
	go cp(b, a, errc)
	if err := <-errc; err != nil {
		log.Println(err)
	}
	a.Close()
	b.Close()
}
func cp(w io.Writer, r io.Reader, errc chan<- error) {
	_, err := io.Copy(w, r)
	errc <- err
}
