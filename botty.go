package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// perms 590416021220433
/*
guildee
123456ABCDEF
Tank-> 12345asfsd
*/

// Weapon = a selectable subclass under a Class
type Weapon struct {
	WeaponID       string `json:"ID"`                       // internal key, e.g. "hammer"
	WeaponName     string `json:"Name"`                     // display name, e.g. "Hammer"
	WeaponDiscID   string `json:"RoleID"`                   // Discord Role snowflake for this weapon
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

func load_config(config_path string) (GuildConfig, error) {
	var cfg GuildConfig

	data, err := os.ReadFile(config_path)
	if err != nil {
		return cfg, fmt.Errorf("error reading config file: %w", err)
	}
	err = json.Unmarshal(data, &cfg) // unpack data json in cfg
	if err != nil {
		return cfg, fmt.Errorf("error parsing config JSON: %w", err)
	}
	// json validation
	for label, class := range cfg.Classes {
		if class.ClassID == "" || class.ClassName == "" || class.RoleID == "" {
			return cfg, fmt.Errorf("invalid class '%s': missing ID, Name or RoleID", label)
		}
		for _, weapon := range class.Weapon {
			if weapon.WeaponID == "" || weapon.WeaponName == "" || weapon.WeaponDiscID == "" {
				return cfg, fmt.Errorf("invalid weapin in class '%s': Missing ID, Name or Role ID", label)
			}
		}
	}
	return cfg, nil
}

func main() {
	token := os.Getenv("ALBION_HELPER")
	config := "./config.json"
	if token == "" {
		fmt.Println(":warning: No bot token found in env.")
		return
	}
	cfg, err := load_config(config)
	if err != nil {
		fmt.Println("Error, failed to load config:", err)
		return
	}
	fmt.Println("Loaded config successfully !")
	fmt.Printf("Loaded %d classes from %s\n", len(cfg.Classes), config)
}
