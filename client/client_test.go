package client

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"testing"
)

const (
	target = "127.0.0.1:8080"
)

func TestClient_SetTentacle(t *testing.T) {
	client, err := NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	defer client.Close(context.Background())

	// TODO:
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
}

func TestClient_GetTentacle(t *testing.T) {
	client, err := NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	defer client.Close(context.Background())

	// TODO:
}

func TestClient_DelTentacle(t *testing.T) {
	client, err := NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	defer client.Close(context.Background())

	// TODO:
}

func TestClient_BatchGetTentacle(t *testing.T) {
	client, err := NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	defer client.Close(context.Background())

	// TODO:
}

func TestClient_BatchSetTentacle(t *testing.T) {
	client, err := NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	defer client.Close(context.Background())

	// TODO:
}

func TestClient_BatchDelTentacle(t *testing.T) {
	client, err := NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	defer client.Close(context.Background())

	// TODO:
}
