package dodecl

import (
	"fmt"

	"github.com/pkg/errors"
)

type runBook struct {
	name     string
	runBooks []RunBook
	action   string
}

var _ RunBook = (*runBook)(nil)

func (r *runBook) Name() string {
	return r.name
}

func (r *runBook) RunBooks() []RunBook {
	return r.runBooks
}

func (r *runBook) Action() {
	fmt.Printf("%s: %s\n", r.name, r.action)
}

// Planner creates a plan of execution given a Dodecl description.
type Planner interface {
	Plan(d *Dodecl) (RunBook, error)
}

type simplePlanner struct{}

var _ Planner = (*simplePlanner)(nil)

// NewPlanner returns an instance of Planner.
func NewPlanner() Planner {
	return &simplePlanner{}
}

func (p *simplePlanner) Plan(d *Dodecl) (RunBook, error) {
	rb := &runBook{
		name:     "root",
		runBooks: []RunBook{},
	}

	for _, resource := range d.Resources {
		action, ok := actions[resource.Type]
		if !ok {
			return nil, errors.Errorf("unknown action type: %s",
				resource.Type)
		}

		resourceRunBook, err := action(d.ID, resource)
		if err != nil {
			return nil, errors.Wrapf(err,
				"resource error for %s", resource.Type)
		}

		rb.runBooks = append(rb.runBooks, resourceRunBook)
	}

	return rb, nil
}
