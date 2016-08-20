package dodecl

import (
	"fmt"

	"github.com/pkg/errors"
)

var (
	actions = map[string]resourceAction{
		"compute.v2.droplet":     dropletAction,
		"compute.v2.floating-ip": fipAction,
	}
)

type resourceAction func(Resource) (*RunBook, error)

// RunBook is a set of actions to be run.
type RunBook struct {
	Name     string
	RunBooks []RunBook
	Action   string
}

// Planner creates a plan of execution given a Dodecl description.
type Planner interface {
	Plan(d *Dodecl) (*RunBook, error)
}

type simplePlanner struct{}

var _ Planner = (*simplePlanner)(nil)

// NewPlanner returns an instance of Planner.
func NewPlanner() Planner {
	return &simplePlanner{}
}

func (p *simplePlanner) Plan(d *Dodecl) (*RunBook, error) {
	rb := &RunBook{
		Name:     "root",
		RunBooks: []RunBook{},
	}

	for _, resource := range d.Resources {
		action, ok := actions[resource.Type]
		if !ok {
			return nil, errors.Errorf("unknown action type: %s",
				resource.Type)
		}

		resourceRunBook, err := action(resource)
		if err != nil {
			return nil, errors.Wrapf(err,
				"resource error for %s", resource.Type)
		}

		rb.RunBooks = append(rb.RunBooks, *resourceRunBook)
	}

	return rb, nil
}

func dropletAction(r Resource) (*RunBook, error) {
	return &RunBook{
		Name: "create droplet",
		Action: fmt.Sprintf("region: %s, count: %d, image: %s, size: %s, keys: %+v",
			r.Properties["region"], r.Properties["count"], r.Properties["image"],
			r.Properties["size"], r.Properties["keys"]),
	}, nil
}

func fipAction(r Resource) (*RunBook, error) {
	return &RunBook{
		Name:   "create fip",
		Action: fmt.Sprintf("region: %s", r.Properties["region"]),
	}, nil
}
