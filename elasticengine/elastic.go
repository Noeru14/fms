/*******************/
/* INTERFACE FUNCS */
/*******************/

package engine

import (
	elastic "gopkg.in/olivere/elastic.v3"
)

// [DESC] : retrieves the files to insert and inserts the
// corresponding data  to the ElasticSearch database
func UpdateDatabase(client *elastic.Client) {
	insertedStores, insertedFiles := DataRecovery(client)
	DataUpdate(client, insertedStores, insertedFiles)
}

// [DESC] : returns the closest store for each request of
// the input as an array of StoreOutput (see objects.go)
// [PARAM - input] : retrieved Input from the request (see objects.go)
func InitSearch(client *elastic.Client, input Input) (string, error) {
	output, err := HasInput(client, input)
	
	return output, err
}