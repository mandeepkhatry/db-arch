package model

//Query struct
type Query struct {
	Database   string                 `json:"db"`
	Collection string                 `json:"collection"`
	Namespace  string                 `json:"namespace"`
	Querydata  map[string]interface{} `json:"querydata"`
}
