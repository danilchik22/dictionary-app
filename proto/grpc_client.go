package proto

import (
	"context"
	pb "dictionary_app/proto/gen"
	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AuthClient struct {
	client pb.AuthServiceClient
}

func NewAuthClient() *AuthClient {
	conn, err := grpc.NewClient(
		"localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatal(err)
	}

	conn.Connect()

	return &AuthClient{
		client: pb.NewAuthServiceClient(conn),
	}

}

func (c *AuthClient) ValidateToken(token string) (response *pb.ResponseToken, err error) {
	res, err := c.client.ValidateToken(context.Background(), &pb.RequestToken{Token: token})
	if err != nil || !res.Valid {
		errFmt := ""
		if err != nil {
			errFmt = fmt.Sprintf("%v", err)
		}
		if res.NeedRefresh == true {
			return &pb.ResponseToken{Username: "", Roles: []string{}, Error: errFmt}, fmt.Errorf("need refresh")
		}
		return &pb.ResponseToken{Username: "", Roles: []string{}, Error: errFmt}, fmt.Errorf("invalid token: %v", err)
	}
	return &pb.ResponseToken{Username: res.Username, Roles: []string{}, Error: ""}, nil
}
