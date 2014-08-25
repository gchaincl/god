package manager

import(
	"time"
	"testing"
)

func assert(t *testing.T, expr bool, format string, args ...interface{}) {
	if expr == false {
		t.Fatalf(format, args...)
	}
}

func TestHandleNotifyWhenProcessFinish(t *testing.T) {
	m := New()
	i := NewInstance("ls", nil)
	m.Handle(i)

	select {
	case _ = <- time.After(1 * time.Second):
		t.Fatal("Notification wait expired")
	case _ = <- m.Chan:
	}
}

func TestHandleWithInvalidCommand(t *testing.T) {
	m := New()
	i := NewInstance("_INVALID_CMD_", nil)
	err := m.Handle(i)
	assert(t, err != nil, "Should return an error")
}

func TestHandleRegisterInstance(t *testing.T) {
	m := New()
	i := NewInstance("ls", nil)
	m.Handle(i)
	<- m.Chan

	instances := m.GetInstances()

	assert(t, len(instances) == 1, "Instance not registered")

	expected := instances["ls"];
	assert(t, i.Name == expected.Name,
		"Failed asserting %#v == %#v", i, &expected)
}

func TestHandleSetsStartAndEndTimestamps(t *testing.T) {
	m := New()
	i := NewInstance("ls", nil)
	i.Max = 1
	m.Handle(i)

	<- m.Chan
	assert(t, i.EndedAt.UnixNano() > i.StartedAt.UnixNano(),
		"Failed asserting %#v > %#v", i.EndedAt, i.StartedAt)
}

func TestHandleFailsRegisteringTheSameInstance(t *testing.T) {
	m := New()
	i := NewInstance("ls", nil)

	if err := m.Handle(i); err != nil { t.Fatal(err) }

	err := m.Handle(i);
	assert(t, err != nil, "Failed asserting err != nil")
}

func TestHandleRestartsInstanceExecutionNTimes(t *testing.T) {
	m := New()
	i := NewInstance("ls", nil)
	max := i.Max
	m.Handle(i);

	for c := 0; c < max; c++ {
		select {
		case _ = <- time.After(1 * time.Second):
			t.Fatal("Notification wait expired")
		case _ = <- m.Chan:
		}
	}

	assert(t, i.Attempts() == 3,
		"Failed asserting 3 == %d", i.Attempts())
}

func TestHandleInstanceMax0(t *testing.T) {
	m := New()
	i := NewInstance("ls", nil)
	i.Max = 0
	x := m.Handle(i);

	assert(t, x != nil, "should return an error")

}
