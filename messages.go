package main

import "fmt"

func printDiagnostics(diag Diagnostics, rolesByID map[string]RoleInfo) {
	fmt.Println("=== Validation: === ")
	fmt.Println("Total of : ", len(rolesByID), "roles.")
	if len(diag.Missing) > 0 {
		fmt.Println("Missing:")
		for _, item := range diag.Missing {
			fmt.Printf(" - %s (role %s)\n", item.Path, item.RoleID)
		}
	}
	if len(diag.TooHigh) > 0 {
		fmt.Println("Too high:")
		for _, item := range diag.TooHigh {
			fmt.Printf(" - %s (role %s)\n", item.Path, item.RoleID)
		}
	}
	if len(diag.Managed) > 0 {
		fmt.Println("Managed:")
		for _, item := range diag.Managed {
			fmt.Printf(" - %s (role %s)\n", item.Path, item.RoleID)
		}
	}
}
