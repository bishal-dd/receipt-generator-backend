package resolver

import (
	"github.com/bishal-dd/receipt-generator-backend/graph/resolver/user"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct{
	*user.UserResolver
}
