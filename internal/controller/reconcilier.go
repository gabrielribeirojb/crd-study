package controller

import (
	"context"
	"fmt"
	"time"

	"github.com/gabrielribeirojb/crd-study/internal/state"
)

type Reader interface {
	Get(ctx context.Context, namespace, name string) (state.CurrentState, error)
}

type Reconciler struct {
	reader Reader
}

func NewReconciler(reader Reader) *Reconciler {
	return &Reconciler{reader: reader}
}

func (r *Reconciler) Wait(ctx context.Context, desired state.DesiredSpec, interval time.Duration) error {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case <-ticker.C:
			cur, err := r.reader.Get(ctx, desired.Namespace, desired.Name)
			if err != nil {
				return err
			}

			if !cur.Exists {
				fmt.Println("Current: does not exist yet")
				continue
			}

			fmt.Printf("Current phase: %s\n", cur.Phase)

			switch cur.Phase {
			case "SUCCEEDED":
				fmt.Println("Done: restore succeeded")
				return nil
			case "FAILED":
				return fmt.Errorf("restore failed")
			}
		}
	}
}
