package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/andreshungbz/cmps4191-lab-01-synchronous-api/internal/data"
	"github.com/andreshungbz/cmps4191-lab-01-synchronous-api/internal/validator"
	"github.com/google/uuid"
)

// createJobHandler reads JSON input to create a job record, then returns JSON of the created record.
func (app *application) createJobHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		ConsumerID uuid.UUID       `json:"consumer_id"`
		JobType    string          `json:"job_type"`
		Payload    json.RawMessage `json:"payload"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	job := &data.Job{
		ConsumerID: input.ConsumerID,
		JobType:    input.JobType,
		Payload:    input.Payload,
	}

	v := validator.New()
	if data.ValidateJob(v, job); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Job.Insert(job)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusCreated, envelope{"job": job}, nil)
}

// showJobHandler reads the UUID in the URL path, then returns JSON of the matching job record.
func (app *application) showJobHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readUUIDParam("public_id", r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	job, err := app.models.Job.GetByPublicID(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"job": job}, nil)
}

// updateJobHandler updates the job record matching UUID in the URL path, then returns JSON of the updated record.
func (app *application) updateJobHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readUUIDParam("public_id", r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	job, err := app.models.Job.GetByPublicID(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	var input struct {
		Status       *string          `json:"status"`
		Result       *json.RawMessage `json:"result"`
		ErrorMessage *string          `json:"error_message"`
		StartedAt    *time.Time       `json:"started_at"`
		CompletedAt  *time.Time       `json:"completed_at"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Status != nil {
		job.Status = data.JobStatus(*input.Status)
	}
	if input.Result != nil {
		job.Result = input.Result
	}
	if input.ErrorMessage != nil {
		job.ErrorMessage = input.ErrorMessage
	}
	if input.StartedAt != nil {
		job.StartedAt = input.StartedAt
	}
	if input.CompletedAt != nil {
		job.CompletedAt = input.CompletedAt
	}

	v := validator.New()
	if data.ValidateJob(v, job); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Job.Update(job)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"job": job}, nil)
}

// deleteJobHandler deletes the job record matching UUID in the URL path.
func (app *application) deleteJobHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readUUIDParam("id", r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.Job.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"message": "job successfully deleted"}, nil)
}
