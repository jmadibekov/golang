package main

import (
	"context"
	"example/hello/grpc/api"
	"log"
	"time"

	"google.golang.org/grpc"
)

const (
	port = ":8080"
)

func main() {
	ctx := context.Background()

	connStartTime := time.Now()
	conn, err := grpc.Dial("localhost"+port, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("could not connect to %s: %v", port, err)
	}
	log.Printf("connected in %d microsec", time.Since(connStartTime).Microseconds())

	mailSenderClient := api.NewMailServiceClient(conn)
	personalAccounterClient := api.NewPersonalAccountServiceClient(conn)

	mailSendStartTime := time.Now()
	_, err = mailSenderClient.MailSend(ctx, &api.MailSendRequest{
		To:      "orochimaru@mail.ru",
		Message: "please, send password",
	})
	if err != nil {
		log.Fatalf("could not send mail: %v", err)
	}
	log.Printf("sended mail in %d microsec", time.Since(mailSendStartTime).Microseconds())

	validAccountID := int64(3)
	invalidAccountID := int64(123)

	validAccount, err := personalAccounterClient.PersonalAccount(ctx, &api.PersonalAccountRequest{Id: validAccountID})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("got account: %+v", validAccount)

	_, err = personalAccounterClient.PersonalAccount(ctx, &api.PersonalAccountRequest{Id: invalidAccountID})
	if err != nil {
		log.Printf("got err: %v", err)
	}
}
