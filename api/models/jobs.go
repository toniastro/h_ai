package models

import "github.com/Kamva/mgm"

type Jobs struct {
	mgm.DefaultModel `bson:",inline"`
	Status           string `json:"status" bson:"status"`
	TimeTaken        int    `bson:"time_taken" json:"time_taken"`
	ObjectID         string `bson:"object_id" json:"object_id"`
	JobID            string `bson:"job_id" json:"job_id"`
}

func NewJob() *Jobs {
	return &Jobs{}
}
