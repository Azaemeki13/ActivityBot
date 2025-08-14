package main

// Weapon = a selectable subclass under a Class
type Weapon struct {
	WeaponID       string `json:"WeaponID"`                 // internal key, e.g. "hammer"
	WeaponName     string `json:"WeaponName"`               // display name, e.g. "Hammer"
	WeaponDiscID   string `json:"WeaponDiscID"`             // Discord Role snowflake for this weapon
	Specialisation int    `json:"Specialisation,omitempty"` // optional, defaults to 0
}

// Class = top-level category (Tank, Heal, etc.)
type Class struct {
	ClassID   string   `json:"ID"`      // internal key, e.g. "tank"
	ClassName string   `json:"Name"`    // display name, e.g. "TANK"
	RoleID    string   `json:"RoleID"`  // Discord Role snowflake for the class
	Weapon    []Weapon `json:"Weapons"` // weapons/subclasses under this class
}

// GuildConfig = everything the bot needs for one guild/server
type GuildConfig struct {
	Classes     map[string]Class    `json:"Classes"`               // key = class label ("Tank", "Heal", etc.) → Class
	ClassOrder  []string            `json:"ClassOrder,omitempty"`  // optional display order of class labels
	WeaponOrder map[string][]string `json:"WeaponOrder,omitempty"` // optional display order of weapon IDs per class label
}

// to pass role from discordgo.Role
type RoleInfo struct {
	ID       string
	Name     string
	Position int
	Managed  bool
}

type DiagItem struct {
	Path   string
	RoleID string
	Reason string
}

type Diagnostics struct {
	Missing []DiagItem
	TooHigh []DiagItem
	Managed []DiagItem
}
