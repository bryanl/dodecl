package dodecl

// RunBook is a set of actions to be run.
type RunBook interface {
	Name() string
	RunBooks() []RunBook
	Action()
}
