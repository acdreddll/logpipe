// Package checkpoint tracks the last successfully processed log position
// so that logpipe can resume from where it left off after a restart.
package checkpoint

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
)

// State holds the persisted checkpoint data.
type State struct {
	Offset int64  `json:"offset"`
	Source string `json:"source"`
}

// Checkpoint manages reading and writing checkpoint state to disk.
type Checkpoint struct {
	mu   sync.Mutex
	path string
	state State
}

// New creates a Checkpoint backed by the file at path.
// If the file exists its state is loaded immediately.
func New(path string) (*Checkpoint, error) {
	if path == "" {
		return nil, errors.New("checkpoint: path must not be empty")
	}
	cp := &Checkpoint{path: path}
	if err := cp.load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}
	return cp, nil
}

// Get returns a copy of the current checkpoint state.
func (c *Checkpoint) Get() State {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.state
}

// Save persists the given state to disk atomically.
func (c *Checkpoint) Save(s State) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	data, err := json.Marshal(s)
	if err != nil {
		return err
	}
	tmp := c.path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return err
	}
	if err := os.Rename(tmp, c.path); err != nil {
		return err
	}
	c.state = s
	return nil
}

// Reset clears the persisted state and removes the checkpoint file.
func (c *Checkpoint) Reset() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.state = State{}
	err := os.Remove(c.path)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return err
}

func (c *Checkpoint) load() error {
	data, err := os.ReadFile(c.path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &c.state)
}
