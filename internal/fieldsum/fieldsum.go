// Package fieldsum provides a processor that sums multiple numeric fields
// from a JSON log event and writes the result to a destination field.
package fieldsum

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// Processor sums a set of source fields and stores the result in a destination field.
type Processor struct {
	sources []string
	dest    string
}

// New creates a Processor that sums the given source fields into dest.
// Returns an error if dest is empty or no source fields are provided.
func New(dest string, sources []string) (*Processor, error) {
	if dest == "" {
		return nil, fmt.Errorf("fieldsum: dest field must not be empty")
	}
	if len(sources) == 0 {
		return nil, fmt.Errorf("fieldsum: at least one source field is required")
	}
	for _, s := range sources {
		if s == "" {
			return nil, fmt.Errorf("fieldsum: source field name must not be empty")
		}
	}
	return &Processor{sources: sources, dest: dest}, nil
}

// Apply reads each source field from the JSON line, sums their numeric values,
// and writes the total to the destination field. Missing fields are treated as 0.
// Returns an error if the input is not valid JSON or a field is non-numeric.
func (p *Processor) Apply(line []byte) ([]byte, error) {
	var obj map[string]interface{}
	if err := json.Unmarshal(line, &obj); err != nil {
		return nil, fmt.Errorf("fieldsum: invalid JSON: %w", err)
	}

	var total float64
	for _, src := range p.sources {
		v, ok := obj[src]
		if !ok {
			continue
		}
		f, err := toFloat(v)
		if err != nil {
			return nil, fmt.Errorf("fieldsum: field %q: %w", src, err)
		}
		total += f
	}

	obj[p.dest] = total

	out, err := json.Marshal(obj)
	if err != nil {
		return nil, fmt.Errorf("fieldsum: marshal: %w", err)
	}
	return out, nil
}

func toFloat(v interface{}) (float64, error) {
	switch x := v.(type) {
	case float64:
		return x, nil
	case string:
		f, err := strconv.ParseFloat(x, 64)
		if err != nil {
			return 0, fmt.Errorf("cannot parse %q as number", x)
		}
		return f, nil
	case bool:
		if x {
			return 1, nil
		}
		return 0, nil
	default:
		return 0, fmt.Errorf("unsupported type %T", v)
	}
}
