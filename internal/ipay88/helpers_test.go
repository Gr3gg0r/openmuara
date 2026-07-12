package ipay88

import (
	"io"
	"net/http"

	"github.com/Gr3gg0r/openmuara/internal/engine"
)

func readBody(req *http.Request) string {
	if req.Body == nil {
		return ""
	}
	defer func() { _ = req.Body.Close() }()
	data, _ := io.ReadAll(req.Body)
	return string(data)
}

// fakeErrorStore is a TransactionStore that returns configured errors for
// GetByReference and CreateOrGet so error paths can be exercised.
type fakeErrorStore struct {
	getErr    error
	createErr error
}

func (f *fakeErrorStore) GetByReference(string) (engine.Transaction, bool, error) {
	return engine.Transaction{}, false, f.getErr
}

func (f *fakeErrorStore) CreateOrGet(engine.Transaction) (engine.Transaction, bool, error) {
	return engine.Transaction{}, false, f.createErr
}

func (f *fakeErrorStore) GetByID(string) (engine.Transaction, bool, error) {
	return engine.Transaction{}, false, nil
}

func (f *fakeErrorStore) List(int, int) ([]engine.Transaction, error) {
	return nil, nil
}

func (f *fakeErrorStore) Clear() error {
	return nil
}
