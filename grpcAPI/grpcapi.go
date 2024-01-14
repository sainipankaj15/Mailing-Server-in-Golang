package grpcapi

import (
	"context"
	"database/sql"
	"log"
	"mailingServer/mailDatabase"
	pb "mailingServer/proto"
	"net"
	"time"

	"google.golang.org/grpc"
)

type MailServer struct {
	pb.UnimplementedMailingListServiceServer
	db *sql.DB
}

// Proto Entry to Mail database entry
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

// Inverse of pbEntryToMailDatabaseEntry function : Mail Database entry to proto entry
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
		Id:          mdbEntry.Id,
		Email:       mdbEntry.Email,
		ConfirmedAt: mdbEntry.ConfirmedAt.Unix(), // We changed the time in Unix : Epoch time in int64
		OptOut:      mdbEntry.OptOut,
	}

	//Mentioned both approches for better understanding
}

func emailResponse(db *sql.DB, email string) (*pb.EmailResponse, error) {

	entry, err := mailDatabase.GetEmailFromDB(db, email)

	if err != nil {
		return &pb.EmailResponse{}, err
	}
	if entry == nil {
		return &pb.EmailResponse{}, nil
	}

	resp := mailDatabaseEntryToProtoEntry(entry)
	return &pb.EmailResponse{EmailEntry: &resp}, nil
}

func (s *MailServer) GetEmail(ctx context.Context, req *pb.GetEmailRequest) (*pb.EmailResponse, error) {
	log.Printf("gRPC GetEmail : %v\n", req)
	return emailResponse(s.db, req.EmailAddr)
}

func (s *MailServer) GetEmailBatch(ctx context.Context, req *pb.GetEmailBatchRequest) (*pb.GetEmailBatchResponse, error) {
	log.Printf("gRPC GetEmailBatch : %v\n", req)

	paramas := mailDatabase.GetEmailBatchQueryParams{
		Page:  int(req.Page),
		Count: int(req.Count),
	}

	mailEntries, err := mailDatabase.GetEmailBatchFromDB(s.db, paramas)
	if err != nil {
		return &pb.GetEmailBatchResponse{}, nil
	}

	pbEntries := make([]*pb.EmailEntry, 0, len(mailEntries))
	for _, entry := range mailEntries {
		pbentry := mailDatabaseEntryToProtoEntry(&entry)
		pbEntries = append(pbEntries, &pbentry)
	}

	return &pb.GetEmailBatchResponse{EmailEntries: pbEntries}, nil
}

func (s *MailServer) CreateEmail(ctx context.Context, req *pb.CreateEmailRequest) (*pb.EmailResponse, error) {
	log.Printf("gRPC CreateEmail : %v\n", req)

	err := mailDatabase.CreateEmail(s.db, req.EmailAddr)

	if err != nil {
		return &pb.EmailResponse{}, err
	}

	return emailResponse(s.db, req.EmailAddr)
}

func (s *MailServer) UpdateEmail(ctx context.Context, req *pb.UpdateEmailRequest) (*pb.EmailResponse, error) {
	log.Printf("gRPC UpdateEmail : %v\n", req)

	entry := pbEntryToMailDatabaseEntry(req.EmailEntry)
	err := mailDatabase.UpdateEmail(s.db, entry)

	if err != nil {
		return &pb.EmailResponse{}, err
	}

	return emailResponse(s.db, entry.Email)
}

func (s *MailServer) DeleteEmail(ctx context.Context, req *pb.DeleteEmailRequest) (*pb.EmailResponse, error) {
	log.Printf("gRPC UpdateEmail : %v\n", req)

	err := mailDatabase.DeleteEmail(s.db, req.EmailAddr)

	if err != nil {
		return &pb.EmailResponse{}, err
	}

	return emailResponse(s.db, req.EmailAddr)
}

func Serve(db *sql.DB, bind string) {

	listener, err := net.Listen("tcp", bind)

	if err != nil {
		log.Fatalf("gRPC server error : failur on bind %v and error : %v\n", bind, err)
	}

	grpcServer := grpc.NewServer()

	mailServer := MailServer{db: db}

	pb.RegisterMailingListServiceServer(grpcServer, &mailServer)

	log.Printf("gRPC API server listening on %v\n", bind)

	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatalf("gRPC server error : %v\n", err)
	}
}
