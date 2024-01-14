package main

import (
	"context"
	"log"
	pb "mailingServer/proto"
	"time"

	"github.com/alexflint/go-arg"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func logResponse(resp *pb.EmailResponse, err error) {

	if err != nil {
		log.Fatalf(" Error : %v", err)
	}

	if resp.EmailEntry == nil {
		log.Println(" Email not Found")
	} else {
		log.Printf(" Response :%v", resp.EmailEntry)
	}
}

func createEmail(client pb.MailingListServiceClient, email string) *pb.EmailEntry {
	log.Println("create email")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := client.CreateEmail(ctx, &pb.CreateEmailRequest{EmailAddr: email})
	logResponse(res, err)

	return res.EmailEntry
}

func getEmail(client pb.MailingListServiceClient, email string) *pb.EmailEntry {
	log.Println("Get email")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := client.GetEmail(ctx, &pb.GetEmailRequest{EmailAddr: email})
	logResponse(res, err)
	return res.EmailEntry
}

func getEmailBatch(client pb.MailingListServiceClient, count int, page int) {
	log.Println("Get email batch")
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	resEntries, err := client.GetEmailBatch(ctx, &pb.GetEmailBatchRequest{Page: int32(page), Count: int32(count)})

	if err != nil {
		log.Fatalf(" Error : %v", err)
	}

	for index, entry := range resEntries.EmailEntries {
		log.Printf("%v (th) Item : %v ", index, entry)
	}
}

func updateEmail(client pb.MailingListServiceClient, entry *pb.EmailEntry) *pb.EmailEntry {
	log.Println("update email")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := client.UpdateEmail(ctx, &pb.UpdateEmailRequest{EmailEntry: entry})
	logResponse(res, err)
	return res.EmailEntry
}

func deleteEmail(client pb.MailingListServiceClient, email string) *pb.EmailEntry {
	log.Println("delete email")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := client.DeleteEmail(ctx, &pb.DeleteEmailRequest{EmailAddr: email})
	logResponse(res, err)
	return res.EmailEntry
}

var args struct {
	GrpcAddr string `arg:"env:MAILINGLIST_GRPC_ADDR"`
}

func main() {

	arg.MustParse(&args)

	if args.GrpcAddr == "" {
		args.GrpcAddr = ":8081"
	}

	conn, err := grpc.Dial(args.GrpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf(" It is not connected : %v", err)
	}

	defer conn.Close()

	client := pb.NewMailingListServiceClient(conn)

	// Test case 1 : Creating a mail 
	newEmailCreated := createEmail(client, "pankajsaini@gmail.com")
	
	// Test case 2 : updating a mail 
	newEmailCreated.ConfirmedAt = 15081998
	updateEmail(client, newEmailCreated)
	
	//Test case 3 : get All email 
	getEmailBatch(client, 1, 2)

	// Test case 4  : delete the mail
	deleteEmail(client , newEmailCreated.Email)

	//Test case 5 : get All email 
	getEmailBatch(client, 1, 2)
}
