package events

type (
	CartCreated struct{}

	CartNameChanged struct {
		Name string
	}
)
