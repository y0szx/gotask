package cli

const (
	help          = "help"
	addContact    = "add"
	deleteContact = "delete"
	listContact   = "list"
	findContact   = "find"
)

type command struct {
	name        string
	description string
}
