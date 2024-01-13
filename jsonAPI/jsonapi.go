package jsonapi

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"log"
	"mailingServer/mailDatabase"
	"net/http"
)

func setJsonHeader(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
}

func fromJSON[T any](body io.Reader, target T) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(body)
	json.Unmarshal(buf.Bytes(), &target)
}

func returnJson[T any](w http.ResponseWriter, withData func() (T, error)) {

	setJsonHeader(w)

	data, serverErr := withData()

	if serverErr != nil {
		w.WriteHeader(500)

		serverErrJson, err := json.Marshal(&serverErr)

		if err != nil {
			log.Println(err)
			return
		}

		w.Write(serverErrJson)
		return
	}

	datajson, err := json.Marshal(&data)

	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}

	w.Write(datajson)
}

func returnErr(w http.ResponseWriter, err error, code int) {

	// Define a JSON structure to encapsulate the error message.
	errorMessage := struct {
		Err string
	}{
		Err: err.Error(),
	}

	// Set the HTTP status code.
	w.WriteHeader(code)

	returnJson(w, func() (interface{}, error) {
		return errorMessage, nil
	})
}

func CreateEmail(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

		if req.Method != "POST" {
			log.Println("Request method is not Allowed")
			return
		}

		entry := mailDatabase.EmailEntry{}

		fromJSON(req.Body, &entry)

		err := mailDatabase.CreateEmail(db, entry.Email)
		if err != nil {
			returnErr(w, err, 400)
			return
		}

		returnJson(w, func() (interface{}, error) {
			log.Printf("JSON CreateEmail : %v\n", entry.Email)
			return mailDatabase.GetEmailFromDB(db, entry.Email)
		})
	})
}

func GetEmail(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

		if req.Method != "GET" {
			log.Println("Request method is not Allowed")
			return
		}

		entry := mailDatabase.EmailEntry{}
		fromJSON(req.Body, &entry)

		returnJson(w, func() (interface{}, error) {
			log.Printf("JSON GetEmail : %v\n", entry.Email)
			return mailDatabase.GetEmailFromDB(db, entry.Email)
		})
	})
}

func UpdateEmail(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

		if req.Method != "PUT" {
			log.Println("Request method is not Allowed")
			return
		}

		entry := mailDatabase.EmailEntry{}

		fromJSON(req.Body, &entry)

		err := mailDatabase.UpdateEmail(db, entry)

		if err != nil {
			returnErr(w, err, 400)
			return
		}

		returnJson(w, func() (interface{}, error) {
			log.Printf("JSON UpdateEmail : %v\n", entry.Email)
			return mailDatabase.GetEmailFromDB(db, entry.Email)
		})
	})
}

func DeleteEmail(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

		if req.Method != "POST" {
			log.Println("Request method is not Allowed")
			return
		}

		entry := mailDatabase.EmailEntry{}

		fromJSON(req.Body, &entry)

		err := mailDatabase.DeleteEmail(db, entry.Email)

		if err != nil {
			returnErr(w, err, 400)
			return
		}

		returnJson(w, func() (interface{}, error) {
			log.Printf("JSON DeleteEmail : %v\n", entry.Email)
			return mailDatabase.GetEmailFromDB(db, entry.Email)
		})
	})
}

func GetEmailBatch(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

		if req.Method != "GET" {
			log.Println("Request method is not Allowed")
			return
		}

		queryOptions := mailDatabase.GetEmailBatchQueryParams{}

		fromJSON(req.Body, &queryOptions)

		if queryOptions.Count <= 0 || queryOptions.Page == 0 {
			returnErr(w, errors.New("page and count fields are requried and must be > 0"), 400)
			return
		}

		returnJson(w, func() (interface{}, error) {
			log.Printf("JSON GetEmailBatch : %v\n", queryOptions)
			return mailDatabase.GetEmailBatchFromDB(db, queryOptions)
		})
	})
}

func Serve(db *sql.DB, bind string) {

	http.Handle("/email/create", CreateEmail(db))
	http.Handle("/email/get", GetEmail(db))
	http.Handle("/email/getbatch", GetEmailBatch(db))
	http.Handle("/email/update", UpdateEmail(db))
	http.Handle("/email/delete", DeleteEmail(db))

	err := http.ListenAndServe(bind , nil)

	if err != nil {
		log.Fatalf("JSON server error : %v", err)
	}
	
}
