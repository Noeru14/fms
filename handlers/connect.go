/***********/
/* CONNECT */
/***********/

package handler

import (
	auth "../auth"
	
	"github.com/gin-gonic/gin"
)

// [DESC] : authentication process, if
// the login and password are right, returns a token
// [CURL] : curl -X POST http://rogue.kdata.fr:8080/connect -d '{"username" : "admin", "password" : "83c0224fdfd3e78088b21008bdd92bfe"}'
func Connect(c *gin.Context) {
	var login auth.Login
	var isRight bool
	var token string
	
	c.BindJSON(&login)
	isRight, token = auth.Connection(login)
	
	if isRight {
		c.JSON(http200, Response{http200, token})
		
	} else {
		c.JSON(http401, Response{http401, authenticationFailedError})
		
	}
}