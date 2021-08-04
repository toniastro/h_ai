package main

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/toniastro/h_ai/api/handlers"
	"github.com/toniastro/h_ai/api/models"
	"log"
	"net/http"
	"strings"
	"testing"
)

var object = "1234"
var JobId string

type StatusResponse struct {
	Data models.Jobs
	Status int
	Message string
}

type CreateResponse struct {
	Data handlers.Response
	Status int
	Message string
}

func TestCreateJob(t *testing.T) {
	request := `{"object_id":"` + object + `"}`
	req, err := http.NewRequest(http.MethodPost, "http://localhost:8100/api/job", strings.NewReader(request))
	if err != nil {
		log.Fatal(err)
	}

	client := http.Client{}
	response, err := client.Do(req)
	if err != nil {
		log.Println("Check that application is running")
		log.Fatal(err)
	}

	var res CreateResponse
	err = json.NewDecoder(response.Body).Decode(&res)
	if err != nil {
		log.Println("Something went wrong decoding")
		log.Fatal(err)
	}

	JobId = res.Data.JobID
	assert.NotNil(t, JobId)
	response.Body.Close()
}

func TestCheckJobStatus(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "http://localhost:8100/api/job/"+JobId, nil)
	if err != nil {
		log.Fatal(err)
	}

	client := http.Client{}
	response, err := client.Do(req)
	if err != nil {
		log.Println("Check that application is running")
		log.Fatal(err)
	}

	var res StatusResponse
	err = json.NewDecoder(response.Body).Decode(&res)
	if err != nil {
		log.Println("Something went wrong decoding")
		log.Fatal(err)
	}

	if assert.NotNil(t, res) {

		assert.Equal(t, "Job Details", res.Message)
		assert.Equal(t, 200, res.Status)
		assert.NotNil(t, res.Data)

	}
	response.Body.Close()
}
