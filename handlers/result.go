/**********/
/* RESULT */
/**********/

package handler

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
	"time"
)

// [DESC] : returns the state of the researched request
// [CURL] : curl -i http://rogue.kdata.fr:8080/result/<insert id here>
func Result(c *gin.Context) {
	rawId := c.Param("id")
	specifiedId, err := strconv.Atoi(rawId)
	
	if err != nil {
		c.JSON(http404, Response{http404, invalidIdError})
		
	}
	
	resMutex.RLock()
	result, isFound := results[specifiedId]
	resMutex.RUnlock()
	
	if !isFound {
		c.JSON(http404, Response{http404, idNotFoundError})
		
	} else {
		// processing, error or done
		if result.Code == http202 {
			c.JSON(http202, result)
			
		} else if result.Code == http400 {
			c.JSON(http400, result)
			
		} else if strings.Contains(result.Content, "\"request_id\"") {
			c.JSON(http200, result)
			
			// deletes the retrieved result 10 seconds later
			go func() {
				time.Sleep(10 * time.Second)
				
				resMutex.Lock()
				delete(results, specifiedId)
				resMutex.Unlock()
				
			} ()
		}
	}
}