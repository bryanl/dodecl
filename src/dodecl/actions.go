package dodecl

var actions = actionMap{}

type actionMap map[string]resourceAction

type resourceAction func(projectID string, r Resource) (RunBook, error)
