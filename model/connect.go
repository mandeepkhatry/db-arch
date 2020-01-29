package model

//Connection defines the struct for Connection client tries to establish
type Connection struct {
	Database  string `json:"db"`
	Namespace string `json:"namespace"`
}
