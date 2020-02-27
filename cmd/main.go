package main

import (
	"cloud.google.com/go/spanner"
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/lei-at-anz/spanner-lab/pkg/core"
	"log"
)

func main() {
	fmt.Println(uuid.New().String())

	connectionString := "projects/xenon-muse-269523/instances/payments/databases/test"

	ctx := context.Background()

	client, err := spanner.NewClient(ctx, connectionString)
	if err != nil {
		log.Fatal(err)
	}

	if err = core.CleanAllProcessIDs(client, ctx); err != nil {
		log.Fatal(err)
	}
}

