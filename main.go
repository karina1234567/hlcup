package main

import (

	"gopkg.in/mgo.v2"
	"goji.io"
	"goji.io/pat"
	"net/http"
	"fmt"
	//"github.com/valyala/fasthttp"

)


func ErrorWithJSON(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	fmt.Fprintf(w, "{message: %q}", message)
}

func ResponseWithJSON(w http.ResponseWriter, json []byte, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	w.Write(json)
}


func main() {
	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	ensureIndex(session)

	//first time if you want to load db from zip (works long!!!)
	//session.DB("testss").DropDatabase()
	//Load(session)

	session.SetMode(mgo.Monotonic, true)
	ensureIndex(session)

	mux := goji.NewMux()
	//users requsts
	mux.HandleFunc(pat.Get("/users"), allUsers(session))
	mux.HandleFunc(pat.Post("/users"), addUser(session))
	mux.HandleFunc(pat.Get("/users/:id"), userById(session))

	//users requsts
	mux.HandleFunc(pat.Get("/locations"), allLocations(session))
	mux.HandleFunc(pat.Get("/locations/:id"), locationById(session))

	// visits requsts
	mux.HandleFunc(pat.Get("/visits"), allVisits(session))
	mux.HandleFunc(pat.Get("/visits/:id"), visitById(session))

	// user + visits requsts
	mux.HandleFunc(pat.Get("/users/:id/:visits/*"), userVisitsById(session))
	mux.HandleFunc(pat.Get("/users/:id/:visits"), userVisitsById(session))

	// average requsts
	mux.HandleFunc(pat.Get("/locations/:id/:avg/*"), locationAvgById(session))
	mux.HandleFunc(pat.Get("/locations/:id/:avg"), locationAvgById(session))

	/*
	mux.HandleFunc(pat.Put("/users/:id"), updateUser(session))
	mux.HandleFunc(pat.Delete("/users/:id"), deleteUser(session))*/
	http.ListenAndServe("localhost:8080", mux)
}

func ensureIndex(s *mgo.Session) {
	session := s.Copy()
	defer session.Close()
	c := session.DB("testss").C("userss")
	index := mgo.Index{
		Key:        []string{"id"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	err := c.EnsureIndex(index)
	if err != nil {
		panic(err)
	}
}










