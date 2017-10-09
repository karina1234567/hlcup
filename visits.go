package main

import (
	"gopkg.in/mgo.v2"
	"net/http"
	"gopkg.in/mgo.v2/bson"
	"log"
	"encoding/json"
	"goji.io/pat"
	"fmt"
	"strconv"
	//"os/user"
)

func allVisits(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		session := s.Copy()
		defer session.Close()

		c := session.DB("testss").C("visits")

		var visits []Visit
		err := c.Find(bson.M{}).All(&visits )
		if err != nil {
			ErrorWithJSON(w, "Database error", http.StatusInternalServerError)
			log.Println("Failed get all visits: ", err)
			return
		}

		respBody, err := json.MarshalIndent(visits , "", "  ")
		if err != nil {
			log.Fatal(err)
		}

		ResponseWithJSON(w, respBody, http.StatusOK)
	}
}

func visitById(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		session := s.Copy()
		defer session.Close()

		id := pat.Param(r, "id")
		var i int
		fmt.Sscanf(id, "%5d", &i)

		c :=  session.DB("testss").C("visits")

		var visit Visit
		err := c.Find(bson.M{"id": i}).One(&visit)
		if err != nil {
			ErrorWithJSON(w, "Database error", http.StatusInternalServerError)
			log.Println("Failed find location: ", err)
			return
		}

		if strconv.Itoa(visit.Id) == "" {
			ErrorWithJSON(w, "User not found", http.StatusNotFound)
			return
		}

		respBody, err := json.MarshalIndent(visit, "", "  ")
		if err != nil {
			log.Fatal(err)
		}

		ResponseWithJSON(w, respBody, http.StatusOK)
	}
}

