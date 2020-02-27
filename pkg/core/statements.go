package core

import (
	"cloud.google.com/go/spanner"
	"context"
	"fmt"
	"log"
)

const (
	tableNamePayment = "Payment"
	colPaymentID = "id"
	colPaymentAmount = "amount"
	colPaymentState = "state"
	colPaymentProcessID = "process_id"
	colPaymentProcessTime = "process_time"
)

var (
	colsPayment = []string{colPaymentID, colPaymentAmount, colPaymentState}
)

func InsertNewPayment(client *spanner.Client, ctx context.Context, process_id string) error {
	values := []interface{} {
		NewRandomID(),
		process_id,
		"NEW",
	}
	mutation := spanner.Insert(tableNamePayment, colsPayment, values)
	_, err := client.Apply(ctx, []*spanner.Mutation {mutation})
	return err
}

func CleanAllProcessIDs(client *spanner.Client, ctx context.Context) error {
	sql := fmt.Sprintf("UPDATE %s SET %s = NULL, %s = NULL WHERE 1 = 1", tableNamePayment, colPaymentProcessID, colPaymentProcessTime)
	stmt := spanner.NewStatement(sql)
	_, err := client.ReadWriteTransaction(ctx, func(ctx context.Context, transaction *spanner.ReadWriteTransaction) error {
		rowCount, err := transaction.Update(ctx, stmt)
		if err != nil {
			return err
		}
		log.Println(rowCount, " rows updated.")
		return nil
	})
	return err
}
