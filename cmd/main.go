package main

import (
	"cloud.google.com/go/spanner"
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/lei-at-anz/spanner-lab/pkg/core"
	"log"
)

const RoutineCount = 10

func main() {
	fmt.Println(uuid.New().String())

	connectionString := "projects/xenon-muse-269523/instances/payments/databases/test"

	ctx := context.Background()

	client, err := spanner.NewClient(ctx, connectionString)
	if err != nil {
		log.Fatal(err)
	}

	count := make(chan int)
	for i := 0; i < RoutineCount; i++ {
		lockHolder := core.LockHolder{
			Client: client,
			ID: core.NewRandomID(),
		}
		go lockHolder.GoLock(ctx, count)
	}

	for i := 0; i < RoutineCount; {
		select {
		case <-count:
			i++
		}
	}
	log.Println("all done.")
}

