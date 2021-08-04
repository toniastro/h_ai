package handlers

import (
	"context"
	"encoding/json"
	"github.com/Kamva/mgm"
	"github.com/gorilla/mux"
	"github.com/toniastro/h_ai/api/models"
	"github.com/toniastro/h_ai/api/utils"
	"github.com/twinj/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"time"
)

type RequestJobCreate struct {
	ObjectID string `json:"object_id"`
}

type Response struct {
	JobID string `json:"job_id"`
}

const InProgress = "QUEUED"

func (h *Handler) CreateJob(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var createReq RequestJobCreate

	err := decoder.Decode(&createReq)
	if err != nil {
		response := utils.Message(http.StatusBadRequest, "Something went wrong with input passed")
		utils.Respond(w, http.StatusBadRequest, response)
		return
	}

	if createReq.ObjectID == "" {
		response := utils.Message(http.StatusBadRequest, "ObjectID is required")
		utils.Respond(w, http.StatusBadRequest, response)
		return
	}

	mutexInit := h.RedisSync.NewMutex(createReq.ObjectID)
	if err = mutexInit.Lock(); err != nil {
		log.Println("Another request tried to access a locked object ID")
		log.Fatal(err)
		return
	}

	defer mutexInit.Unlock()

	findJob := models.NewJob()
	opts := options.FindOne().SetSort(bson.D{{"created_at", -1}})
	err = mgm.Coll(findJob).FindOne(context.TODO(), bson.D{{"object_id", createReq.ObjectID}}, opts).Decode(&findJob)
	if err != nil && err != mongo.ErrNoDocuments {
		response := utils.Message(http.StatusBadRequest, "Something went wrong getting job")
		utils.Respond(w, http.StatusBadRequest, response)
		return
	}

	//A task cannot be scheduled twice within 5 minutes
	if err != mongo.ErrNoDocuments && findJob.CreatedAt.Add(5*time.Minute).After(time.Now().UTC()) {
		response := utils.Message(http.StatusOK, "Job is still within 5 minute range")
		response["data"] = &Response{JobID: findJob.JobID}
		utils.Respond(w, http.StatusOK, response)
		return
	}

	newJobEntity := &models.Jobs{
		Status:    InProgress,
		TimeTaken: utils.GenerateRandomSleep(),
		ObjectID:  createReq.ObjectID,
		JobID:     uuid.NewV4().String(),
	}
	err = mgm.Coll(newJobEntity).Create(newJobEntity)
	if err != nil {
		log.Println(err)
		response := utils.Message(http.StatusBadRequest, "Something went wrong creating job. Kindly try again")
		utils.Respond(w, http.StatusBadRequest, response)
		return
	}

	byteBody, err := json.Marshal(newJobEntity)
	if err != nil {
		log.Println(err)
		response := utils.Message(http.StatusBadRequest, "Something went wrong converting job to byte")
		utils.Respond(w, http.StatusBadRequest, response)
		return
	}

	if err = h.Queue.PublishJob(byteBody); err != nil {
		//TODO Crete routine to delete created job
		log.Println(err)
		response := utils.Message(http.StatusBadRequest, "Something went wrong adding job to queue")
		utils.Respond(w, http.StatusBadRequest, response)
		return

	}

	response := utils.Message(http.StatusOK, "Job has been added to the queue")
	response["data"] = Response{JobID: newJobEntity.JobID}
	utils.Respond(w, http.StatusOK, response)
	return
}

func (h *Handler) GetJobStatus(w http.ResponseWriter, r *http.Request) {
	getParams := mux.Vars(r)
	jobId, exists := getParams["id"]

	if !exists {
		response := utils.Message(http.StatusBadRequest, "JobID not passed to path")
		utils.Respond(w, http.StatusBadRequest, response)
		return
	}

	var findJob models.Jobs
	err := mgm.Coll(&models.Jobs{}).FindOne(context.TODO(), bson.M{"job_id": jobId}).Decode(&findJob)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			response := utils.Message(http.StatusNotFound, "This Job doesn't exist.")
			utils.Respond(w, http.StatusNotFound, response)
			return
		}
		log.Println(err)
		response := utils.Message(http.StatusBadRequest, "Something went wrong while searching for job")
		utils.Respond(w, http.StatusBadRequest, response)
		return
	}

	response := utils.Message(http.StatusOK, "Job Details")
	response["data"] = findJob
	utils.Respond(w, http.StatusOK, response)
	return
}
