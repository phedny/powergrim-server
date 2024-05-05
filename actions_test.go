package main_test

import (
	"encoding/json"
	"testing"

	powergrim "github.com/phedny/powergrim-server"
)

func TestMarshalAction(t *testing.T) {
	data, err := json.Marshal(powergrim.AddPlayer{
		Action: "addPlayer",
		Id:     1,
	})
	if err != nil {
		t.Fatal(err)
	}
	expected := `{"action":"addPlayer","id":1}`
	if string(data) != expected {
		t.Fatalf("json.Marshal() returned %q; expected %q", data, expected)
	}
}

func TestUnmarshalAction(t *testing.T) {
	data := `{"action":"addPlayer","id":1}`
	var wa powergrim.WrappedAction
	err := json.Unmarshal([]byte(data), &wa)
	if err != nil {
		t.Fatal(err)
	}
	expected := powergrim.AddPlayer{
		Action: "addPlayer",
		Id:     1,
	}
	if wa.Action != expected {
		t.Fatalf("json.Unmarshal() wrote %#v; expected %#v", wa.Action, expected)
	}
}
