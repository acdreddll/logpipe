// Package jsonpath provides dot-notation field extraction and mutation
// for JSON log lines, supporting nested key traversal.
package jsonpath

import (
	"encoding/json"
	"errors"
	"strings"
)

// ErrNotFound is returned when the requested path does not exist.
var ErrNotFound = errors.New("jsonpath: field not found")

// Get extracts the value at the dot-separated path from a JSON object.
// Returns ErrNotFound if any segment along the path is absent.
func Get(line []byte, path string) (any, error) {
	var root map[string]any
	if err := json.Unmarshal(line, &root); err != nil {
		return nil, err
	}
	return walk(root, strings.Split(path, "."))
}

// Set returns a new JSON line with the value at path set to val.
// Intermediate maps are created as needed.
func Set(line []byte, path string, val any) ([]byte, error) {
	var root map[string]any
	if err := json.Unmarshal(line, &root); err != nil {
		return nil, err
	}
	setNested(root, strings.Split(path, "."), val)
	return json.Marshal(root)
}

// Delete returns a new JSON line with the field at path removed.
// A missing path is silently ignored.
func Delete(line []byte, path string) ([]byte, error) {
	var root map[string]any
	if err := json.Unmarshal(line, &root); err != nil {
		return nil, err
	}
	deleteNested(root, strings.Split(path, "."))
	return json.Marshal(root)
}

func walk(m map[string]any, parts []string) (any, error) {
	v, ok := m[parts[0]]
	if !ok {
		return nil, ErrNotFound
	}
	if len(parts) == 1 {
		return v, nil
	}
	child, ok := v.(map[string]any)
	if !ok {
		return nil, ErrNotFound
	}
	return walk(child, parts[1:])
}

func setNested(m map[string]any, parts []string, val any) {
	if len(parts) == 1 {
		m[parts[0]] = val
		return
	}
	child, ok := m[parts[0]].(map[string]any)
	if !ok {
		child = make(map[string]any)
		m[parts[0]] = child
	}
	setNested(child, parts[1:], val)
}

func deleteNested(m map[string]any, parts []string) {
	if len(parts) == 1 {
		delete(m, parts[0])
		return
	}
	child, ok := m[parts[0]].(map[string]any)
	if !ok {
		return
	}
	deleteNested(child, parts[1:])
}
