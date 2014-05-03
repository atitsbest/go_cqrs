package application

import (
	"github.com/atitsbest/go_cqrs/storage"
	"github.com/atitsbest/go_cqrs/tests/domain"
)

// CreateCart erstellt einen neuen Warenkorb und speichert ihn.
func CreateCart(name string) error {
	cart := domain.NewCart()
	cart.SetName(name)
	st := storage.NewMemoryStore()
	return st.Save(cart.ID, cart.UncommitedChanges())
}
