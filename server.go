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
	"time"
)

const backendUrl string = "http://localhost:9000"

type User struct {
	Name string
	Status string
	Online string
}

type Message struct {
	MessageID int
	Content string
	Author Author
	MessageType int
	PostedAt string
}

type Author struct {
	UserName string
	LastStatusCode int
}

var People map[string]User
var Wall []map[string]interface{}

var ActiveClients = make(map[ClientConn] int)
var ActiveClientsRWMutex sync.RWMutex

type ClientConn struct {
	websocket *websocket.Conn
	clientIP  net.Addr
}

func main() {
	m := martini.Classic()

	UpdatePeopleList()
	UpdateWall()

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
	m.Get("/socket", RequireLogin, GetSocket)
	m.Get("/wall", RequireLogin, GetWall)
	m.Post("/message", RequireLogin, PostMessage)

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
		user := User{People[userName].Name, People[userName].Status, "1"}
		People[userName] = user
		res := map[string]interface{}{}
		res["type"] = "login"
		res["user"] = user
		broadcastMessage(&res)

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
	userName := s.Get("userName").(string)
	s.Delete("userName")
	user := User{People[userName].Name, People[userName].Status, "0"}
	People[userName] = user
	res := map[string]interface{}{}
	res["type"] = "logout"
	res["user"] = user
	broadcastMessage(&res)

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
		user := User{People[userName].Name, People[userName].Status, "1"}
		People[userName] = user
		res := map[string]interface{}{}
		res["type"] = "login"
		res["user"] = user
		broadcastMessage(&res)

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
			user := User{userName, "0", "1"}
			People[userName] = user
			res := map[string]interface{}{}
			res["type"] = "signup"
			res["user"] = user
			broadcastMessage(&res)

			ren.Redirect("/")
		default:
			ren.HTML(200, "signup", "User name is taken. Please try another one.")
		}
	}
}

func PostMessage(ren render.Render, r *http.Request, s sessions.Session) {
	message := r.FormValue("message")
	userName := s.Get("userName")
	jsonString := map[string]interface{}{"content": message}
	body, _ := json.Marshal(jsonString)
	url := backendUrl + "/message/" + userName.(string)
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
	log.Print(res.StatusCode)

	if res.StatusCode == 201 {
		ren.Redirect("/wall")
	}

}

func GetPeople(ren render.Render, user *User) {
	ren.HTML(200, "people", user)
}

func GetSocket(ren render.Render, rw http.ResponseWriter, r *http.Request) {
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
			user := req["user"].(map[string]interface{})
			temp := User{user["name"].(string), user["status"].(string), "1"}
			People[user["name"].(string)] = temp
			postReq := res["user"].(map[string]interface{})
			status, _ := strconv.Atoi(postReq["status"].(string))
			PostStatusBack(float64(status), postReq["name"].(string))

			temp2 := map[string]interface{}{}
			temp2["type"] = "chat"
			temp2["message"] = map[string]interface{}{}
			temp3 := temp2["message"].(map[string]interface{})
			temp3["sender"] = user["name"].(string)
			temp3["receiver"] = "public"
			temp3["time"] = time.Now().Format("2006-01-02 15:04")
			s := postReq["status"].(string)
			switch s {
			case "0":
				s = "<font color=\"green\">I'm OK</font>"
				case "1":
			s="<font color=\"yellow\">I need help</font>"
			case "2":
			s="<font color=\"red\">Emergency</font>"
			}
			temp3["content"] = "Just change the status to " + s

			broadcastMessage(&temp2)
		case "people":
			res := map[string]interface{}{}
			res["type"] = "people"
			res["users"] = People
			ws.WriteJSON(&res)
		case "wall":
			res := map[string]interface{}{}
			res["type"] = "wall"
			res["messages"] = Wall
			ws.WriteJSON(&res)
		case "chat":
			temp := req["message"].(map[string]interface{})
			temp["time"] = time.Now().Format("2006-01-02 15:04")
			if (len(Wall) == 10) {
				Wall = append(Wall[1:10], temp)
			} else {
				Wall = append(Wall, temp)
			}
			broadcastMessage(&req)
		}
	}
}

func UpdatePeopleList() {
	People = make(map[string]User)
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
	for _, item := range data {
		name := item["userName"].(string)
		People[name] = User{name, fmt.Sprintf("%v", int(item["lastStatusCode"].(float64))), "0"}
	}
	return
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
	}
}

func GetWall(ren render.Render, user *User) {
	url := backendUrl + "/messages/wall"
	req, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer req.Body.Close()
	decoder := json.NewDecoder(req.Body)
	var data []Message
	errDecode := decoder.Decode(&data)
	if errDecode != nil {
		panic(errDecode)
	}

	log.Println(data)

	ren.HTML(200, "wall", map[string]interface {}{"user": user, "posts":data})
}

func UpdateWall() {
	message := map[string]interface{}{"sender": "xxx", "receiver": "public", "content": "fdafafa", "time": "20:44"}
	Wall = append(Wall, message)
	message = map[string]interface{}{"sender": "yyy", "receiver": "public", "content": "adsfadf", "time": "22:22"}
	Wall = append(Wall, message)
	message = map[string]interface{}{"sender": "zzz", "receiver": "public", "content": "rfdasfev", "time": "23:23"}
	Wall = append(Wall, message)
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

func broadcastMessage(message interface{}) {
	ActiveClientsRWMutex.RLock()
	defer ActiveClientsRWMutex.RUnlock()
	for client, _ := range ActiveClients {
		if err := client.websocket.WriteJSON(&message); err != nil {
			panic(err)
		}
	}
}
