package manager

import (
	"errors"
	//"fmt"
	//"os"
)

type Manager struct {
	Chan chan *Instance
	instances map[string]Instance
}

func New() *Manager {
	return &Manager{
		Chan: make(chan *Instance),
		instances: make(map[string]Instance),
	}
}

func (self Manager) registerInstance(instance *Instance) (err error) {
	name := instance.Name
	if _, ok := self.instances[name]; ok {
		err = errors.New("Instance already registered")
	} else {
		self.instances[instance.Name] = *instance
	}

	return
}

func (self Manager) Handle(instance *Instance) error {
	if err := self.registerInstance(instance); err != nil {
		return err
	}

	if err, ok := instance.start(); err != nil {
		return err
	} else if !ok {
		return errors.New("Unknown error")
	}

	go func() {
		for {
			instance.wait()
			self.Chan <- instance
			if _, ok := instance.start(); !ok {
				return
			}
		}
	}()

	return nil
}

func (self Manager) GetInstances() map[string]Instance {
	return self.instances
}
