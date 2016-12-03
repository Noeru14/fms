/********/
/* MAIN */
/********/

package main

import (
	handler "./handlers"
	logger "./logger"
	
	"github.com/gin-gonic/gin"
)

func main() {
	logger.Print("##################################################")
	logger.Print("####                                          ####")
	logger.Print("####    FindMyFpStore v1.4.2 - William Ung    ####")
	logger.Print("####                                          ####")
	logger.Print("##################################################")
	
	handler.InitElastic(0)
	
	// [RELEASE MODE] : debug mode minus recovery and logger
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	
	// [DEBUG MODE] : uses recovery and logger 
	// logger : instances a Logger middleware that will write the logs 
	// to gin.DefaultWriter By default gin.DefaultWriter = os.Stdout    
	// recovery :  returns a middleware that recovers from any panics and writes a 500 if there was one
	// router := gin.Default()
	
	// routes
	router.GET("/", handler.Index)
	router.POST("/connect", handler.Connect)
	router.POST("/getnearest", handler.GetNearest)
	router.GET("/result/:id", handler.Result)
	
	// run
	router.Run(":8080")
}
