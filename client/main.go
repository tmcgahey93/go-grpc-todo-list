package main

import (
	"context"
	//"crypto/tls"
	"fmt"
	todo "go-grpc-todo-list/go-grpc-todo/proto"
	"log"

	"google.golang.org/grpc"
	//"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	//creds := credentials.NewTLS(&tls.Config{})
	//conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(creds))

	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	client := todo.NewTodoServiceClient(conn)

	//add task
	task := &todo.Task{
		Id:          "1",
		Description: "Learn gRPC in Go",
		Done:        false,
	}
	resp, _ := client.AddTask(context.Background(), task)

	fmt.Println("AddTask response:", resp.Message)

	//list tasks
	list, _ := client.ListTasks(context.Background(), &todo.Empty{})
	for _, t := range list.Tasks {
		fmt.Printf("Task: %s (%s)\n", t.Description, t.Id)
	}

	//stream tasks
	fmt.Println("Streaming tasks...")
	stream, _ := client.StreamTasks(context.Background(), &todo.Empty{})
	for {
		t, err := stream.Recv()
		if err != nil {
			break
		}
		fmt.Printf("Streamed: %s (%s)\n", t.Description, t.Id)
	}

	//delete task
	delResp, _ := client.DeleteTask(context.Background(), &todo.TaskID{Id: "1"})
	fmt.Println("DeleteTask response:", delResp.Message)
}
