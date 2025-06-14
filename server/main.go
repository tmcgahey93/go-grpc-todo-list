package main

import (
	"context"
	"fmt"
	todo "go-grpc-todo-list/go-grpc-todo/proto"
	"log"
	"net"
	"sync"

	"google.golang.org/grpc"
)

type server struct {
	todo.UnimplementedTodoServiceServer
	mu    sync.Mutex
	tasks map[string]*todo.Task
}

func (s *server) AddTask(ctx context.Context, task *todo.Task) (*todo.TaskResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.tasks[task.Id] = task
	return &todo.TaskResponse{Message: "Task added!"}, nil
}

func (s *server) ListTasks(ctx context.Context, _ *todo.Empty) (*todo.TaskList, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	var list []*todo.Task

	for _, t := range s.tasks {
		list = append(list, t)
	}

	return &todo.TaskList{Tasks: list}, nil
}

func (s *server) DeleteTask(ctx context.Context, id *todo.TaskID) (*todo.TaskResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.tasks, id.Id)

	return &todo.TaskResponse{Message: "Task deleted!"}, nil
}

func (s *server) StreamTasks(_ *todo.Empty, stream todo.TodoService_StreamTasksServer) error {
	s.mu.Lock()
	tasks := make([]*todo.Task, 0, len(s.tasks))

	for _, task := range s.tasks {
		tasks = append(tasks, task)
	}

	s.mu.Unlock()

	for _, task := range tasks {
		if err := stream.Send(task); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	todo.RegisterTodoServiceServer(s, &server{tasks: make(map[string]*todo.Task)})

	fmt.Println("Server listening on :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
