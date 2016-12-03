/****************/
/* ENGINE FUNCS */
/****************/

package engine

import (
	logger "../logger"
	
	elastic "gopkg.in/olivere/elastic.v3"
	"github.com/cheggaaa/pb"
	"path/filepath"
	"encoding/json"
	"strconv"
	"strings"
	"errors"
	"time"
	"os"
)

// [DESC] : updates the database with the recovered data from the files to insert
// [VAR - insertedStores] : contains the id of the stores that have already been inserted
// [VAR - insertedFiles] : contains the name of the files that have already been inserted
func DataUpdate(client *elastic.Client, insertedStores map[int]bool, insertedFiles map[string]bool) {
	var fpStores map[int]Store
	var fileSplit []string
	searchDir := "/home/findMyFpStore/filesToInsert"
	
	logger.Print("Starting DataUpdate phase...")

	fileList := []string{}
	err := filepath.Walk(searchDir, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			logger.PrintError(err)
		}
		
		if !f.IsDir() {
			fileList = append(fileList, path)
		}
		
		return nil
	})

	for _, filename := range fileList {
		fileSplit = strings.Split(filename, "/")
		// fileFound := insertedFiles[fileSplit[2]]
		
		curFilename = fileSplit[len(fileSplit) - 1]
		fileFound := insertedFiles[curFilename]
		logger.Print("Current file is:", curFilename)
		
		if !fileFound {
			logger.Print("Inserting datas from... " + curFilename)
			
			file := File{curFilename, time.Now().Format("2006-01-02 15:04:05")}
			
			// file inserts to Elastic
			_, err = client.Index().
				Index(indexName).
				Type(fileType).
				BodyJson(file).
				Refresh(true).
				Do()
			CheckErr(err)
		
			fpStores, err = CsvReader(filename)
			if err != nil {
				logger.PrintError(err)
				continue
			}
			
			progressBar := pb.StartNew(len(fpStores))
			
			// insertions from retrieved data
			for storeId, store := range fpStores {
				isInserted := insertedStores[storeId]
				
				if !isInserted {
					// store inserts to Elastic
					_, err = client.Index().
						Index(indexName).
						Type(storeType).
						Id(strconv.Itoa(storeId)).
						BodyJson(store).
						Refresh(true).
						Do()
					CheckErr(err)
					
				}
				
				progressBar.Increment()
			}
		}
	}
	
	logger.Print("DataUpdate has ended")
}

// [DESC] : recovers data from the files to insert
func DataRecovery (client *elastic.Client) (map[int]bool, map[string]bool) {
	insertedStores := make(map[int]bool)
	insertedFiles := make(map[string]bool)
	
	logger.Print("Starting DataRecovery phase...")
	
	// prevents shards failed exception
	time.Sleep(500 * time.Millisecond)
	
	/* STORES */
	matchAllQuery := elastic.NewMatchAllQuery()
	searchResult, err := client.Search().
		Index(indexName).
		Type(storeType).
		Query(matchAllQuery). 
		Do()
	CheckErr(err)
	
	// iterate through results with full control over each step
	if searchResult.Hits.TotalHits > 0 {

		// iterate through results (hit is of type SearchHit)
		for _, hit := range searchResult.Hits.Hits {
			intHidId, err := strconv.Atoi(hit.Id)
			CheckErr(err)
			
			insertedStores[intHidId] = true
		}
	}
	
	/* FILES */
	matchAllQuery = elastic.NewMatchAllQuery()
	searchResult, err = client.Search().
		Index(indexName).
		Type(fileType).
		Query(matchAllQuery). 
		Do()
	CheckErr(err)
	
	// iterate through results with full control over each step
	if searchResult.Hits.TotalHits > 0 {

		// tterate through results (hit is of type SearchHit)
		for _, hit := range searchResult.Hits.Hits {
			// deserialize hit.Source into a file (could also be just a map[string]interface{})
			var f File
			err := json.Unmarshal(*hit.Source, &f)
			CheckErr(err)
			
			insertedFiles[f.Filename] = true
		}
	}
	
	logger.Print("DataRecovery has ended")
	
	return insertedStores, insertedFiles
}

// [DESC] : calls GetStore to retrieve and return the
// closest store for each request of the input
func HasInput(client *elastic.Client, input Input) (string, error) {
	var output string
	var err error
	
	switch input.Mode {
		case "getstore":
			output, err = GetStore(client, input.Coords, input.Distance)
			
			return output, err
			
		// case "insertOtherModeHere" : 
		
		// mode isn't recognized
		default:
			logger.PrintError("Specified mode isn't recognized : " + input.Mode)
			return output, err
			
	}
	
	return output, errors.New("Unexpected behaviour")
}

// [DESC] : retrieves the closest store for each coord
// and returns the marshalized StoreOutput array
func GetStore(client *elastic.Client, coords []Coord, distance int) (string, error) {
	var output []StoreOutput
	var intStoreId int
	
	logger.Print("Query of " + strconv.Itoa(len(coords)) + " coords has started")
	
	for _, aCoord := range coords {
		var aStoreOutput StoreOutput
		
		// sorter : closest location first
		geoDistanceSorter := elastic.NewGeoDistanceSort("location").
			Point(aCoord.Lat, aCoord.Lon).
			Unit("m").
			GeoDistance("plane").
			Asc()

		// query execution
		matchAllQuery := elastic.NewMatchAllQuery()
		searchResult, err := client.Search().
			Index(indexName).
			Type(storeType).
			Query(matchAllQuery).
			Size(1).
			SortBy(geoDistanceSorter).
			Do()
			
		if err != nil {
			return "err", err
		}
		
		// aStoreOutput.RequestId = aCoord.RequestId
		if searchResult.Hits.TotalHits > 0 {
			for _, hit := range searchResult.Hits.Hits {
				var s Store
				var tooFar bool
				
				err := json.Unmarshal(*hit.Source, &s)
				if err != nil {
					return "err", err
				}

				intStoreId, err = strconv.Atoi(hit.Id)
				if err != nil {
					return "err", err
				}
				
				// distance checking if radius restriction
				// radius restriction if distance bigger than zero
				if distance > 0 {
					
					// if the distance separating the current coord from the 
					// store is bigger than the radius restriction, nil return
					distFromCoord, _ := hit.Sort[0].(float64)
					
					if distFromCoord > float64(distance) {
						tooFar = true
						
					}
				}
				
				if !tooFar {
					aStoreOutput.RequestId = aCoord.RequestId
					aStoreOutput.StoreId = intStoreId
					aStoreOutput.Address = s.Address
					aStoreOutput.Postal = s.Postal
					aStoreOutput.City = s.City
				}
			}
		} 
		
		output = append(output, aStoreOutput) 
	}
	
	marshalizedOutput, err := json.Marshal(output)
	if err != nil {
		return "err", err
	}
	
	stringifiedOutput := string(marshalizedOutput[:])
	
	return stringifiedOutput, nil
}