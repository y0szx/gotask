package cli

const (
	help         = "help"
	acceptOrder  = "acceptOrder"
	returnOrder  = "returnOrder"
	issueOrder   = "issueOrder"
	listOrders   = "listOrders"
	acceptReturn = "acceptReturn"
	listReturns  = "listReturns"
)

type command struct {
	name        string
	description string
}
