package mailDatabase

import (
	"database/sql"
	"log"
	"time"

	"github.com/mattn/go-sqlite3"
)

type EmailEntry struct {
	Id          int64
	Email       string
	ConfirmedAt *time.Time
	OptOut      bool
}

// Will try to create a DB table
func TryCreate(db *sql.DB) {

	_, err := db.Exec(`
		CREATE TABLE emails (
			id 			INTEGER PRIMARY KEY,
			email		TEXT UNIQUE,
			confirmed_at INTEGER,
			opt_out 	 INTEGER
		)
	`)

	if err != nil {
		sqlError, ok := err.(sqlite3.Error)
		if ok {
			// if code = 1 it means table already exists
			if sqlError.Code == 1 {
				log.Println(sqlError)
			} else {
				log.Fatal(err)
			}
		}
	}
}

func emailEntryFromRow(row *sql.Rows) (*EmailEntry, error) {
	var id int64
	var email string
	var confirmedAt int64
	var optOut bool

	err := row.Scan(&id, &email, &confirmedAt, &optOut)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	t := time.Unix(confirmedAt, 0)

	return &EmailEntry{id, email, &t, optOut}, nil
}

// Setting the Email in the DB
func CreateEmail(db *sql.DB, email string) error {
	_, err := db.Exec(`
		INSERT INTO
		emails(email, confirmed_at, opt_out)
		VALUES(?, 0, false)`, email)

	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

// Getting the Email from the DB
func GetEmailFromDB(db *sql.DB, email string) (*EmailEntry, error) {

	rows, err := db.Query(`
	SELECT id, email, confirmed_at , opt_out
	FROM emails
	WHERE email = ? `, email)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		return emailEntryFromRow(rows)
	}

	return nil, nil
}

// To update the email entry in Database
func UpdateEmail(db *sql.DB, entry EmailEntry) error {

	var t int64
	if entry.ConfirmedAt != nil {
	t = entry.ConfirmedAt.Unix()
	}else{
		t = 0
	}

	_, err := db.Exec(`
		INSERT INTO 
		emails(email, confirmed_at, opt_out)
		VALUES(?, ? ,?)
		ON CONFLICT(email) DO UPDATE SET
			confirmed_at = ?,
			opt_out = ?`, entry.Email, t, entry.OptOut, t, entry.OptOut)

	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

// To delete the email entry from Database
func DeleteEmail(db *sql.DB, email string) error {

	// Here instead of deleting we are setting it as opt_out = true
	_, err := db.Exec(`
	UPDATE emails SET opt_out=true WHERE email=?`, email)

	if err != nil {
		log.Println("Error updating row:", err)
		log.Println(err)
		return err
	}

	return nil
}

type GetEmailBatchQueryParams struct {
	Page  int
	Count int
}

func GetEmailBatchFromDB(db *sql.DB, params GetEmailBatchQueryParams) ([]EmailEntry, error) {

	var emptyList []EmailEntry

	rows, err := db.Query(`
	SELECT id , email, confirmed_at , opt_out
	FROM emails 
	WHERE opt_out = false
	ORDER BY id ASC
	LIMIT ? OFFSET ?`, params.Count, (params.Page-1)*params.Count)

	if err != nil {

		log.Println(err)
		return emptyList, err
	}

	defer rows.Close()

	emails := make([]EmailEntry, 0, params.Count)

	for rows.Next() {
		email, err := emailEntryFromRow(rows)

		if err != nil {
			return nil, err
		}

		// Here we derefernce above email and appened in the slice
		emails = append(emails, *email)
	}
	return emails, nil
}
