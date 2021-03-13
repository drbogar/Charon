package main

import "github.com/diamondburned/arikawa/v2/discord"

//Player represents a playing user
type Player struct {
	Role   Role
	Member discord.Member
	Alive  bool
}

// GetNick returns the player nick, if it is presented. Otherwise it is return the username.
func (p *Player) GetNick() string {
	if p.Member.Nick == "" {
		return p.Member.User.Username
	}
	return p.Member.Nick
}

//Role is a role in the game
type Role struct {
	Name    string
	Ability Ability
}

//Ability is the ability of a Role
type Ability struct {
	ActiveTime ActiveTime
	Using      func()
}

//ActiveTime is an enum for the ability active time
type ActiveTime uint

const (
	// Day indicatctives that the ability active only during the day
	Day ActiveTime = iota
	// Night indicatctives that the ability active only at night
	Night
	// Both indicatctives that the ability allway active
	Both
)

func (a ActiveTime) String() string {
	return [...]string{"Day", "Night", "Both"}[a]
}
