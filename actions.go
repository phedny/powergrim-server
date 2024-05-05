package main

import (
	"encoding/json"
)

type AddPlayer struct {
	Action      string    `json:"action"`
	Id          int       `json:"id"`
	AfterPlayer int       `json:"afterPlayer,omitempty"`
	Character   string    `json:"character,omitempty"`
	Alignment   Alignment `json:"alignment,omitempty"`
}

type RemovePlayer struct {
	Action string `json:"action"`
	Id     int    `json:"id"`
}

type MovePlayer struct {
	Action      string `json:"action"`
	Id          int    `json:"id"`
	AfterPlayer int    `json:"afterPlayer"`
}

type UpdatePlayer struct {
	Action    string    `json:"action"`
	Id        int       `json:"id"`
	Character string    `json:"character,omitempty"`
	Alignment Alignment `json:"alignment,omitempty"`
}

type AddReminder struct {
	Action    string           `json:"action"`
	Character string           `json:"character"`
	Token     string           `json:"token"`
	Position  ReminderPosition `json:"position"`
}

type RemoveReminder struct {
	Action    string           `json:"action"`
	Character string           `json:"character"`
	Token     string           `json:"token"`
	Position  ReminderPosition `json:"position"`
}

type MoveReminder struct {
	Action       string           `json:"action"`
	Character    string           `json:"character"`
	Token        string           `json:"token"`
	FromPosition ReminderPosition `json:"fromPosition"`
	ToPosition   ReminderPosition `json:"toPosition"`
}

type WrappedAction struct {
	Action any
}

func (wa *WrappedAction) UnmarshalJSON(data []byte) error {
	var actionType struct {
		Action string `json:"action"`
	}
	err := json.Unmarshal(data, &actionType)
	if err != nil {
		return err
	}
	switch actionType.Action {
	case "addPlayer":
		var addPlayer AddPlayer
		err := json.Unmarshal(data, &addPlayer)
		wa.Action = addPlayer
		return err
	case "removePlayer":
		var removePlayer RemovePlayer
		err := json.Unmarshal(data, &removePlayer)
		wa.Action = removePlayer
		return err
	case "movePlayer":
		var movePlayer MovePlayer
		err := json.Unmarshal(data, &movePlayer)
		wa.Action = movePlayer
		return err
	case "updatePlayer":
		var updatePlayer UpdatePlayer
		err := json.Unmarshal(data, &updatePlayer)
		wa.Action = updatePlayer
		return err
	case "addReminder":
		var addReminder AddReminder
		err := json.Unmarshal(data, &addReminder)
		wa.Action = addReminder
		return err
	case "removeReminder":
		var removeReminder RemoveReminder
		err := json.Unmarshal(data, &removeReminder)
		wa.Action = removeReminder
		return err
	case "moveReminder":
		var moveReminder MoveReminder
		err := json.Unmarshal(data, &moveReminder)
		wa.Action = moveReminder
		return err
	default:
		return ErrInvalidAction
	}
}
