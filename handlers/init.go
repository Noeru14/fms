/********/
/* INIT */
/********/

package handler

import (
	engine "../elasticengine"
	logger "../logger"
	auth "../auth"
	
	elastic "gopkg.in/olivere/elastic.v3"
	"github.com/gin-gonic/gin"
)

// [DESC] : creates a client to connect to the ElasticSearch db
// the created client will be used by the handlers 
func InitElastic(mode int) {
	var err error
	
	client, err = elastic.NewClient(
		elastic.SetSniff(false),
		elastic.SetURL("http://localhost:9200"),
	)
	
	if err != nil {
		logger.PrintError(err)
	}
		
	// setup the index
	exists, err := client.IndexExists(engine.GetIndexName()).Do()
	if err != nil {
		logger.PrintError(err)
	}
	
	if !exists {
		// creates the index with a mapping
		_, err = client.CreateIndex(engine.GetIndexName()).BodyString(engine.GetMapping()).Do()
		if err != nil {
			logger.PrintError(err)
		}
	}
	
	// updates the database with the possible new datas
	engine.UpdateDatabase(client)
	
	/* BENCHMARK */
	auth.SetToken("sudo")
	/* BENCHMARK */
}

// [DESC] : index func 
// [CURL] : curl -i http://rogue.kdata.fr:8080/
func Index(c *gin.Context) {
	greetings := greetings1 + greetings2 + greetings3 + greetings4
	
	c.String(http200, greetings)
}