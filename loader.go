package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/bwmarrin/discordgo"
)

// perms 590416021220433

func loadConfig(config_path string) (GuildConfig, error) {
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
				return cfg, fmt.Errorf("invalid weapon in class '%s': Missing ID, Name or Role ID", label)
			}
		}
	}
	return cfg, nil
}

func fetchGuildID(session *discordgo.Session) (string, error) {
	if len(session.State.Guilds) == 0 {
		return "", fmt.Errorf("bot is not connected, please retry")
	}
	return session.State.Guilds[0].ID, nil
}

func serverConfig(discordSession *discordgo.Session, guildID string) (map[string]RoleInfo, error) {
	roles, err := discordSession.GuildRoles(guildID)
	if err != nil {
		return nil, err
	}
	rolesByID := make(map[string]RoleInfo, len(roles))
	for _, r := range roles {
		rolesByID[r.ID] = RoleInfo{ID: r.ID, Name: r.Name, Position: r.Position, Managed: r.Managed}
	}
	return rolesByID, nil
}

func botMaxPosition(session *discordgo.Session, guildID string, rolesByID map[string]RoleInfo) (int, error) {
	botID := session.State.User.ID
	isGuild, err := session.GuildMember(guildID, botID)
	if err != nil {
		return 0, fmt.Errorf("bot isn't a part of the server: %w", err)
	}
	maxPos := -1
	for _, rid := range isGuild.Roles {
		if r, ok := rolesByID[rid]; ok && r.Position > maxPos {
			maxPos = r.Position
		}
	}
	if maxPos < 0 {
		return 0, fmt.Errorf("bot has no roles in this guild")
	}
	return maxPos, nil
}

func validateRoles(serverRoles map[string]RoleInfo, cfg GuildConfig, botMax int) Diagnostics {
	var d Diagnostics

	for label, class := range cfg.Classes {
		if r, ok := serverRoles[class.RoleID]; !ok {
			d.Missing = append(d.Missing, DiagItem{Path: "Class " + label, RoleID: class.RoleID, Reason: "missing"})
		} else {
			if r.Managed {
				d.Managed = append(d.Managed, DiagItem{Path: "Class " + label, RoleID: r.ID, Reason: "managed"})
			}
			if r.Position >= botMax {
				d.TooHigh = append(d.TooHigh, DiagItem{Path: "Class " + label, RoleID: r.ID, Reason: "too-high"})
			}
		}
		for _, w := range class.Weapon {
			path := fmt.Sprintf("Weapon %s/%s", label, w.WeaponName)
			if r, ok := serverRoles[w.WeaponDiscID]; !ok {
				d.Missing = append(d.Missing, DiagItem{Path: path, RoleID: w.WeaponDiscID, Reason: "missing"})
				continue
			} else {
				if r.Managed {
					d.Managed = append(d.Managed, DiagItem{Path: path, RoleID: r.ID, Reason: "managed"})
				}
				if r.Position >= botMax {
					d.TooHigh = append(d.TooHigh, DiagItem{Path: path, RoleID: r.ID, Reason: "too-high"})
				}
			}
		}
	}
	return d
}
