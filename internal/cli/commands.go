package cli

const (
	help         = "help"
	acceptOrder  = "acceptOrder"
	returnOrder  = "returnOrder"
	issueOrder   = "issueOrder"
	listOrders   = "listOrders"
	acceptReturn = "acceptReturn"
	listReturns  = "listReturns"
	exit         = "exit"
)

type command struct {
	name        string
	description string
}
