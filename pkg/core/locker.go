package core

import (
	"cloud.google.com/go/spanner"
	"context"
	"fmt"
	"log"
	"time"
)

type LockHolder struct {
	Client *spanner.Client
	ID string
}

func (l *LockHolder) GoLock(ctx context.Context, report chan int) {
	log.Printf("lock holder - %s: started.\n", l.ID)
	for {
		err := l.GrabLock(ctx)
		if err != nil {
			log.Println(err)
			time.Sleep(time.Second * 5)
		} else {
			break
		}
	}
	log.Printf("lock holder - %s: I got the lock!\n", l.ID)
	for i := 0; i < 30; i++ {
		time.Sleep(time.Second * 1)
		if err := l.CheckLock(ctx); err != nil {
			log.Println(err)
		}
	}
	report <- 1
}

func (l *LockHolder) GrabLock(ctx context.Context) error {
	_, err := l.Client.ReadWriteTransaction(ctx, func(c context.Context, transaction *spanner.ReadWriteTransaction) error {
		sql := `UPDATE Payment SET process_id = @processID, process_time = @currentTime 
WHERE id = (SELECT id FROM Payment WHERE process_id IS NULL 
OR process_id = @processID 
OR TIMESTAMP_ADD(process_time, INTERVAL 30 SECOND) < @currentTime 
ORDER BY id LIMIT 1)`
		stmt := spanner.Statement{
			SQL: sql,
			Params: map[string]interface{} {
				"processID": l.ID,
				"currentTime": time.Now(),
			},
		}
		rowCount, err := transaction.Update(c, stmt)
		if err != nil {
			return err
		} else if rowCount != 1 {
			return fmt.Errorf("lock holder - %s: grab lock failed", l.ID)
		}
		return nil
	})
	return err
}

func (l *LockHolder) CheckLock(ctx context.Context) error {
	sql := `SELECT id FROM Payment 
WHERE process_id = @processID AND TIMESTAMP_ADD(process_time, INTERVAL 30 SECOND) >= @currentTime`
	stmt := spanner.Statement{
		SQL: sql,
		Params: map[string]interface{} {
			"processID": l.ID,
			"currentTime": time.Now(),
		},
	}
	rows := l.Client.Single().Query(ctx, stmt)
	_, err := rows.Next()
	if err != nil {
		return fmt.Errorf("lock holder - %s: %v", l.ID, err)
	}
	log.Printf("lock holder - %s: lock checked.\n", l.ID)
	return nil
}
