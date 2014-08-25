package manager

import (
//	"fmt"
	"os/exec"
	"time"
)

type Instance struct {
	Name string
	Args []string

	// Register last Started and End time
	StartedAt time.Time
	EndedAt time.Time

	// Consider normal run after Threshold (in seconds)
	Threshold int

	// Max number attempts to restart Process
	Max int

	cmd *exec.Cmd
	attempts int
}

func (self *Instance) start() (err error, ok bool) {
	ok = false
	err = nil

	if self.attempts == self.Max {
		return
	}

	self.cmd = exec.Command(self.Name, self.Args...)
	err = self.cmd.Start();
	if err == nil {
		ok = true
		self.StartedAt = time.Now()
	}

	self.attempts += 1

	return
}

func (self *Instance) wait() {
	err := self.cmd.Wait()
	if err != nil {
		panic(err)
	}
	self.EndedAt = time.Now()
}

func (self *Instance) Attempts() int {
	return self.attempts
}

func NewInstance(name string, args []string) *Instance {
	return &Instance{
		Name: name,
		Args: args,
		Threshold: 1,
		Max: 3,
	}
}
