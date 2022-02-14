package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/a2ush/sample-grpc-server-with-appmesh-xray/rpc"
	"google.golang.org/grpc"
)

func main() {

	log.Print("HTTP Server is running with :8080 port...")

	http.HandleFunc("/", getTime)
	http.ListenAndServe(":8080", nil)
}

func getTime(w http.ResponseWriter, r *http.Request) {

	log.Println(r)
	if r.Method == http.MethodGet {

		var marshaled_result []byte

		switch r.RequestURI {
		case "/utc":
			marshaled_result, _ = json.Marshal(gRPCRequest("UTC"))
		case "/pst":
			marshaled_result, _ = json.Marshal(gRPCRequest("PST"))
		case "/jst":
			marshaled_result, _ = json.Marshal(gRPCRequest("JST"))
		default:
			marshaled_result, _ = json.Marshal(gRPCRequest("UTC"))
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write(marshaled_result)
	}
}

func gRPCRequest(timezone string) *rpc.ServerResponse {

	path, found := os.LookupEnv("SERVER_PATH")
	if !found {
		path = "localhost"
	}
	address := path + ":50051"

	conn, err := grpc.Dial(
		address,
		grpc.WithInsecure(),
		grpc.WithBlock(),
	)
	if err != nil {
		// log.Fatal("Connection failed.")
		log.Println("Connection failed.")
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
	case "PST":
		req = rpc.ClientRequest{
			TimezoneFormat: rpc.Timezone_PST,
		}
	}

	result, err := client.ConvertTime(ctx, &req)
	log.Println(result)
	log.Println(err)
	if err != nil {
		// log.Fatal("Request failed.")
		log.Println("Request failed.")
		return nil
	}

	return result
}
