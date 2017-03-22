package vault

import (
	"log"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/hashicorp/go-uuid"
)

// mockRollback returns a mock rollback manager
func mockRollback(t *testing.T) (*RollbackManager, *NoopBackend) {
	backend := new(NoopBackend)
	mounts := new(MountTable)
	router := NewRouter()

	mounts.Entries = []*MountEntry{
		&MountEntry{
			Path: "foo",
		},
	}
	meUUID, err := uuid.GenerateUUID()
	if err != nil {
		t.Fatal(err)
	}
	if err := router.Mount(backend, "foo", &MountEntry{UUID: meUUID}, nil); err != nil {
		t.Fatalf("err: %s", err)
	}

	mountsFunc := func() []*MountEntry {
		return mounts.Entries
	}

	logger := log.New(os.Stderr, "", log.LstdFlags)
	rb := NewRollbackManager(logger, mountsFunc, router)
	rb.period = 10 * time.Millisecond
	return rb, backend
}

func TestRollbackManager(t *testing.T) {
	m, backend := mockRollback(t)
	if len(backend.Paths) > 0 {
		t.Fatalf("bad: %#v", backend)
	}

	m.Start()
	time.Sleep(50 * time.Millisecond)
	m.Stop()

	count := len(backend.Paths)
	if count == 0 {
		t.Fatalf("bad: %#v", backend)
	}
	if backend.Paths[0] != "" {
		t.Fatalf("bad: %#v", backend)
	}

	time.Sleep(50 * time.Millisecond)

	if count != len(backend.Paths) {
		t.Fatalf("should stop requests: %#v", backend)
	}
}

func TestRollbackManager_Join(t *testing.T) {
	m, backend := mockRollback(t)
	if len(backend.Paths) > 0 {
		t.Fatalf("bad: %#v", backend)
	}

	m.Start()
	defer m.Stop()

	wg := &sync.WaitGroup{}
	wg.Add(3)

	go func() {
		defer wg.Done()
		err := m.Rollback("foo")
		if err != nil {
			t.Fatalf("err: %v", err)
		}
	}()

	go func() {
		defer wg.Done()
		err := m.Rollback("foo")
		if err != nil {
			t.Fatalf("err: %v", err)
		}
	}()

	go func() {
		defer wg.Done()
		err := m.Rollback("foo")
		if err != nil {
			t.Fatalf("err: %v", err)
		}
	}()
	wg.Wait()
}
