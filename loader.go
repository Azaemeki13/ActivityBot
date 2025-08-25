package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"

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

func isRoleValid(roleID string, rolesByID map[string]RoleInfo, botMax int) (RoleInfo, bool) {
	role, ok := rolesByID[roleID]
	if !ok {
		return RoleInfo{}, false
	}
	if role.Managed || role.Position > botMax {
		return role, false
	}
	return role, true
}

func validClassLabels(cfg GuildConfig, rolesByID map[string]RoleInfo, botMax int) []string {
	result := []string{}
	order := cfg.ClassOrder
	if len(order) == 0 {
		for id := range cfg.Classes {
			order = append(order, id)
		}
	}
	// here actually declaring an inline function that sorts during my loop if i < j and sort slice accordingly
	sort.Slice(order, func(i, j int) bool {
		return cfg.Classes[order[i]].ClassName < cfg.Classes[order[j]].ClassName
	})
	for _, classID := range order {
		class := cfg.Classes[classID]
		if _, ok := isRoleValid(class.RoleID, rolesByID, botMax); ok {
			result = append(result, classID)
		}
	}
	return result
}

func validWeaponLabel(ClassID string, cfg GuildConfig, rolesByID map[string]RoleInfo, botMax int) []string {
	// 1 - sort
	// 2 - compare if it's valid
	class, ok := cfg.Classes[ClassID]
	if !ok {
		return nil
	}
	// 1) build quick lookup: weaponID -> Weapon, refaire des exos dessus
	byID := make(map[string]Weapon, len(class.Weapon))
	for _, w := range class.Weapon {
		byID[w.WeaponID] = w
	}
	result := []string{}
	order := cfg.WeaponOrder[ClassID]
	if len(order) == 0 {
		ids := make([]string, 0, len(class.Weapon))
		for _, w := range class.Weapon {
			ids = append(ids, w.WeaponID)
		}
		sort.Slice(ids, func(x, y int) bool {
			wx, wy := byID[ids[x]], byID[ids[y]]
			return strings.ToLower(wx.WeaponName) > strings.ToLower((wy.WeaponName))
		})
		order = ids
	}
	for _, wid := range order {
		w, ok := byID[wid]
		if !ok {
			continue
		}
		if _, ok := isRoleValid(w.WeaponDiscID, rolesByID, botMax); ok {
			result := append(result, wid)
		}
	}
	return result
}
