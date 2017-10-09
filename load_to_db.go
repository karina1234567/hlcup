package main



import (
	"fmt"
	"strings"
	"gopkg.in/mgo.v2"
	"archive/zip"
	"bytes"
	"encoding/json"
	"gopkg.in/mgo.v2/bson"

)


func Load(s *mgo.Session) error {
	read, err := zip.OpenReader("./data.zip")
	if err != nil {
		return fmt.Errorf("Cannot read zip %s", err)
	}
	defer fmt.Printf("load compleate")

	for _, file := range read.File {
		if strings.HasPrefix(file.Name, "data/users") {
			loadUsers(file, s)
		} else if strings.HasPrefix(file.Name, "data/locations") {
			loadLocations(file,s)
		} else if strings.HasPrefix(file.Name, "data/visits") {
			loadVisits(file,s)
		}
			if err != nil {
				fmt.Errorf("Cannot open file %s,  %s", file.Name, err)
				return nil
			}
		}

	VisitOfUser(s)
	VisitOfLocations(s)

	return nil
}

func VisitOfUser(s *mgo.Session){
	session := s.Copy()
	defer session.Close()
	c := session.DB("testss").C("userss")
	c2 := session.DB("testss").C("visits")
	var users[] User
	c.Find(bson.M{}).All(&users)

	for _,u := range users {
		var visits [] Visit
		c2.Find(bson.M{"user": u.Id}).All(&visits)
		u.VisitsOfUser = visits
		c.Remove(bson.M{"id":u.Id })
		c.Insert(u)
	}
	 }

func VisitOfLocations(s *mgo.Session){
	session := s.Copy()
	defer session.Close()
	c := session.DB("testss").C("locations")
	c2 := session.DB("testss").C("visits")
	var locations[] Location
	c.Find(bson.M{}).All(&locations)

	for _,l := range locations {
		var visits [] Visit
		c2.Find(bson.M{"user": l.Id}).All(&visits)
		l.VisitofLocations  = visits
		c.Remove(bson.M{"id":l.Id })
		c.Insert(l)
	}
}

func loadUsers(file *zip.File,s *mgo.Session ) {
	fileReader, err := file.Open()
	buf := new(bytes.Buffer)
	buf.ReadFrom(fileReader)

	//fmt.Printf(buf.String())

	if err != nil {
		fmt.Errorf("Cannot open file %s. Reason %s", file.Name, err)
		return
	}
	defer fileReader.Close()

	session := s.Copy()
	c := session.DB("testss").C("userss")

	data := JsonFileUsers{}
	err = json.Unmarshal(buf.Bytes(), &data)
	if err != nil {
		fmt.Errorf("Cannot unmarshal user file. Reason %s", err)
		return
	}
	for _, user1 := range data.Users {
		c.Insert(user1)
	}
}

func loadLocations(file *zip.File,s *mgo.Session ) {
	fileReader, err := file.Open()
	buf := new(bytes.Buffer)
	buf.ReadFrom(fileReader)

	//fmt.Printf(buf.String())

	if err != nil {
		fmt.Errorf("Cannot open file %s. Reason %s", file.Name, err)
		return
	}
	defer fileReader.Close()

	session := s.Copy()
	c := session.DB("testss").C("locations")

	data := JsonFileLocations{}
	err = json.Unmarshal(buf.Bytes(), &data)
	if err != nil {
		fmt.Errorf("Cannot unmarshal user file. Reason %s", err)
		return
	}


	for _, user1 := range data.Locations {
		c.Insert(user1)
	}
}

func loadVisits(file *zip.File,s *mgo.Session ) {
	fileReader, err := file.Open()
	buf := new(bytes.Buffer)
	buf.ReadFrom(fileReader)
	if err != nil {
		fmt.Errorf("Cannot open file %s. Reason %s", file.Name, err)
		return
	}
	defer fileReader.Close()

	session := s.Copy()
	c := session.DB("testss").C("visits")

	data := JsonFileVisits{}
	err = json.Unmarshal(buf.Bytes(), &data)
	if err != nil {
		fmt.Errorf("Cannot unmarshal user file. Reason %s", err)
		return
	}

	for _, user1 := range data.Visits {
		c.Insert(user1)
	}


}







