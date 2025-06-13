package main

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

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
		fmt.Printf("Task: %s (%s)\n", t.Description, t.ID)
	}

	//stream tasks
	fmt.Println("Streaming tasks...")
	stream, _ := client.StreamTasks(context.Background(), &todo.Empty{})
	for {
		t, err := stream.Recv()
		if err != nil {
			break
		}
		fmt.Printf("Streamed: %s (%s)\n", t.Description, t.ID)
	}

	//delete task
	delResp, _ := client.DeleteTask(context.Background(), &todo.TaskID{Id: "1"})
	fmt.Println("DeleteTask response:", delResp.Message)
}
