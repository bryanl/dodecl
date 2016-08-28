package state

import (
	"dodecl/do"
	"fmt"

	"github.com/digitalocean/godo"
	"github.com/pkg/errors"
)

// App is where app state is managed.
var App State

// State is an interface for storing state. It could be
// a remote data store, or a local file.
type State interface {
	Verify() error
}

func init() {
	App = &LocalState{}
}

// LocalState is State hosted locally in a $HOME/.local/share/dodecl.
type LocalState struct {
}

var _ State = (*LocalState)(nil)

// Verify checks for local state.
func (s *LocalState) Verify() error {
	fmt.Println("verifying local state")
	return nil
}

// DOState is State hosted on DigitalOcean.
type DOState struct {
	doClient *godo.Client
}

var _ State = (*DOState)(nil)

// Verify verifies the presence of the state do resource.
func (s *DOState) Verify() error {
	hasState, err := s.hasState()
	if err != nil {
		return errors.Wrap(err, "unable to check for state tag")
	}

	if hasState {
		return nil
	}

	fmt.Println("creating state resource")
	return nil
}

// hasState returns true if the state resource is created.
func (s *DOState) hasState() (bool, error) {
	ds := do.NewDropletsService(s.doClient)
	droplets, err := ds.ListByTag("dodecl:state")
	if err != nil {
		return false, errors.Wrap(err, "unable to list droplets by tag")
	}

	return len(droplets) > 0, nil
}
