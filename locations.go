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
	"strings"
)

func allLocations(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		session := s.Copy()
		defer session.Close()

		c := session.DB("testss").C("locations")

		var locations []Location
		err := c.Find(bson.M{}).All(&locations )
		if err != nil {
			ErrorWithJSON(w, "Database error", http.StatusInternalServerError)
			log.Println("Failed get all locations: ", err)
			return
		}

		respBody, err := json.MarshalIndent(locations , "", "  ")
		if err != nil {
			log.Fatal(err)
		}

		ResponseWithJSON(w, respBody, http.StatusOK)
	}
}

func locationById(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		session := s.Copy()
		defer session.Close()

		id := pat.Param(r, "id")
		var i int
		fmt.Sscanf(id, "%5d", &i)

		c :=  session.DB("testss").C("locations")

		var location Location
		err := c.Find(bson.M{"id": i}).One(&location)
		if err != nil {
			ErrorWithJSON(w, "Database error", http.StatusInternalServerError)
			log.Println("Failed find location: ", err)
			return
		}

		if strconv.Itoa(location.Id) == "" {
			ErrorWithJSON(w, "User not found", http.StatusNotFound)
			return
		}

		respBody, err := json.MarshalIndent(location, "", "  ")
		if err != nil {
			log.Fatal(err)
		}

		ResponseWithJSON(w, respBody, http.StatusOK)
	}
}

func locationAvgById(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		session := s.Copy()
		defer session.Close()

		id := pat.Param(r, "id")
		var i int
		fmt.Sscanf(id, "%5d", &i)

		req := r.RequestURI

		var location Location
		var user User
		result := 0.0
		count := 0

		c := session.DB("testss").C("locations")
		c2 := session.DB("testss").C("userss")
		c.Find(bson.M{"id": i}).One(&location)

		for _,lv := range location.VisitofLocations {
			if strings.Contains(req, "fromDate") {
				fromDateStr := strings.Split(req, "=")
				date, _ := strconv.Atoi(fromDateStr[1])
				if lv.VisitedAt > date {
					result += float64(lv.Mark)
					count += 1
				}
			} else if strings.Contains(req, "toDate") {
				fromDateStr := strings.Split(req, "=")
				date, _ := strconv.Atoi(fromDateStr[1])
				if lv.VisitedAt < date {
					result += float64(lv.Mark)
					count += 1
				}
			} else if strings.Contains(req, "fromAge") {
				fromDateStr := strings.Split(req, "=")
				c2.Find(bson.M{"id": lv.User}).One(&user)
				dist, _ := strconv.Atoi(fromDateStr[1])
				if user.Birthdate > dist {
					result += float64(lv.Mark)
					count += 1
				}

			} else if strings.Contains(req, "toAge") {
				fromDateStr := strings.Split(req, "=")
				c2.Find(bson.M{"id": lv.User}).One(&user)
				dist, _ := strconv.Atoi(fromDateStr[1])
				if user.Birthdate < dist {
					result += float64(lv.Mark)
					count += 1
				}
			} else if strings.Contains(req, "gender") {
				gender := strings.Split(req, "=")
				c2.Find(bson.M{"id": lv.User}).One(&user)

				if user.Gender == gender[1] {
					result += float64(lv.Mark)
					count += 1
				}
				}else {
				result += float64(lv.Mark)
				count += 1
			}
		}
		if count != 0{
			result = result/float64(count)
		} else{
			result = 0
		}

		respBody, err := json.MarshalIndent(result , "", "  ")
		if err != nil {
			log.Fatal(err)
		}

		ResponseWithJSON(w, respBody, http.StatusOK)



		/*c :=  session.DB("testss").C("userss")
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
	*/
		}
}