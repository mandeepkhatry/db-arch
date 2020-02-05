package response

//QueryResponse struct represents response of a query
type QueryResponse struct{
	Result []map[string]interface{} 				`json:"data"`
	Metadata  QueryMetaResponse  					`json:"meta"`
}

//QueryMetaResponse struct represents meta field of response
type QueryMetaResponse struct{
	Status 			bool 								`json:"status"`
	Code 			string								`json:"code"`
	Description 	string								`json:"desc"`
}
