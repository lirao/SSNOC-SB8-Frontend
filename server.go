package main

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessions"
	"fmt"
	"net/http"
	"bytes"
)

type User struct {
	Name   string
	Status int
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
	urlBack := fmt.Sprintf("http://localhost:9000/user/%s", name)
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
	jsonReqBack := []byte(fmt.Sprintf("{\"password\":\"%s\"}", password))
	urlBack := fmt.Sprintf("http://localhost:9000/user/%s/authenticate", userName)

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
	status := PostUserAuthenticate(userName, r.FormValue("password"))
	switch status {
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
		jsonReqBack := []byte(fmt.Sprintf("{\"userName\":\"%s\",\"password\":\"%s\",\"createdAt\":\"%s\"}",
			userName, password, "2014-10-05 09:12"))
		urlBack := "http://localhost:9000/user/signup"
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
