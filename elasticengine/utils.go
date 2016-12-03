/***************/
/* UTILS FUNCS */
/***************/

package engine

import (
	logger "../logger"
	
	"encoding/csv"
	"strconv"
	"os"
)

// [DESC] : retrieves data from the corresponding csv file
func CsvReader(filenameToLoad string) (map[int]Store, error) {
	var curId, intPostal int
	var floatLat, floatLon float64
	isFirstLine := true // prevents the func to read the title line
	fpStores := make(map[int]Store)
	// fileSplit := strings.Split(filenameToLoad, "/")
	
	logger.Print("[INFO]\t Retrieving data from... " + curFilename)
	
	csvFile, err := os.Open(filenameToLoad)
	if err != nil {
		return fpStores, err
	}
	
	defer csvFile.Close()
	
	reader := csv.NewReader(csvFile)
	reader.Comma = '|'
	reader.FieldsPerRecord = -1
	
	csvData, err := reader.ReadAll()
	if err != nil {
		return fpStores, err
	}
	
	for _, aStore := range csvData {
		if isFirstLine {
			isFirstLine = false
			
		} else {
			curId, err = strconv.Atoi(aStore[0])
			CheckErr(err)
			floatLat, err = strconv.ParseFloat(aStore[1], 64) 
			CheckErr(err)
			floatLon, err = strconv.ParseFloat(aStore[2], 64)
			CheckErr(err)
			intPostal, err = strconv.Atoi(aStore[4])
			CheckErr(err)
			
			fpStores[curId] = Store{Location{floatLat, floatLon}, aStore[3], intPostal, aStore[5]}
		}
	}
	
	return fpStores, nil
}

// [DESC] : converts a variable of type float64 to string
func FloatToString(input_num float64) string {
    return strconv.FormatFloat(input_num, 'f', 6, 64)
}

// [DESC] : standard error checking func
func CheckErr(err error) {
	if err != nil {
		logger.PrintError(err)
	}
}