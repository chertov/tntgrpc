package main

import (
    pb_helloworld "helloworld"
    "log"
    "google.golang.org/grpc"
    "context"
    "time"
    "encoding/json"
)

func main() {
    address := "127.0.0.1:50051"
    conn, err := grpc.Dial(address, grpc.WithInsecure())
    if err != nil {
        log.Fatalf("did not connect: %v", err)
    }
    defer conn.Close()
    client := pb_helloworld.NewExampleServiceClient(conn)

    ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
    defer cancel()

    req := pb_helloworld.HelloRequest{"World"}
    reply, err := client.SayHello(ctx, &req)
    if err != nil {
        log.Fatalf("could not to say hello: %v", err)
    }
    replyjson, _ := json.MarshalIndent(reply, "", "    ")
    log.Println("Reply:", string(replyjson))
}
