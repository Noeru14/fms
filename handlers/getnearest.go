/**************/
/* GETNEAREST */
/**************/

package handler

import (
	engine "../elasticengine"
	logger "../logger"
	auth "../auth"
	
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)

// [DESC] : get the nearest store for each request of the input
// within a specified range if distance is greather than zero
// [CURL] : curl -X POST http://rogue.kdata.fr:8080/getnearest -d '{"token" : "<insert token here>", "mode" : "getstore", "distance" : 2000, "coord" : [{"request_id" : 1, "lat" : 48.9014356, "lon" : 2.2665344}, {"request_id" : 2, "lat" : 46.811319, "lon" : 1.6980988}]}'
// [CURL (file)] : curl -X POST http://rogue.kdata.fr:8080/getnearest -d '@./<filename>'
func GetNearest(c *gin.Context) {
	/* BENCHMARK */
	start := time.Now()
	/* BENCHMARK */
	
	var output, processingMessage string
	var input engine.Input
	var err error
	var myId int
	var resultMessage Response
	
	err = c.BindJSON(&input)
	if err != nil {
		logger.PrintError("BindJSON :", err)
		
		c.JSON(http400, Response{http400, incorrectJsonError})
		
		return
	}
	
	if auth.CheckToken(input.Token) {
		var mode string
		
		// mutex
		mutex.Lock()
		curId++
		myId = curId
		mutex.Unlock()
		
		// data mgt
		if input.Distance > 0 {
			mode = "LimitStore"
			
		} else if input.Distance < 0 {
			mode = "GetStore"
			
		}
		
		processingMessage = "[MODE : " + mode + "] Request " + strconv.Itoa(myId) + " is being processed!"
		
		resultMessage = Response{http200, processingMessage}
		logger.Print(processingMessage)
		
		// returning message
		returnMessage := GetNearestResponse{http200, mode, myId}
		c.JSON(http200, returnMessage)
		
		go func() {
			// status code while InitSearch is processing
			resultMessage.Code = http202
			
			resMutex.Lock()
			results[myId] = resultMessage
			resMutex.Unlock()
			
			// starts the research
			output, err = engine.InitSearch(client, input)
			
			if err != nil {
				resultMessage.Code = http400
				resultMessage.Content = "[ERROR] " + err.Error()
				
				resMutex.Lock()
				results[myId] = resultMessage
				resMutex.Unlock()
				
			} else {
				/* BENCHMARK */
				elapsed := time.Since(start)
				logger.Print("[BENCHMARK] time to complete :", elapsed)
				/* BENCHMARK */
				
				// resultMessage.Code = http200
				resultMessage.Code = int(elapsed)
				resultMessage.Content = output
				
				resMutex.Lock()
				results[myId] = resultMessage
				resMutex.Unlock()
				
			}
		} ()
		
	} else {
		c.JSON(http401, GetNearestResponse{http401, unauthorizedTokenError, 0})
		
	}
}