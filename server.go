package main

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessions"
	"github.com/gorilla/websocket"
	"fmt"
	"net/http"
	"bytes"
	"encoding/json"
	"sync"
	"log"
	"net"

)

const backendUrl string = "http://localhost:9000"


var ActiveClients = make(map[ClientConn] int)
var ActiveClientsRWMutex sync.RWMutex

type ClientConn struct {
	websocket *websocket.Conn
	clientIP  net.Addr
}

func addClient(cc ClientConn) {
	ActiveClientsRWMutex.Lock()
	ActiveClients[cc] = 0
	ActiveClientsRWMutex.Unlock()
}

func deleteClient(cc ClientConn) {
	ActiveClientsRWMutex.Lock()
	delete(ActiveClients, cc)
	ActiveClientsRWMutex.Unlock()
}

func broadcastMessage(messageType int, message []byte) {
	ActiveClientsRWMutex.RLock()
	defer ActiveClientsRWMutex.RUnlock()

	for client, _ := range ActiveClients {
		log.Println(client)
		if err := client.websocket.WriteMessage(messageType, message); err != nil {
			return
		}
	}

}


func main() {
	m := martini.Classic()

	store := sessions.NewCookieStore([]byte("rooftrellen"))
	m.Use(sessions.Sessions("frontend", store))

	m.Use(render.Renderer(render.Options{
		Layout: "layout",
	}))

	m.Get("/", RequireLogin, GetWelcome)
	m.Get("/login", GetLogin)
	m.Post("/login", PostLogin)
	m.Get("/logout", GetLogout)
	m.Get("/signup", GetSignup)
	m.Post("/signup", PostSignup)
	m.Get("/people", RequireLogin, GetPeople)

	m.Get("/chat", GetChat)

	m.Get("/chat/sock", func(w http.ResponseWriter, r *http.Request) {
			log.Println(ActiveClients)
			ws, err := websocket.Upgrade(w, r, nil, 1024, 1024)
			if _, ok := err.(websocket.HandshakeError); ok {
				http.Error(w, "Not a websocket handshake", 400)
				return
			} else if err != nil {
				log.Println(err)
				return
			}
			client := ws.RemoteAddr()
			sockCli := ClientConn{ws, client}
			addClient(sockCli)

			for {
				log.Println(len(ActiveClients), ActiveClients)
				messageType, p, err := ws.ReadMessage()
				log.Println("P: ", p)
				if err != nil {
					deleteClient(sockCli)
					log.Println("bye")
					log.Println(err)
					return
				}
				broadcastMessage(messageType, p)
			}
		})

	m.Run()
}

func RequireLogin(rw http.ResponseWriter, r *http.Request, s sessions.Session, c martini.Context) {
	userName := s.Get("userName")
	if userName == nil || GetUser(userName.(string)) == false {
		http.Redirect(rw, r, "/login", http.StatusFound)
		return
	}
	c.Map(userName)
	return
}

func GetUser(name string) bool {
	urlBack := fmt.Sprintf("%s/user/%s", backendUrl, name)
	reqBack, err := http.NewRequest("GET", urlBack, nil)
	if err != nil {
		panic(err)
	}
	clientBack := &http.Client{}
	resBack, err := clientBack.Do(reqBack)
	if err != nil {
		panic(err)
	}
	defer resBack.Body.Close()
	return resBack.StatusCode == 200
}

func GetWelcome(ren render.Render, s string) {
	ren.HTML(200, "welcome", s)
}

func GetLogin(ren render.Render, s sessions.Session) {
	s.Delete("userName")
	ren.HTML(200, "login", nil)
}

func PostUserAuthenticate(userName string, password string) int {
	jsonReqBack, _ := json.Marshal(map[string]string{"password": password})
	urlBack := fmt.Sprintf("%s/user/%s/authenticate", backendUrl, userName)
	reqBack, err := http.NewRequest("POST", urlBack, bytes.NewBuffer(jsonReqBack))
	if err != nil {
		panic(err)
	}
	reqBack.Header.Set("Content-Type", "application/json")
	clientBack := &http.Client{}
	resBack, err := clientBack.Do(reqBack)
	if err != nil {
		panic(err)
	}
	defer resBack.Body.Close()
	return resBack.StatusCode
}

func PostLogin(ren render.Render, r *http.Request, s sessions.Session) {
	userName := r.FormValue("username")
	statusCode := PostUserAuthenticate(userName, r.FormValue("password"))
	switch statusCode {
	case 401:
		ren.Redirect("/login")
	case 404:
		ren.Redirect("/login")
	default:
		s.Set("userName", userName)
		ren.Redirect("/")
	}
}

func GetLogout(ren render.Render, s sessions.Session) {
	s.Delete("userName")
	ren.Redirect("/login")
}

func GetSignup(ren render.Render, s sessions.Session) {
	s.Delete("userName")
	ren.HTML(200, "signup", nil)
}

func PostSignup(ren render.Render, r *http.Request, s sessions.Session) {
	userName := r.FormValue("username")
	password := r.FormValue("password")
	status := PostUserAuthenticate(userName, password)
	if (status == 401 || status == 404) {
		jsonReqBack, _ := json.Marshal(map[string]string{"userName": userName, "password": password, "createdAt": "2014-10-05 09:12"})
		urlBack := backendUrl + "/user/signup"
		reqBack, err := http.NewRequest("POST", urlBack, bytes.NewBuffer(jsonReqBack))
		if err != nil {
			panic(err)
		}
		reqBack.Header.Set("Content-Type", "application/json")
		clientBack := &http.Client{}
		resBack, err := clientBack.Do(reqBack)
		if err != nil {
			panic(err)
		}
		defer resBack.Body.Close()
		switch resBack.StatusCode {
		case 201:
			s.Set("userName", userName)
			ren.Redirect("/")
		default:
			ren.Redirect("/signup")
		}
	} else {
		s.Set("userName", userName)
		ren.Redirect("/")
	}
}

func GetPeople(ren render.Render) {
	ren.HTML(200, "people", nil)
}

func GetChat(ren render.Render) {
	ren.HTML(200, "chat", nil)
}


