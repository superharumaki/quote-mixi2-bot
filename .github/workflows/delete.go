package main	
	
import (	
	context
	log
	os
	
	github.com/mixigroup/mixi2-application-sdk-go/auth
	application_apiv1 "github.com/mixigroup/mixi2-application-sdk-go/gen/go/social/mixi/application/service/application_api/v1"
	
	google.golang.org/grpc
	google.golang.org/grpc/credentials
)	
	
func main() {	
	authenticator, err := auth.NewAuthenticator(
	
	
	
	)
	if err != nil {
	
	}
	
	authCtx, err := authenticator.AuthorizedContext(context.Background())
	if err != nil {
	
	}
	
	conn, err := grpc.NewClient(
	
	
	
	
	)
	if err != nil {
	
	}
	defer conn.Close()
	
	client := application_apiv1.NewApplicationServiceClient(conn)
	
	postID := "119bacca-1a05-4b68-9111-94429d25c5b6"
	
	_, err = client.DeletePost(authCtx, &application_apiv1.DeletePostRequest{
	
	})
	if err != nil {
	
	}
	
	log.Println("削除成功")
}	
