/***********/
/* GLOBALS */
/***********/

package engine

/* ELASTIC DB MGT */
type File struct {
	Filename 	string 		`json:"filename"`
	InsertDate 	string		`json:"insertdate"`
}

type Location struct {
	Lat			float64 	`json:"lat"`
	Lon			float64		`json:"lon"`
}

type Store struct {
	Loc			Location	`json:"location"`
	Address		string		`json:"address"`
	Postal		int			`json:"postal"`
	City		string		`json:"city"`
}

/* INPUT FILE MGT */
type Coord struct {
	RequestId	int			`json:"request_id"`
	Lat			float64		`json:"lat" binding:"required"`
	Lon			float64		`json:"lon" binding:"required"`
}

type Input struct {
	Token		string		`json:"token" binding:"required"`
	Mode		string		`json:"mode" binding:"required"`
	Distance 	int			`json:"distance" binding:"required"`
	Coords		[]Coord		`json:"coord" binding:"required"`
}

/* OUTPUT FILE MGT */
type StoreOutput struct {
	RequestId	int			`json:"request_id"`
	StoreId		int			`json:"store_id"`
	Address 	string		`json:"address"`
	Postal		int			`json:"postal"`
	City		string		`json:"city"`
}

var (
	curFilename 		string
)

const (
	indexName 			= "findmyfpstore"
	storeType  			= "store"
	fileType 			= "file"
	mapping 			= `{
		"settings":{
			"number_of_shards":1,
			"number_of_replicas":0
		},
		"mappings":{
			"store":{
				"properties":{
					"location":{
						"type":"geo_point"
					},
					"address":{
						"type":"string"
					},
					"postal":{
						"type":"long"
					},
					"city":{
						"type":"string"
					}
				}
			}
		}
	}`
)

/* CONSTS GETTERS */
func GetIndexName() string {
	return indexName
}

func GetMapping() string {
	return mapping
}