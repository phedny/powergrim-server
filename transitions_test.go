package main_test

import (
	"reflect"
	"testing"

	powergrim "github.com/phedny/powergrim-server"
)

func TestAddPlayer(t *testing.T) {
	game := powergrim.Game{
		Players: []powergrim.Player{
			{Id: 1},
		},
	}

	got, err := game.AddPlayer(powergrim.AddPlayer{Id: 2})
	if err != nil {
		t.Fatal(err)
	}
	expected := powergrim.Game{
		Players: []powergrim.Player{
			{Id: 1},
			{Id: 2, Alive: true, FirstNight: true},
		},
	}

	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("AddPlayer() returned %#v; expected %#v", got, expected)
	}
}

func TestAddPlayerAfterPlayer(t *testing.T) {
	game := powergrim.Game{
		Players: []powergrim.Player{
			{Id: 1},
			{Id: 2},
			{Id: 3},
		},
	}

	got, err := game.AddPlayer(powergrim.AddPlayer{Id: 4, AfterPlayer: 1})
	if err != nil {
		t.Fatal(err)
	}
	expected := powergrim.Game{
		Players: []powergrim.Player{
			{Id: 1},
			{Id: 4, Alive: true, FirstNight: true},
			{Id: 2},
			{Id: 3},
		},
	}

	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("AddPlayer() returned %#v; expected %#v", got, expected)
	}
}

func TestAddPlayerUniqueId(t *testing.T) {
	game := powergrim.Game{
		Players: []powergrim.Player{
			{Id: 1},
		},
	}

	got, err := game.AddPlayer(powergrim.AddPlayer{Id: 1})
	if err != powergrim.ErrUniqueId {
		t.Fatalf("AddPlayer() returned (%#v, %s); expected error %s", got, err, powergrim.ErrUniqueId)
	}
}

func TestAddPlayerInvalidAfterPlayer(t *testing.T) {
	game := powergrim.Game{
		Players: []powergrim.Player{
			{Id: 1},
		},
	}

	got, err := game.AddPlayer(powergrim.AddPlayer{Id: 2, AfterPlayer: 3})
	if err != powergrim.ErrOptionalAfterPlayer {
		t.Fatalf("AddPlayer() returned (%#v, %s); expected error %s", got, err, powergrim.ErrOptionalAfterPlayer)
	}
}

func TestRemovePlayer(t *testing.T) {
	game := powergrim.Game{
		Players: []powergrim.Player{
			{Id: 1},
			{Id: 2, Character: "Something"},
			{Id: 3},
		},
		Reminders: []powergrim.Reminder{
			{Character: "Something", Token: "Some Token", Position: powergrim.ReminderPosition{1}},
			{Character: "Other", Token: "Must Stay", Position: powergrim.ReminderPosition{1}},
			{Character: "Other", Token: "Some Token", Position: powergrim.ReminderPosition{2}},
			{Character: "Other", Token: "Shared Token", Position: powergrim.ReminderPosition{2, 3}},
		},
	}

	got, err := game.RemovePlayer(powergrim.RemovePlayer{Id: 2})
	if err != nil {
		t.Fatal(err)
	}

	expected := powergrim.Game{
		Players: []powergrim.Player{
			{Id: 1},
			{Id: 3},
		},
		Reminders: []powergrim.Reminder{
			{Character: "Other", Token: "Must Stay", Position: powergrim.ReminderPosition{1}},
		},
	}
	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("RemovePlayer() returned %#v; expected %#v", got, expected)
	}
}

func TestMovePlayer(t *testing.T) {
	game := powergrim.Game{
		Players: []powergrim.Player{
			{Id: 1},
			{Id: 2},
			{Id: 3},
			{Id: 4},
		},
		Reminders: []powergrim.Reminder{
			{Character: "Something", Token: "Some Token", Position: powergrim.ReminderPosition{2}},
			{Character: "Something", Token: "Other Token", Position: powergrim.ReminderPosition{4, 1}},
		},
	}

	got, err := game.MovePlayer(powergrim.MovePlayer{Id: 1, AfterPlayer: 3})
	if err != nil {
		t.Fatal(err)
	}

	expected := powergrim.Game{
		Players: []powergrim.Player{
			{Id: 2},
			{Id: 3},
			{Id: 1},
			{Id: 4},
		},
		Reminders: []powergrim.Reminder{
			{Character: "Something", Token: "Some Token", Position: powergrim.ReminderPosition{2}},
			{Character: "Something", Token: "Other Token", Position: powergrim.ReminderPosition{1, 4}},
		},
	}
	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("MovePlayer() returned %#v; expected %#v", got, expected)
	}
}

func TestMovePlayerBreakingSharedToken(t *testing.T) {
	game := powergrim.Game{
		Players: []powergrim.Player{
			{Id: 1},
			{Id: 2},
			{Id: 3},
			{Id: 4},
		},
		Reminders: []powergrim.Reminder{
			{Character: "Something", Token: "Some Token", Position: powergrim.ReminderPosition{2}},
			{Character: "Something", Token: "Other Token", Position: powergrim.ReminderPosition{3, 4}},
		},
	}

	got, err := game.MovePlayer(powergrim.MovePlayer{Id: 1, AfterPlayer: 3})
	if err != powergrim.ErrMovingWithSharedReminder {
		t.Fatalf("MovePlayer() returned (%#v, %s); expected error %s", got, err, powergrim.ErrMovingWithSharedReminder)
	}
}

func TestUpdatePlayer(t *testing.T) {
	game := powergrim.Game{
		Players: []powergrim.Player{
			{Id: 1},
			{Id: 2},
			{Id: 3},
		},
	}

	got, err := game.UpdatePlayer(powergrim.UpdatePlayer{Id: 1, Character: "Something", Alignment: "good"})
	if err != nil {
		t.Fatal(err)
	}

	expected := powergrim.Game{
		Players: []powergrim.Player{
			{Id: 1, Character: "Something", Alignment: "good", FirstNight: true},
			{Id: 2},
			{Id: 3},
		},
	}
	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("UpdatePlayer() returned %#v; expected %#v", got, expected)
	}
}

func TestUpdatePlayeCharacterChange(t *testing.T) {
	game := powergrim.Game{
		Players: []powergrim.Player{
			{Id: 1, Character: "Something"},
			{Id: 2},
			{Id: 3},
		},
		Reminders: []powergrim.Reminder{
			{Character: "Something", Token: "Some Token", Position: powergrim.ReminderPosition{2}},
			{Character: "Other Character", Token: "Other Token", Position: powergrim.ReminderPosition{2}},
		},
	}

	got, err := game.UpdatePlayer(powergrim.UpdatePlayer{Id: 1, Character: "New Character", Alignment: "good"})
	if err != nil {
		t.Fatal(err)
	}

	expected := powergrim.Game{
		Players: []powergrim.Player{
			{Id: 1, Character: "New Character", Alignment: "good", FirstNight: true},
			{Id: 2},
			{Id: 3},
		},
		Reminders: []powergrim.Reminder{
			{Character: "Other Character", Token: "Other Token", Position: powergrim.ReminderPosition{2}},
		},
	}
	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("UpdatePlayer() returned %#v; expected %#v", got, expected)
	}
}

func TestAddCentralReminder(t *testing.T) {
	game := powergrim.Game{
		Players: []powergrim.Player{
			{Id: 1},
		},
	}

	got, err := game.AddReminder(powergrim.AddReminder{Character: "Something", Token: "Some Token", Position: powergrim.ReminderPosition{}})
	if err != nil {
		t.Fatal(err)
	}

	expected := powergrim.Game{
		Players: []powergrim.Player{
			{Id: 1},
		},
		Reminders: []powergrim.Reminder{
			{Character: "Something", Token: "Some Token", Position: powergrim.ReminderPosition{}},
		},
	}
	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("AddReminder() returned %#v; expected %#v", got, expected)
	}
}

func TestAddPlayerReminder(t *testing.T) {
	game := powergrim.Game{
		Players: []powergrim.Player{
			{Id: 1},
		},
	}

	got, err := game.AddReminder(powergrim.AddReminder{Character: "Something", Token: "Some Token", Position: powergrim.ReminderPosition{1}})
	if err != nil {
		t.Fatal(err)
	}

	expected := powergrim.Game{
		Players: []powergrim.Player{
			{Id: 1},
		},
		Reminders: []powergrim.Reminder{
			{Character: "Something", Token: "Some Token", Position: powergrim.ReminderPosition{1}},
		},
	}
	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("AddReminder() returned %#v; expected %#v", got, expected)
	}
}

func TestAddReminderToNonExistingPlayer(t *testing.T) {
	game := powergrim.Game{
		Players: []powergrim.Player{
			{Id: 1},
		},
	}

	got, err := game.AddReminder(powergrim.AddReminder{Character: "Something", Token: "Some Token", Position: powergrim.ReminderPosition{2}})
	if err != powergrim.ErrReminderPosition {
		t.Fatalf("AddReminder() returned (%#v, %s); expected error %s", got, err, powergrim.ErrReminderPosition)
	}
}

func TestAddSharedPlayerReminder(t *testing.T) {
	game := powergrim.Game{
		Players: []powergrim.Player{
			{Id: 1},
			{Id: 2},
			{Id: 3},
		},
	}

	got, err := game.AddReminder(powergrim.AddReminder{Character: "Something", Token: "Some Token", Position: powergrim.ReminderPosition{2, 3}})
	if err != nil {
		t.Fatal(err)
	}

	expected := powergrim.Game{
		Players: []powergrim.Player{
			{Id: 1},
			{Id: 2},
			{Id: 3},
		},
		Reminders: []powergrim.Reminder{
			{Character: "Something", Token: "Some Token", Position: powergrim.ReminderPosition{2, 3}},
		},
	}
	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("AddReminder() returned %#v; expected %#v", got, expected)
	}
}

func TestAddSharedPlayerReminderSwitchedOrder(t *testing.T) {
	game := powergrim.Game{
		Players: []powergrim.Player{
			{Id: 1},
			{Id: 2},
			{Id: 3},
		},
	}

	got, err := game.AddReminder(powergrim.AddReminder{Character: "Something", Token: "Some Token", Position: powergrim.ReminderPosition{1, 3}})
	if err != nil {
		t.Fatal(err)
	}

	expected := powergrim.Game{
		Players: []powergrim.Player{
			{Id: 1},
			{Id: 2},
			{Id: 3},
		},
		Reminders: []powergrim.Reminder{
			{Character: "Something", Token: "Some Token", Position: powergrim.ReminderPosition{3, 1}},
		},
	}
	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("AddReminder() returned %#v; expected %#v", got, expected)
	}
}

func TestAddReminderToNonExistingSharedPlayer(t *testing.T) {
	game := powergrim.Game{
		Players: []powergrim.Player{
			{Id: 1},
			{Id: 2},
			{Id: 3},
		},
	}

	got, err := game.AddReminder(powergrim.AddReminder{Character: "Something", Token: "Some Token", Position: powergrim.ReminderPosition{2, 4}})
	if err != powergrim.ErrReminderPosition {
		t.Fatalf("AddReminder() returned (%#v, %s); expected error %s", got, err, powergrim.ErrReminderPosition)
	}
}

func TestAddReminderToNonNeighbouringSharedPlayer(t *testing.T) {
	game := powergrim.Game{
		Players: []powergrim.Player{
			{Id: 1},
			{Id: 2},
			{Id: 3},
			{Id: 4},
		},
	}

	got, err := game.AddReminder(powergrim.AddReminder{Character: "Something", Token: "Some Token", Position: powergrim.ReminderPosition{1, 3}})
	if err != powergrim.ErrReminderPosition {
		t.Fatalf("AddReminder() returned (%#v, %s); expected error %s", got, err, powergrim.ErrReminderPosition)
	}
}

func TestRemoveReminder(t *testing.T) {
	game := powergrim.Game{
		Players: []powergrim.Player{
			{Id: 1},
			{Id: 2},
			{Id: 3},
		},
		Reminders: []powergrim.Reminder{
			{Character: "Something", Token: "Some Token", Position: powergrim.ReminderPosition{2, 3}},
		},
	}

	got, err := game.RemoveReminder(powergrim.RemoveReminder{Character: "Something", Token: "Some Token", Position: powergrim.ReminderPosition{3, 2}})
	if err != nil {
		t.Fatal(err)
	}

	expected := powergrim.Game{
		Players: []powergrim.Player{
			{Id: 1},
			{Id: 2},
			{Id: 3},
		},
		Reminders: []powergrim.Reminder{},
	}
	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("RemoveReminder() returned %#v; expected %#v", got, expected)
	}
}

func TestRemoveNonExistingReminder(t *testing.T) {
	game := powergrim.Game{
		Players: []powergrim.Player{
			{Id: 1},
			{Id: 2},
			{Id: 3},
		},
		Reminders: []powergrim.Reminder{
			{Character: "Something", Token: "Some Token", Position: powergrim.ReminderPosition{2, 3}},
		},
	}

	got, err := game.RemoveReminder(powergrim.RemoveReminder{Character: "Something", Token: "Some Token", Position: powergrim.ReminderPosition{1}})
	if err != powergrim.ErrExistingReminder {
		t.Fatalf("RemoveReminder() returned (%#v, %s); expected error %s", got, err, powergrim.ErrExistingReminder)
	}
}

func TestMoveReminder(t *testing.T) {
	game := powergrim.Game{
		Players: []powergrim.Player{
			{Id: 1},
			{Id: 2},
			{Id: 3},
		},
		Reminders: []powergrim.Reminder{
			{Character: "Something", Token: "Some Token", Position: powergrim.ReminderPosition{2, 3}},
		},
	}

	got, err := game.MoveReminder(powergrim.MoveReminder{Character: "Something", Token: "Some Token", FromPosition: powergrim.ReminderPosition{2, 3}, ToPosition: powergrim.ReminderPosition{1}})
	if err != nil {
		t.Fatal(err)
	}

	expected := powergrim.Game{
		Players: []powergrim.Player{
			{Id: 1},
			{Id: 2},
			{Id: 3},
		},
		Reminders: []powergrim.Reminder{
			{Character: "Something", Token: "Some Token", Position: powergrim.ReminderPosition{1}},
		},
	}
	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("MoveReminder() returned %#v; expected %#v", got, expected)
	}
}

func TestMoveNonExistingReminder(t *testing.T) {
	game := powergrim.Game{
		Players: []powergrim.Player{
			{Id: 1},
			{Id: 2},
			{Id: 3},
		},
		Reminders: []powergrim.Reminder{
			{Character: "Something", Token: "Some Token", Position: powergrim.ReminderPosition{2, 3}},
		},
	}

	got, err := game.MoveReminder(powergrim.MoveReminder{Character: "Something", Token: "Some Token", FromPosition: powergrim.ReminderPosition{1}, ToPosition: powergrim.ReminderPosition{}})
	if err != powergrim.ErrExistingReminder {
		t.Fatalf("MoveReminder() returned (%#v, %s); expected error %s", got, err, powergrim.ErrExistingReminder)
	}
}

func TestMovePReminderToNonExistingSharedPlayer(t *testing.T) {
	game := powergrim.Game{
		Players: []powergrim.Player{
			{Id: 1},
			{Id: 2},
			{Id: 3},
		},
		Reminders: []powergrim.Reminder{
			{Character: "Something", Token: "Some Token", Position: powergrim.ReminderPosition{2, 3}},
		},
	}

	got, err := game.MoveReminder(powergrim.MoveReminder{Character: "Something", Token: "Some Token", FromPosition: powergrim.ReminderPosition{2, 3}, ToPosition: powergrim.ReminderPosition{2, 4}})
	if err != powergrim.ErrReminderPosition {
		t.Fatalf("MoveReminder() returned (%#v, %s); expected error %s", got, err, powergrim.ErrReminderPosition)
	}
}

func TestMoveReminderToNonNeighbouringSharedPlayer(t *testing.T) {
	game := powergrim.Game{
		Players: []powergrim.Player{
			{Id: 1},
			{Id: 2},
			{Id: 3},
			{Id: 4},
		},
		Reminders: []powergrim.Reminder{
			{Character: "Something", Token: "Some Token", Position: powergrim.ReminderPosition{2, 3}},
		},
	}

	got, err := game.MoveReminder(powergrim.MoveReminder{Character: "Something", Token: "Some Token", FromPosition: powergrim.ReminderPosition{2, 3}, ToPosition: powergrim.ReminderPosition{1, 3}})
	if err != powergrim.ErrReminderPosition {
		t.Fatalf("AddReminder() returned (%#v, %s); expected error %s", got, err, powergrim.ErrReminderPosition)
	}
}
