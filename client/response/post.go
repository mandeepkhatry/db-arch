package response

//PostResponse struct represents response of a document post
type PostResponse struct {
	Message  string           `json:"message"`
	Metadata PostMetaResponse `json:"meta"`
}

//PostMetaResponse struct represents meta field of response
type PostMetaResponse struct {
	Status      bool   `json:"status"`
	Code        string `json:"code"`
	Description string `json:"desc"`
}
