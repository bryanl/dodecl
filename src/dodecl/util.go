package dodecl

import "github.com/pkg/errors"

func extractString(in interface{}) (string, error) {
	s, ok := in.(string)
	if !ok {
		return "", errors.New("unable to extract string from interface")
	}

	return s, nil
}

func extractInt(in interface{}) (int, error) {
	i, ok := in.(int)
	if !ok {
		return 0, errors.New("unable to extract int from interface")
	}

	return i, nil
}

func extractSSHKeys(in interface{}) ([]sshKey, error) {
	out, ok := in.([]interface{})
	if !ok {
		return nil, errors.New("unable to extract key list from interface")
	}

	keys := []sshKey{}

	for _, k := range out {
		switch k.(type) {
		case string:
			keys = append(keys, sshKey{Fingerprint: k.(string)})
		case int:
			keys = append(keys, sshKey{ID: k.(int)})
		default:
			return nil, errors.Errorf("unknown key type found: %v", k)
		}
	}

	return keys, nil
}
