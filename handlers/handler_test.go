package handler 

import (
	"testing"
	
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"bytes"
	"time"
)

type Coord struct {
	RequestId	int			`json:"request_id"`
	Lat			float64		`json:"lat"`
	Lon			float64		`json:"lon"`
}

type Res struct {
	RequestID	int			`json:"request_id"`
	StoreID		int			`json:"store_id"`
	Address 	string		`json:"address"`
	Postal		int			`json:"postal"`
	City		string		`json:"city"`
}

// data structure sent to FMS
type Input struct {
	Token		string		`json:"token"`
	Mode		string		`json:"mode"`
	Distance 	int			`json:"distance"`
	Coords		[]Coord		`json:"coord"`
}

const (
	authenticationURL 		= "http://rogue.kdata.fr:8080/connect"
	queryURL 				= "http://rogue.kdata.fr:8080/getnearest"
	resultURL				= "http://rogue.kdata.fr:8080/result/"
	
)

var (
	httpClient 				= &http.Client{}
)

// Connect func test with right user infos
// > Check that an OK code (200) is returned 
func TestConnectOK(t *testing.T){
	tokenCode := ConnectTest(t, "admin", "83c0224fdfd3e78088b21008bdd92bfe")
	
	if tokenCode != 200 {
		t.Error("Unexpected returned code, is:", tokenCode, "; expected : 200")
	}
}

// Connect func test with wrong user infos
// > Check that an unauthorized code (401) is returned
func TestConnectKO(t *testing.T) {
	tokenCode := ConnectTest(t, "testKO", "testKO")
	
	if tokenCode != 401 {
		t.Error("Unexpected returned code, is:", tokenCode, "; expected: 401")
	}
}

// GetNearest func test with no distance limit set
// > Check that the returned storeID is the one expected based on the used coordinates
func TestGetNearestNoLimit(t *testing.T) {
	// case 1 : single coord
	var coordOne []Coord
	coordOne = append(coordOne, Coord{1, 48.84858, 2.55261})
	
	reqNum := Query(t, 0, coordOne)
	results := Recovery(t, reqNum)
	
		// single storeID check
	if results[0].StoreID != 6045 {
		t.Error("Unexpected returned storeID, is:", results[0].StoreID, "; expected: 6045")
		
	}
	
	// case 2 : multiple coords
	var coordTwo []Coord
	coordTwo = append(coordTwo, Coord{1, 48.797764, 2.1353993}, Coord{2, 48.6569781, 2.4104558})
	
	reqNum = Query(t, 0, coordTwo)
	results = Recovery(t, reqNum)
	
		// first storeID check
	if results[0].StoreID != 5602 {
		t.Error("Unexpected returned storeID, is:", results[0].StoreID, "; expected: 5602")
		
	} 
	
		// second storeID check
	if results[1].StoreID != 5316 {
		t.Error("Unexpected returned storeID, is:", results[1].StoreID, "; expected: 5316")
		
	}
}

/*
func TestGetNearestLimit(t *testing.T) {
	// one coord
		// check expected store
		// coord that should be returning no value
	
	// multiple coords
		// check expected stores for every coords
		// check stores + no stores
		// check no stores
	
}

func TestResult() {
	// error
	// processing
	// done
	
}
*/

// TestConnectOK, TestConnectKO
func ConnectTest(t *testing.T, username string, password string) int {
	var tokenResponse Response
	var jsonStr = []byte(`{"username": "` + username + `", "password" : "`+ password + `"}`)
	
	req, err := http.NewRequest("POST", authenticationURL, bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal("(Connect) NewRequest:", err)
	}
	
	resp, err := httpClient.Do(req)
	if err != nil {
		t.Fatal("(Connect) Do:", err)
	}
	
	defer resp.Body.Close()
	
	htmlData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal("(Connect) ReadAll:", err)
	}
	
	err = json.Unmarshal(htmlData, &tokenResponse)
	if err != nil {
		t.Fatal("(Connect) Unmarshal:", err)
	}
	
	return tokenResponse.Code
}

// TestGetNearestLimit, TestGetNearestNoLimit, TestResult
func Query(t *testing.T, distance int, coords []Coord) int {
	var fmsResponse GetNearestResponse
	anInput := Input{"sudo", "getstore", distance, coords}
	
	byteInput, err := json.Marshal(anInput)
	if err != nil {
		t.Fatal("(Query) Marshal:", err)
	}
	
	req, err := http.NewRequest("POST", queryURL, bytes.NewBuffer(byteInput))
	if err != nil {
		t.Fatal("(Query) NewRequest:", err)
	}
	
	resp, err := httpClient.Do(req)
	if err != nil {
		t.Fatal("(Query) Do:", err)
	}
	defer resp.Body.Close()
	
	htmlData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal("(Query) ReadAll:", err)
	}
	
	err = json.Unmarshal(htmlData, &fmsResponse)
	if err != nil {
		t.Fatal("(Query) Unmarshal:", err)
	}
	
	return fmsResponse.ReqNum
}

// TestGetNearestLimit, TestGetNearestNoLimit, TestResult
func Recovery(t *testing.T, requestID int) []Res {
	var results []Res
	
	for {
		req, err := http.NewRequest("GET", resultURL + strconv.Itoa(requestID), nil)
		if err != nil {
			t.Fatal("(Recovery) NewRequest:", err)
		}
		
		resp, err := httpClient.Do(req)
		if err != nil {
			t.Fatal("(Recovery) Do:", err)
		}
		defer resp.Body.Close()
		
		htmlData, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatal("(Recovery) ReadAll:", err)
		}
		
		var IntermediateResults Response
		
		err = json.Unmarshal(htmlData, &IntermediateResults)
		if err != nil {
			t.Fatal("(Recovery) Unmarshal:", err)
		}
		
		if IntermediateResults.Code == 400 || IntermediateResults.Code == 404 {
			t.Log("(Recovery) Error:", IntermediateResults.Content)
			break
			
		} else if IntermediateResults.Code == 202 {
			time.Sleep(200 * time.Millisecond)
			
		} else {
			convertedResults := []byte(IntermediateResults.Content)
			json.Unmarshal(convertedResults, &results)
			break
		}
			
		/*
		} else {
			t.Fatal("(Recovery) Unexpected returned value:", IntermediateResults.Code)
			
		}
		*/
	}
	
	return results
}