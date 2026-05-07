package fieldround

import (
	"encoding/json"
	"testing"
)

func decode(t *testing.T, data []byte) map[string]interface{} {
	t.Helper()
	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("decode: %v", err)
	}
	return m
}

func TestNew_EmptyField(t *testing.T) {
	_, err := New("", 2)
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNew_NegativePlaces(t *testing.T) {
	_, err := New("value", -1)
	if err == nil {
		t.Fatal("expected error for negative places")
	}
}

func TestApply_RoundsToPlaces(t *testing.T) {
	r, _ := New("score", 2)
	out, err := r.Apply([]byte(`{"score":3.14159}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := decode(t, out)
	if m["score"].(float64) != 3.14 {
		t.Fatalf("expected 3.14, got %v", m["score"])
	}
}

func TestApply_ZeroPlaces(t *testing.T) {
	r, _ := New("val", 0)
	out, err := r.Apply([]byte(`{"val":2.7}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := decode(t, out)
	if m["val"].(float64) != 3 {
		t.Fatalf("expected 3, got %v", m["val"])
	}
}

func TestApply_FieldAbsent_Unchanged(t *testing.T) {
	r, _ := New("missing", 2)
	input := []byte(`{"other":1.5}`)
	out, err := r.Apply(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(out) != string(input) {
		t.Fatalf("expected unchanged output")
	}
}

func TestApply_NonNumericField_Error(t *testing.T) {
	r, _ := New("name", 2)
	_, err := r.Apply([]byte(`{"name":"alice"}`))
	if err == nil {
		t.Fatal("expected error for non-numeric field")
	}
}

func TestApply_InvalidJSON_Error(t *testing.T) {
	r, _ := New("val", 2)
	_, err := r.Apply([]byte(`not-json`))
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}
