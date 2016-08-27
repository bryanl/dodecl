package dodecl

import "github.com/pkg/errors"

// Execer executes a Runbook.
type Execer interface {
	Exec(rb RunBook) error
}

type simpleExecer struct{}

var _ Execer = (*simpleExecer)(nil)

// NewExecer returns an instance of Execer.
func NewExecer() Execer {
	return &simpleExecer{}
}

func (e *simpleExecer) Exec(rb RunBook) error {
	rb.Action()

	for _, child := range rb.RunBooks() {
		err := e.Exec(child)
		if err != nil {
			return errors.Wrapf(err, "failure running node %s", rb.Name)
		}
	}

	return nil
}
