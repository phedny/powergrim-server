package main_test

import (
	"encoding/json"
	"testing"

	powergrim "github.com/phedny/powergrim-server"
)

func TestMarshalCentralReminder(t *testing.T) {
	reminder := powergrim.Reminder{
		Character: "Drunk",
		Token:     "Is the Drunk",
		Position:  powergrim.ReminderPosition{},
	}

	got, err := json.Marshal(reminder)
	if err != nil {
		t.Fatal(err)
	}

	expected := `{"character":"Drunk","token":"Is the Drunk","position":"central"}`
	if string(got) != expected {
		t.Fatalf("json.Marshal() returned %q; expected %q", got, expected)
	}
}

func TestUnmarshalCentralReminder(t *testing.T) {
	var got powergrim.Reminder
	err := json.Unmarshal([]byte(`{"character":"Drunk","token":"Is the Drunk","position":"central"}`), &got)
	if err != nil {
		t.Fatal(err)
	}

	expected := powergrim.Reminder{
		Character: "Drunk",
		Token:     "Is the Drunk",
		Position:  powergrim.ReminderPosition{},
	}
	if got != expected {
		t.Fatalf("json.Unmarshal() returned %#v; expected %#v", got, expected)
	}
}

func TestMarshalPlayerReminder(t *testing.T) {
	reminder := powergrim.Reminder{
		Character: "Drunk",
		Token:     "Is the Drunk",
		Position:  powergrim.ReminderPosition{2},
	}

	got, err := json.Marshal(reminder)
	if err != nil {
		t.Fatal(err)
	}

	expected := `{"character":"Drunk","token":"Is the Drunk","position":2}`
	if string(got) != expected {
		t.Fatalf("json.Marshal() returned %q; expected %q", got, expected)
	}
}

func TestUnmarshalPlayerReminder(t *testing.T) {
	var got powergrim.Reminder
	err := json.Unmarshal([]byte(`{"character":"Drunk","token":"Is the Drunk","position":2}`), &got)
	if err != nil {
		t.Fatal(err)
	}

	expected := powergrim.Reminder{
		Character: "Drunk",
		Token:     "Is the Drunk",
		Position:  powergrim.ReminderPosition{2},
	}
	if got != expected {
		t.Fatalf("json.Unmarshal() returned %#v; expected %#v", got, expected)
	}
}

func TestMarshalSharedReminder(t *testing.T) {
	reminder := powergrim.Reminder{
		Character: "Revolutionary",
		Token:     "Register falsely?",
		Position:  powergrim.ReminderPosition{2, 3},
	}

	got, err := json.Marshal(reminder)
	if err != nil {
		t.Fatal(err)
	}

	expected := `{"character":"Revolutionary","token":"Register falsely?","position":[2,3]}`
	if string(got) != expected {
		t.Fatalf("json.Marshal() returned %q; expected %q", got, expected)
	}
}

func TestUnmarshalSharedReminder(t *testing.T) {
	var got powergrim.Reminder
	err := json.Unmarshal([]byte(`{"character":"Revolutionary","token":"Register falsely?","position":[2,3]}`), &got)
	if err != nil {
		t.Fatal(err)
	}

	expected := powergrim.Reminder{
		Character: "Revolutionary",
		Token:     "Register falsely?",
		Position:  powergrim.ReminderPosition{2, 3},
	}
	if got != expected {
		t.Fatalf("json.Unmarshal() returned %#v; expected %#v", got, expected)
	}
}
