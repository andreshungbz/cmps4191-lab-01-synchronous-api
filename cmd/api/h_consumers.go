package main

import (
	"errors"
	"net/http"

	"github.com/andreshungbz/cmps4191-lab-01-synchronous-api/internal/data"
	"github.com/andreshungbz/cmps4191-lab-01-synchronous-api/internal/validator"
)

// createConsumerHandler reads JSON input to create a consumer record, then returns JSON of the created record.
func (app *application) createConsumerHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	consumer := &data.Consumer{
		Name:  input.Name,
		Email: input.Email,
	}

	v := validator.New()
	v.Check(consumer.Name != "", "name", "must be provided")
	data.ValidateContactEmail(v, consumer.Email)
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Consumer.Insert(consumer)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusCreated, envelope{"consumer": consumer}, nil)
}

// showConsumerHandler reads the UUID in the URL path, then returns JSON of the matching consumer record.
func (app *application) showConsumerHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readUUIDParam("id", r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	consumer, err := app.models.Consumer.GetByID(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"consumer": consumer}, nil)
}

// updateConsumerHandler updates the consumer record matching UUID in the URL path, then returns JSON of the updated record.
func (app *application) updateConsumerHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readUUIDParam("id", r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	consumer, err := app.models.Consumer.GetByID(id)
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
		Name   *string `json:"name"`
		Email  *string `json:"email"`
		Status *string `json:"status"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Name != nil {
		consumer.Name = *input.Name
	}
	if input.Email != nil {
		consumer.Email = *input.Email
	}
	if input.Status != nil {
		consumer.Status = data.ConsumerStatus(*input.Status)
	}

	v := validator.New()
	if data.ValidateConsumer(v, consumer); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Consumer.Update(consumer)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"consumer": consumer}, nil)
}

// deleteConsumerHandler deletes the consumer record matching UUID in the URL path.
func (app *application) deleteConsumerHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readUUIDParam("id", r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.Consumer.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"message": "consumer successfully deleted"}, nil)
}
