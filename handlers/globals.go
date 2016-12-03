/***********/
/* GLOBALS */
/***********/

package handler

import (
	elastic "gopkg.in/olivere/elastic.v3"
	"sync"
)

type Response struct {
	Code 			int			`json:"code"`
	Content 		string		`json:"content"`
}

type GetNearestResponse struct {
	Code 			int			`json:"code"`
	Mode			string		`json:"mode"`
	ReqNum			int			`json:"reqnum"`
}

const (
	greetings1 						= "Welcome to FindMyFpStore!\n"
	greetings2 						= "1) To send any requests, please connect (/connect) and use the generated token in your requests\n"
	greetings3 						= "2) Send your JSON requests (/getnearest), using the appropriate structure\n"
	greetings4 						= "3) Retrieve your result and/or its status (/result/:id)"
	
	unauthorizedTokenError 			= "Specified token hasn't been authorized"
	authenticationFailedError		= "Wrong username and/or password"
	incorrectJsonError 				= "JSON input isn't correct"
	invalidIdError 					= "Specified id has an unvalid format"
	idNotFoundError 				= "Specified id is unknown"
	
	http200 						= 200 // successful http requests
	http202 						= 202 // the request has been accepted for processing
	http400 						= 400 // bad request
	http401 						= 401 // unauthorized
	http404 						= 404 // not found
)

var (
	client 							*elastic.Client
	curId 							int
	results 						= make(map[int]Response)
	mutex 							= &sync.Mutex{}
	resMutex 						= &sync.RWMutex{}
)