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
	"io/ioutil"
	"strconv"
)

const backendUrl string = "http://localhost:9000"

type User struct {
	Name string
	Status string
	Online string
}

var ActiveClients = make(map[ClientConn] int)
var ActiveClientsRWMutex sync.RWMutex

type ClientConn struct {
	websocket *websocket.Conn
	clientIP  net.Addr
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
	m.Get("/status", RequireLogin, GetStatus)

	m.Run()
}

func RequireLogin(ren render.Render, s sessions.Session, c martini.Context) {
	if user := GetUserBack(s.Get("userName")); user == nil {
		ren.Redirect("/login")
	} else {
		c.Map(user)
	}
	return
}

func GetUserBack(name interface {}) *User {
	if (name == nil) {
		return nil
	}
	url := fmt.Sprintf("%s/user/%s", backendUrl, name.(string))
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	if res.StatusCode == 200 {
		var data map[string]interface{}
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			panic(err)
		}
		json.Unmarshal(body, &data)
		user := User{}
		user.Name = data["userName"].(string)
		user.Status = fmt.Sprintf("%v", int(data["lastStatusCode"].(float64)))
		user.Online = "1"
		return &user
	} else {
		return nil
	}
}

func GetWelcome(ren render.Render, user *User) {
	ren.HTML(200, "welcome", user)
}

func GetLogin(ren render.Render, s sessions.Session) {
	s.Delete("userName")
	ren.HTML(200, "login", nil)
}

func PostLogin(ren render.Render, r *http.Request, s sessions.Session) {
	userName := r.FormValue("userName")
	if PostUserAuthenticate(userName, r.FormValue("password")) {
		s.Set("userName", userName)
		ren.Redirect("/")
	} else {
		ren.HTML(200, "login", "Wrong user name or password. Please try again.")
	}
}

func PostUserAuthenticate(userName string, password string) bool {
	body, _ := json.Marshal(map[string]string{"password": password})
	url := fmt.Sprintf("%s/user/%s/authenticate", backendUrl, userName)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	return res.StatusCode == 200
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
	userName := r.FormValue("userName")
	password := r.FormValue("password")
	if PostUserAuthenticate(userName, password) {
		s.Set("userName", userName)
		ren.Redirect("/")
	} else {
		body, _ := json.Marshal(map[string]string{"userName": userName, "password": password})
		url := backendUrl + "/user/signup"
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
		if err != nil {
			panic(err)
		}
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		res, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer res.Body.Close()
		switch res.StatusCode {
		case 201:
			s.Set("userName", userName)
			ren.Redirect("/")
		default:
			ren.HTML(200, "signup", "User name is taken. Please try another one.")
		}
	}
}

func GetPeople(ren render.Render, user *User) {
	ren.HTML(200, "people", user)
}

func GetStatus(ren render.Render, rw http.ResponseWriter, r *http.Request) {
	log.Println(ActiveClients)
	ws, err := websocket.Upgrade(rw, r, nil, 1024, 1024)
	if _, ok := err.(websocket.HandshakeError); ok {
		ren.Redirect("/people")
		return
	} else if err != nil {
		panic(err)
		return
	}
	client := ws.RemoteAddr()
	sockCli := ClientConn{ws, client}
	addClient(sockCli)
	var req map[string]interface{}
	for {
		log.Println(len(ActiveClients), ActiveClients)
		err := ws.ReadJSON(&req)
		log.Println("P: ", req)
		if err != nil {
			deleteClient(sockCli)
			log.Println("bye")
			log.Println(err)
			log.Println(len(ActiveClients), ActiveClients)
			return
		}
		switch (req["type"]) {
		case "status":
			res := map[string]interface{}{}
			res["type"] = "status"
			res["user"] = req["user"]
			broadcastMessage(&res)
			postReq := res["user"].(map[string]interface{})
			status, _ := strconv.Atoi(postReq["status"].(string))
			PostStatusBack(float64(status), postReq["name"].(string))
		case "people":
			res := map[string]interface{}{}
			res["type"] = "people"
			res["users"] = *GetUsersBack()
			ws.WriteJSON(&res)
		}
	}
}

func GetUsersBack() *[]map[string]interface{} {
	url := fmt.Sprintf("%s/users", backendUrl)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	data := []map[string]interface{}{}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	json.Unmarshal(body, &data)
	return &data
}

func PostStatusBack(status float64, name string) {

	s := map[string]interface{}{"createdAt":"2009-09-09 09:09", "statusCode":status, "location":"null"}
	body, _ := json.Marshal(s)
	url := backendUrl + "/status/" + name
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 201 {
		log.Print(res.StatusCode)
	}
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

func broadcastMessage(user interface{}) {
	ActiveClientsRWMutex.RLock()
	defer ActiveClientsRWMutex.RUnlock()
	for client, _ := range ActiveClients {
		if err := client.websocket.WriteJSON(&user); err != nil {
			panic(err)
		}
	}
}
