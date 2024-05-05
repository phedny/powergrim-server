package main

import (
	"encoding/json"
	"errors"
)

var ErrInvalidAlignment = errors.New("invalid Alignment")
var ErrInvalidReminderPosition = errors.New("invalid ReminderPosition")

type Alignment string

func (a *Alignment) UnmarshalJSON(data []byte) error {
	switch string(data) {
	case "good":
		*a = "good"
		return nil
	case "evil":
		*a = "evil"
		return nil
	default:
		return ErrInvalidAlignment
	}
}

type Player struct {
	Id         int       `json:"id"`
	Position   [2]int    `json:"position"`
	Character  string    `json:"character,omitempty"`
	Alignment  Alignment `json:"alignment,omitempty"`
	Alive      bool      `json:"alive"`
	GhostVotes uint      `json:"ghostVotes,omitempty"`
	FirstNight bool      `json:"firstNight,omitempty"`
}

type ReminderPosition [2]int

func (rp ReminderPosition) MarshalJSON() ([]byte, error) {
	switch {
	case rp[1] != 0:
		return json.Marshal(rp[:])
	case rp[0] != 0:
		return json.Marshal(rp[0])
	default:
		return json.Marshal("central")
	}
}

func (rp *ReminderPosition) UnmarshalJSON(data []byte) error {
	if len(data) == 0 {
		return ErrInvalidReminderPosition
	}
	switch data[0] {
	case '"':
		if string(data) != `"central"` {
			return ErrInvalidReminderPosition
		}
		*rp = [2]int{0, 0}
		return nil
	case '[':
		slice := rp[:]
		return json.Unmarshal(data, &slice)
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		rp[1] = 0
		return json.Unmarshal(data, &rp[0])
	default:
		return ErrInvalidReminderPosition
	}
}

type Reminder struct {
	Character string           `json:"character"`
	Token     string           `json:"token"`
	Position  ReminderPosition `json:"position"`
}

type Game struct {
	Script    string     `json:"script"`
	Players   []Player   `json:"players"`
	Reminders []Reminder `json:"reminders"`
}
