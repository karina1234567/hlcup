package main

import (
	"gopkg.in/mgo.v2"
	"net/http"
	"gopkg.in/mgo.v2/bson"
	"log"
	"encoding/json"
	"strconv"
	"goji.io/pat"
	"fmt"
	//"github.com/valyala/fasthttp"
	//"strings"
	"strings"
	//"math"
)

func allUsers(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		session := s.Copy()
		defer session.Close()

		c := session.DB("testss").C("userss")

		var users []User
		err := c.Find(bson.M{}).All(&users )
		if err != nil {
			ErrorWithJSON(w, "Database error", http.StatusInternalServerError)
			log.Println("Failed get all users: ", err)
			return
		}

		respBody, err := json.MarshalIndent(users , "", "  ")
		if err != nil {
			log.Fatal(err)
		}

		ResponseWithJSON(w, respBody, http.StatusOK)
	}
}

func addUser(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		session := s.Copy()
		defer session.Close()

		var user User
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&user)
		if err != nil {
			ErrorWithJSON(w, "Incorrect body", http.StatusBadRequest)
			return
		}

		c := session.DB("testss").C("userss")

		err = c.Insert(user)
		if err != nil {
			if mgo.IsDup(err) {
				ErrorWithJSON(w, "User with this id already exists", http.StatusBadRequest)
				return
			}

			ErrorWithJSON(w, "Database error", http.StatusInternalServerError)
			log.Println("Failed insert book: ", err)
			return
		}

		userId :=  strconv.Itoa(user.Id)



		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Location", r.URL.Path+"/"+userId)
		w.WriteHeader(http.StatusCreated)
	}
}



func userById(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		session := s.Copy()
		defer session.Close()

		id := pat.Param(r, "id")
		var i int
		fmt.Sscanf(id, "%5d", &i)

		c :=  session.DB("testss").C("userss")

		var user User
		err := c.Find(bson.M{"id": i}).One(&user)
		if err != nil {
			ErrorWithJSON(w, "Database error", http.StatusInternalServerError)
			log.Println("Failed find user: ", err)
			return
		}

		if strconv.Itoa(user.Id) == "" {
			ErrorWithJSON(w, "User not found", http.StatusNotFound)
			return
		}

		respBody, err := json.MarshalIndent(user, "", "  ")
		if err != nil {
			log.Fatal(err)
		}

		ResponseWithJSON(w, respBody, http.StatusOK)
	}
}

func userVisitsById(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		session := s.Copy()
		defer session.Close()

		id := pat.Param(r, "id")
		var i int
		fmt.Sscanf(id, "%5d", &i)

		req := r.RequestURI

		c :=  session.DB("testss").C("userss")
		c2 := session.DB("testss").C("locations")

		var user  User
		var location Location
		result := make([]UserVisits, 0)

		err := c.Find(bson.M{"id": i}).One(&user)

		for _,uv := range user.VisitsOfUser{

			c2.Find(bson.M{"id": uv.Location}).One(&location)

			if strings.Contains(req, "fromDate") {
				fromDateStr := strings.Split(req, "=")
				date, _ := strconv.Atoi(fromDateStr[1])
				if uv.VisitedAt > date {
					result = append(result, UserVisits{
						VisitedAt: uv.VisitedAt,
						Mark:      uv.Mark,
						Place:     location.Place,
					})
				}
			} else if strings.Contains(req, "toDate") {
				fromDateStr := strings.Split(req, "=")
				date, _ := strconv.Atoi(fromDateStr[1])
				if uv.VisitedAt < date {
					result = append(result, UserVisits{
						VisitedAt: uv.VisitedAt,
						Mark:      uv.Mark,
						Place:     location.Place,
					})
				}
			}else if strings.Contains(req, "toDistance") {
				fromDateStr := strings.Split(req, "=")
				dist, _ := strconv.Atoi(fromDateStr[1])
				if location.Distance < dist{
					result = append(result, UserVisits{
						VisitedAt: uv.VisitedAt,
						Mark:      uv.Mark,
						Place:     location.Place,
					})
				}
			} else {

				result = append(result, UserVisits{
					VisitedAt: uv.VisitedAt,
					Mark:      uv.Mark,
					Place:     location.Place,
				})}

		}

		visitsArray := UserVisitsArray{
			Visits: result,
		}


		respBody, err := json.MarshalIndent(visitsArray , "", "  ")
		if err != nil {
			log.Fatal(err)
		}

		ResponseWithJSON(w, respBody, http.StatusOK)
	}
}

