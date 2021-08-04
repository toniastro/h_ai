package utils

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"time"
)

func Respond(w http.ResponseWriter, status int, data map[string]interface{}) {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func Message(status int, message string) map[string]interface{} {
	return map[string]interface{}{"status": status, "message": message}
}

func GenerateRandomSleep() int {
	rand.Seed(time.Now().UnixNano())
	min := 15
	max := 40
	return rand.Intn(max-min+1) + min
}
