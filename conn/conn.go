package conn

import "context"

type Conn interface {
	Name() string
	State() string
	Check(ctx context.Context) error
}
