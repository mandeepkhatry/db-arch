package model

//TODO Remove Db Namespace
//Document defines the struct for Message recieved from user
type Document struct {
	Collection string                 `json:"collection"`
	Data       map[string]interface{} `json:"data"`
	Indices    []interface{}          `json:"indices"`
}
