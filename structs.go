package main


type User struct {
	Id        int 	 	`json:"id"`
	Email     string 	`json:"email"`
	FirstName string 	`json:"first_name"`
	LastName  string 	`json:"last_name"`
	Gender    string 	`json:"gender"`
	Birthdate int    	`json:"birth_date"`
	VisitsOfUser 	  []Visit 	`json:"-"`

}

type JsonFileUsers struct {
	Users []*User
}

type UserVisits struct {
	VisitedAt int    `json:"visited_at"`
	Mark      int    `json:"mark"`
	Place     string `json:"place"`
}

type UserVisitsArray struct {
	Visits [] UserVisits
}

type Visit struct {
	Id        int `json:"id"`
	Location  int `json:"location"`
	User      int `json:"user"`
	VisitedAt int `json:"visited_at"`
	Mark      int `json:"mark"`
}

type JsonFileVisits struct {
	Visits []*Visit
}

type Location struct {
	Id       int    `json:"id"`
	Place    string `json:"place"`
	Country  string `json:"country"`
	City     string `json:"city"`
	Distance int    `json:"distance"`
	VisitofLocations []Visit 	`json:"-"`
}

type JsonFileLocations struct {
	Locations []*Location
}










