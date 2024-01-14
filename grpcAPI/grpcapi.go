package grpcapi

import (
	"database/sql"
	"mailingServer/mailDatabase"
	pb "mailingServer/proto"
	"time"
)

type MailServer struct {
	pb.UnimplementedMailingListServiceServer
	db *sql.DB
}

func pbEntryToMailDatabaseEntry(pbEntry *pb.EmailEntry) mailDatabase.EmailEntry {

	//  we have to convert from epoch time (int 64) to Time Format 
	t := time.Unix(pbEntry.ConfirmedAt, 0)

	emailEntry := mailDatabase.EmailEntry{}

	emailEntry.Id = pbEntry.Id
	emailEntry.Email = pbEntry.Email
	emailEntry.ConfirmedAt = &t
	emailEntry.OptOut = pbEntry.OptOut

	return emailEntry
}

// Inverse of pbEntryToMailDatabaseEntry function
func mailDatabaseEntryToProtoEntry(mdbEntry *mailDatabase.EmailEntry) pb.EmailEntry {

	/*
	emailEntry := pb.EmailEntry{}

	emailEntry.Id = mdbEntry.Id
	emailEntry.Email = mdbEntry.Email
	emailEntry.ConfirmedAt = mdbEntry.ConfirmedAt.Unix() // We changed the time in Unix : Epoch time in int64
	emailEntry.OptOut = mdbEntry.OptOut
	
	return emailEntry
	*/

	return pb.EmailEntry{
		Id : mdbEntry.Id,
		Email : mdbEntry.Email,
		ConfirmedAt : mdbEntry.ConfirmedAt.Unix(), // We changed the time in Unix : Epoch time in int64
		OptOut : mdbEntry.OptOut,
	}

	//Mentioned both approches for better understanding
}
