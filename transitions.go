package main

import (
	"errors"
	"slices"
)

var (
	ErrInvalidAction            = errors.New("invalid action")
	ErrIdLength                 = errors.New("id must have length between 1 and 31")
	ErrUniqueId                 = errors.New("id must be unique")
	ErrExistingId               = errors.New("id must be id of existing player")
	ErrOptionalAfterPlayer      = errors.New("afterPlayer must be absent or id of existing player")
	ErrRequiredAfterPlayer      = errors.New("afterPlayer must be id of existing player")
	ErrDistinctIdAfterPlayer    = errors.New("id and afterPlayer must be distinct")
	ErrMovingWithSharedReminder = errors.New("moving a player must not disturb a shared reminder token")
	ErrExistingReminder         = errors.New("reminder must be present")
	ErrReminderPosition         = errors.New("position must be 0, player id, or array with 2 adjacent player ids")
)

func (game Game) ApplyAction(action any) (Game, error) {
	switch action := action.(type) {
	case AddPlayer:
		return game.AddPlayer(action)
	case RemovePlayer:
		return game.RemovePlayer(action)
	case MovePlayer:
		return game.MovePlayer(action)
	case UpdatePlayer:
		return game.UpdatePlayer(action)
	case AddReminder:
		return game.AddReminder(action)
	case RemoveReminder:
		return game.RemoveReminder(action)
	case MoveReminder:
		return game.MoveReminder(action)
	default:
		return Game{}, ErrInvalidAction
	}
}

func (game Game) AddPlayer(addPlayer AddPlayer) (Game, error) {
	if addPlayer.Id == 0 {
		return Game{}, ErrIdLength
	}
	for _, player := range game.Players {
		if player.Id == addPlayer.Id {
			return Game{}, ErrUniqueId
		}
	}
	player := Player{
		Id:         addPlayer.Id,
		Character:  addPlayer.Character,
		Alignment:  addPlayer.Alignment,
		Alive:      true,
		FirstNight: true,
	}
	if addPlayer.AfterPlayer == 0 {
		game.Players = append(slices.Clone(game.Players), player)
		return game, nil
	}
	afterPlayerIdx := slices.IndexFunc(game.Players, playerWithId(addPlayer.AfterPlayer))
	if afterPlayerIdx == -1 {
		return Game{}, ErrOptionalAfterPlayer
	}
	game.Players = slices.Insert(slices.Clone(game.Players), afterPlayerIdx+1, player)
	return game, nil
}

func (game Game) RemovePlayer(removePlayer RemovePlayer) (Game, error) {
	playerIdx := slices.IndexFunc(game.Players, playerWithId(removePlayer.Id))
	if playerIdx == -1 {
		return Game{}, ErrExistingId
	}
	player := game.Players[playerIdx]
	game.Players = append(slices.Clone(game.Players[:playerIdx]), game.Players[playerIdx+1:]...)
	reminders := make([]Reminder, 0, len(game.Reminders))
	for _, reminder := range game.Reminders {
		if player.Character != "" && reminder.Character == player.Character {
			continue
		}
		if reminder.Position[0] == removePlayer.Id || reminder.Position[1] == removePlayer.Id {
			continue
		}
		reminders = append(reminders, reminder)
	}
	game.Reminders = reminders
	return game, nil
}

func (game Game) MovePlayer(movePlayer MovePlayer) (Game, error) {
	if movePlayer.Id == movePlayer.AfterPlayer {
		return Game{}, ErrDistinctIdAfterPlayer
	}
	playerIdx := slices.IndexFunc(game.Players, playerWithId(movePlayer.Id))
	if playerIdx == -1 {
		return Game{}, ErrExistingId
	}
	player := game.Players[playerIdx]
	game.Players = append(slices.Clone(game.Players)[:playerIdx], game.Players[playerIdx+1:]...)
	afterPlayerIdx := slices.IndexFunc(game.Players, playerWithId(movePlayer.AfterPlayer))
	if afterPlayerIdx == -1 {
		return Game{}, ErrRequiredAfterPlayer
	}
	game.Players = slices.Insert(game.Players, afterPlayerIdx+1, player)
	var reminders []Reminder
	for reminderIdx, reminder := range game.Reminders {
		cPos, err := game.canonicalReminderPosition(reminder.Position)
		if err != nil {
			return Game{}, ErrMovingWithSharedReminder
		}
		if reminders == nil {
			reminders = slices.Clone(game.Reminders)
		}
		reminders[reminderIdx].Position = cPos
	}
	if reminders != nil {
		game.Reminders = reminders
	}
	return game, nil
}

func (game Game) UpdatePlayer(updatePlayer UpdatePlayer) (Game, error) {
	playerIdx := slices.IndexFunc(game.Players, playerWithId(updatePlayer.Id))
	if playerIdx == -1 {
		return Game{}, ErrExistingId
	}
	player := game.Players[playerIdx]
	if player.Character != updatePlayer.Character {
		if player.Character != "" {
			reminders := make([]Reminder, 0, len(game.Reminders))
			for _, reminder := range game.Reminders {
				if reminder.Character != player.Character {
					reminders = append(reminders, reminder)
				}
			}
			game.Reminders = reminders
		}
		player.Character = updatePlayer.Character
		player.FirstNight = true
	}
	player.Alignment = updatePlayer.Alignment
	game.Players[playerIdx] = player
	return game, nil
}

func (game Game) AddReminder(addReminder AddReminder) (Game, error) {
	cPos, err := game.canonicalReminderPosition(addReminder.Position)
	if err != nil {
		return Game{}, err
	}
	game.Reminders = append(slices.Clone(game.Reminders), Reminder{
		Character: addReminder.Character,
		Token:     addReminder.Token,
		Position:  cPos,
	})
	return game, nil
}

func (game Game) RemoveReminder(removeReminder RemoveReminder) (Game, error) {
	cPos, err := game.canonicalReminderPosition(removeReminder.Position)
	if err != nil {
		return Game{}, err
	}
	reminderIdx := slices.IndexFunc(game.Reminders, isReminder(Reminder{
		Character: removeReminder.Character,
		Token:     removeReminder.Token,
		Position:  cPos,
	}))
	if reminderIdx == -1 {
		return Game{}, ErrExistingReminder
	}
	game.Reminders = slices.Delete(slices.Clone(game.Reminders), reminderIdx, reminderIdx+1)
	return game, nil
}

func (game Game) MoveReminder(moveReminder MoveReminder) (Game, error) {
	cPos, err := game.canonicalReminderPosition(moveReminder.FromPosition)
	if err != nil {
		return Game{}, err
	}
	reminderIdx := slices.IndexFunc(game.Reminders, isReminder(Reminder{
		Character: moveReminder.Character,
		Token:     moveReminder.Token,
		Position:  cPos,
	}))
	if reminderIdx == -1 {
		return Game{}, ErrExistingReminder
	}
	cPos, err = game.canonicalReminderPosition(moveReminder.ToPosition)
	if err != nil {
		return Game{}, err
	}
	game.Reminders[reminderIdx].Position = cPos
	return game, nil
}

func (game Game) canonicalReminderPosition(position ReminderPosition) (ReminderPosition, error) {
	switch {
	case position[1] != 0:
		if position[0] == position[1] {
			return ReminderPosition{}, ErrReminderPosition
		}
		player1Idx := slices.IndexFunc(game.Players, playerWithId(position[0]))
		player2Idx := slices.IndexFunc(game.Players, playerWithId(position[1]))
		switch {
		case player1Idx == -1 || player2Idx == -1:
			return ReminderPosition{}, ErrReminderPosition
		case (player1Idx+1)%len(game.Players) == player2Idx:
			return position, nil
		case (player2Idx+1)%len(game.Players) == player1Idx:
			return ReminderPosition{position[1], position[0]}, nil
		default:
			return ReminderPosition{}, ErrReminderPosition
		}
	case position[0] != 0:
		for _, player := range game.Players {
			if player.Id == position[0] {
				return position, nil
			}
		}
		return ReminderPosition{}, ErrReminderPosition
	default:
		return ReminderPosition{}, nil
	}
}

func playerWithId(id int) func(Player) bool {
	return func(p Player) bool { return p.Id == id }
}

func isReminder(reminder Reminder) func(Reminder) bool {
	return func(reminder2 Reminder) bool {
		if reminder.Character != reminder2.Character || reminder.Token != reminder2.Token || len(reminder.Position) != len(reminder2.Position) {
			return false
		}
		for i, p := range reminder.Position {
			if p != reminder2.Position[i] {
				return false
			}
		}
		return true
	}
}
