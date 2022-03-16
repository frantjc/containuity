package plugin

import "sync"

var (
	registry Registry = Registry{}
)

func Register(r *Registration) {
	if err := registry.Register(r); err != nil {
		panic(err)
	}
}

type Registry struct {
	sync.RWMutex
	registrations []Registration
}

func (r *Registry) ContainsID(t Type, id string) bool {
	for _, registration := range r.registrations {
		if registration.Type == t && registration.ID == id {
			return true
		}
	}

	return false
}

func (r *Registry) Contains(t Type) bool {
	for _, registration := range r.registrations {
		if registration.Type == t {
			return true
		}
	}

	return false
}

func (reg *Registry) Register(r *Registration) error {
	reg.Lock()
	defer reg.Unlock()

	if r.Type == "" {
		return ErrNoType
	}
	if r.ID == "" {
		return ErrNoPluginID
	}
	if yes := reg.Contains(r.Type); yes {
		return ErrIDRegistered
	}

	reg.registrations = append(reg.registrations, *r)
	return nil
}

func (reg *Registry) Graph() []Registration {
	reg.Lock()
	defer reg.Unlock()

	var (
		ordered = make([]Registration, len(reg.registrations))
		added   = map[Type]bool{}
	)
	for i := range ordered {
		for _, registration := range reg.registrations {
			if yes, ok := added[registration.Type]; !(ok && yes) {
				satisfied := true
				for _, dep := range registration.DependsOn {
					yes, ok := added[dep]
					satisfied = satisfied && ok && yes
				}

				if satisfied {
					ordered[i] = registration
					added[registration.Type] = true
					break
				}
			}
		}
	}

	return ordered
}
