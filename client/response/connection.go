package response

//ConnectionResponse struct represents response of a connection
type ConnectionResponse struct {
	Message  string           `json:"message"`
	Metadata PostMetaResponse `json:"meta"`
}

//ConnectionMetaResponse struct represents meta field of response
type ConnectionMetaResponse struct {
	Status bool `json:"status"`
}
