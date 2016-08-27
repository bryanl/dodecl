package dodecl

import "fmt"

// FloatingIPKey identifies floating ip actions.
const FloatingIPKey = "compute.v2.floating-ip"

func init() {
	actions[FloatingIPKey] = fipAction
}

func fipAction(r Resource) (RunBook, error) {
	return &runBook{
		name:   "create fip",
		action: fmt.Sprintf("region: %s", r.Properties["region"]),
	}, nil
}
