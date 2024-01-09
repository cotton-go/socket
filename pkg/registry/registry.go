package registry

type Registry interface {
	Register() error
	Deregister() error
}
