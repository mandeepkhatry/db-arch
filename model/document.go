package model

//Document defines the struct for Message recieved from user
type Document struct {
	Database   string        `json:"db"`
	Collection string        `json:"collection"`
	Namespace  string        `json:"namespace"`
	Data       interface{}   `json:"data"`
	Indices    []interface{} `json:"indices"`
}
