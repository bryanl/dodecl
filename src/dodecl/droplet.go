package dodecl

import "fmt"

// DropletKey identifies droplet actions.
const DropletKey = "compute.v2.droplet"

func init() {
	actions[DropletKey] = dropletAction
}

// DropletSettings are settings for a droplet.
type DropletSettings struct {
	region string
	image  string
	size   string
	keys   []sshKey
}

type sshKey struct {
	ID          int
	Fingerprint string
}

// AddDropletRunBook is a runbook for adding droplets.
type AddDropletRunBook struct {
	dropletSettings *DropletSettings
	runBooks        []RunBook
}

var _ RunBook = (*AddDropletRunBook)(nil)

// NewAddDropletRunBook creates an instance of AddDropletRunBook.
func NewAddDropletRunBook(r Resource) (*AddDropletRunBook, error) {
	region, err := extractString(r.Properties["region"])
	if err != nil {
		return nil, err
	}

	image, err := extractString(r.Properties["image"])
	if err != nil {
		return nil, err
	}

	size, err := extractString(r.Properties["size"])
	if err != nil {
		return nil, err
	}

	count, err := extractInt(r.Properties["count"])
	if err != nil {
		return nil, err
	}

	keys, err := extractSSHKeys(r.Properties["keys"])
	if err != nil {
		return nil, err
	}

	rb := AddDropletRunBook{
		dropletSettings: &DropletSettings{
			region: region,
			image:  image,
			size:   size,
			keys:   keys,
		},
		runBooks: []RunBook{},
	}

	for i := 1; i <= count; i++ {
		name := createDropletName(r, i)
		child := newCreateDropletRunBook(rb.dropletSettings, name)
		rb.runBooks = append(rb.runBooks, child)
	}

	return &rb, nil
}

func createDropletName(r Resource, i int) string {
	return fmt.Sprintf("%s-%d", r.Name, i)
}

// Name is the name for the runbook.
func (r *AddDropletRunBook) Name() string {
	return "create droplet"
}

// RunBooks are the child runbooks.
func (r *AddDropletRunBook) RunBooks() []RunBook {
	return r.runBooks
}

// Action performs this runbook's action.
func (r *AddDropletRunBook) Action() {
	fmt.Printf("creating droplets with settings: "+
		"region: %s, image: %s, size: %s, keys: %#v\n",
		r.dropletSettings.region, r.dropletSettings.image,
		r.dropletSettings.size, r.dropletSettings.keys)
}

type createDropletRunBook struct {
	dropletSettings *DropletSettings
	name            string
}

func newCreateDropletRunBook(ds *DropletSettings, name string) *createDropletRunBook {
	return &createDropletRunBook{
		dropletSettings: ds,
		name:            name,
	}
}

var _ RunBook = (*createDropletRunBook)(nil)

func (r *createDropletRunBook) Name() string {
	return fmt.Sprintf("create droplet: %s", r.name)
}

func (r *createDropletRunBook) RunBooks() []RunBook {
	return []RunBook{}
}

func (r *createDropletRunBook) Action() {
	fmt.Println(r.Name())
}

func dropletAction(r Resource) (RunBook, error) {
	return NewAddDropletRunBook(r)
}
