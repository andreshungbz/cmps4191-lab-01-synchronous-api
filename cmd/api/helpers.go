package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
)

// envelope is used to enclose a JSON response.
type envelope map[string]any

// writeJSON attempts to encode data into JSON, applies given HTTP headers,
// and writes to the HTTP response with the given HTTP status code.
func (app *application) writeJSON(w http.ResponseWriter, status int, data envelope, headers http.Header) error {
	// encoding data into JSON
	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}
	js = append(js, '\n')

	// apply HTTP headers
	for key, values := range headers {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	// set Content-Type and write to HTTP response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}

// readJSON attempts to decode JSON from the client, writing to a destination object.
// It checks for errors in the JSON input form and errors in applying the appropriate
// types to the destination.
func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	// set 1MB limit for HTTP request body
	r.Body = http.MaxBytesReader(w, r.Body, 1_048_576)

	// create and configure JSON decoder
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	// decode and check for errors
	err := dec.Decode(dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError
		var maxBytesError *http.MaxBytesError

		switch {
		// Badly-formed JSON
		case errors.As(err, &syntaxError):
			return fmt.Errorf("Body contains badly-formed JSON (at character %d)", syntaxError.Offset)
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("Body contains badly-formed JSON")

		// Incorrect JSON types for destination fields
		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("Body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("Body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)

		// Empty HTTP request body
		case errors.Is(err, io.EOF):
			return errors.New("Body must not be empty")

		// Unknown JSON fields for destination fields
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("Body contains unknown key %s", fieldName)

		// Too-large HTTP request
		case errors.As(err, &maxBytesError):
			return fmt.Errorf("Body must not be larger than %d bytes", maxBytesError.Limit)

		// Programmer error: Passing non-nil pointer
		case errors.As(err, &invalidUnmarshalError):
			panic(err)

		default:
			return err
		}
	}

	// check for extraneous input
	err = dec.Decode(&struct{}{})
	if !errors.Is(err, io.EOF) {
		return errors.New("Body must only contain a single JSON value")
	}

	return nil
}

// readUUIDParam extracts a UUID path parameter from the request context and parses it.
func (app *application) readUUIDParam(paramName string, r *http.Request) (uuid.UUID, error) {
	params := httprouter.ParamsFromContext(r.Context())
	paramStr := params.ByName(paramName)

	id, err := uuid.Parse(paramStr)
	if err != nil {
		return uuid.Nil, errors.New("invalid uuid parameter")
	}

	return id, nil
}

// background launches a function as a goroutine and handles panics for it.
func (app *application) background(fn func()) {
	// increment WaitGroup counter
	app.wg.Go(func() {
		defer func() {
			pv := recover()
			if pv != nil {
				app.logger.Error(fmt.Sprintf("%v", pv))
			}
		}()

		fn()
	})
}
