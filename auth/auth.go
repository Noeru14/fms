package auth

import (
	logger "../logger"
	
	"encoding/hex"
	"crypto/md5"
	"math/rand"
	"strings"
	"time"
)

type Login struct {
	Username 	string		`json:"username"`
	Password 	string		`json:"password"`
}

const (
	username = "admin"
	password = "jZ8B69Hr7D"
	
	letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
    letterIdxBits = 6
    letterIdxMask = 1 << letterIdxBits - 1
    letterIdxMax  = 63 / letterIdxBits
)

var authorizedTokens = make(map[string]bool)
var src = rand.NewSource(time.Now().UnixNano())

// [DESC] : if the username and the password are correct
// generates and returns a token
func Connection(login Login) (bool, string) {
	var curToken string
	
	logger.Print("An authentication attempt is being processed...")
	
	md5Password := md5.Sum([]byte(password))
	strMd5Password := hex.EncodeToString(md5Password[:])
	
	if (strings.Compare(login.Username, username) == 0) && (strings.Compare(login.Password, strMd5Password) == 0) {
		curToken = tokenGeneration()
		
		return true, curToken
	}
	
	return false, curToken
}

// [DESC] : adds and returns the generated token
func tokenGeneration() string {
	curToken := RandStringBytesMaskImprSrc(20)
	authorizedTokens[curToken] = true
	
	return curToken
}

// [DESC] : creates a random string of the specified length
func RandStringBytesMaskImprSrc(n int) string {
    b := make([]byte, n)
    for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
        if remain == 0 {
            cache, remain = src.Int63(), letterIdxMax
        }
        if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
            b[i] = letterBytes[idx]
            i--
        }
        cache >>= letterIdxBits
        remain--
    }

    return string(b)
}

// [DESC] : checks if the token is authorized
func CheckToken(token string) bool {
	isAuthorized := authorizedTokens[token]
	
	if isAuthorized {
		return true
	}
	
	return false
}

/* GETTERS */
func GetUsername() string {
	return username
}

func GetPassword() string {
	return password
}

/* BENCHMARK */
func SetToken(desiredToken string) {
	authorizedTokens[desiredToken] = true
}