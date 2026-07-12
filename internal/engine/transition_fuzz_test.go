package engine

import (
	"testing"
)

func FuzzTransition(f *testing.F) {
	statuses := []TransactionStatus{
		TransactionStatusNew,
		TransactionStatusPaid,
		TransactionStatusUnpaid,
		TransactionStatusRefunded,
	}
	for _, from := range statuses {
		for _, to := range statuses {
			f.Add(string(from), string(to))
		}
	}

	f.Fuzz(func(t *testing.T, fromStr, toStr string) {
		from := TransactionStatus(fromStr)
		to := TransactionStatus(toStr)

		tx := Transaction{Status: from}
		err := Transition(&tx, to)

		if from == to {
			if err != nil {
				t.Errorf("expected %s -> %s to be a no-op, got %v", from, to, err)
			}
			if tx.Status != to {
				t.Errorf("expected status %q, got %q", to, tx.Status)
			}
			return
		}

		if CanTransition(from, to) {
			if err != nil {
				t.Errorf("expected %s -> %s to succeed, got %v", from, to, err)
			}
			if tx.Status != to {
				t.Errorf("expected status %q, got %q", to, tx.Status)
			}
			return
		}

		if err == nil {
			t.Errorf("expected %s -> %s to fail", from, to)
		}
	})
}
