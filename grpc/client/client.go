package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/a2ush/sample-grpc-server-with-appmesh-xray/rpc"
	"google.golang.org/grpc"
)

func main() {

	log.Print("HTTP Server is running with :8080 port...")

	http.HandleFunc("/jst", getJSTTime)
	http.HandleFunc("/utc", getUTCTime)
	http.ListenAndServe(":8080", nil)
}

func getJSTTime(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		marshaled_result, _ := json.Marshal(gRPCRequest("JST"))
		fmt.Fprintf(w, string(marshaled_result))
	}
}

func getUTCTime(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		marshaled_result, _ := json.Marshal(gRPCRequest("UTC"))
		fmt.Fprintf(w, string(marshaled_result))
	}
}

func gRPCRequest(timezone string) *rpc.ServerResponse {

	address := "localhost:50051"
	conn, err := grpc.Dial(
		address,
		grpc.WithInsecure(),
		grpc.WithBlock(),
	)
	if err != nil {
		log.Fatal("Connection failed.")
		return nil
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Second,
	)
	defer cancel()

	client := rpc.NewTimeManageClient(conn)

	req := rpc.ClientRequest{}
	switch timezone {
	case "JST":
		req = rpc.ClientRequest{
			TimezoneFormat: rpc.Timezone_JST,
		}
	case "UTC":
		req = rpc.ClientRequest{
			TimezoneFormat: rpc.Timezone_UTC,
		}
	}

	reply, err := client.ConvertTime(ctx, &req)
	if err != nil {
		log.Fatal("Request failed.")
		return nil
	}

	return reply
}
