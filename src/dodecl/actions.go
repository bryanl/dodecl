package dodecl

var actions = map[string]resourceAction{}

type actionMap map[string]resourceAction

type resourceAction func(Resource) (RunBook, error)
