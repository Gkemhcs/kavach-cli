package utils

import (
	"fmt"
	"os"

	"github.com/manifoldco/promptui"
)

// ConfirmSecretGroupCreation prompts the user to confirm a secret group creation or deletion action.
// Returns true if the user confirms, false otherwise.
func ConfirmSecretGroupCreation(message string) bool {
	printKavachLogo()
	prompt := promptui.Select{
		Label: message,
		Items: []string{"✅ Yes, proceed", "❌ No, cancel"},
		Size:  2,
		Templates: &promptui.SelectTemplates{
			Active:   "▶ {{ . | cyan }}",
			Inactive: "  {{ . | faint }}",
			Selected: "✔ {{ . | green }}",
		},
	}
	i, _, err := prompt.Run()
	if err != nil {
		fmt.Println("❌ Prompt cancelled.")
		os.Exit(1)
	}
	return i == 0 // 0 means yes
}

// printKavachLogo prints the Kavach CLI logo to the terminal.
func printKavachLogo() {
	fmt.Println(`🛡️  -------------------------------
🛡️   KAVACH CLI - Secret Syncing
🛡️  -------------------------------`)
}
