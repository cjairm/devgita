package macos

import (
	"context"
	"fmt"
	"os"

	"github.com/cjairm/devgita/pkg/common"
	"github.com/manifoldco/promptui"
)

var languages = []string{
	"All",
	"None",
	"Done",
	"Node.js",
	"Go",
	"PHP",
	"Python",
}

func ChooseLanguages(ctx context.Context) context.Context {
	selectedLanguages, err := multiSelect("Select programming languages", languages)
	if err != nil {
		fmt.Println("\033[31mError: Error selecting languages.\033[0m")
		fmt.Println("Installation stopped.")
		os.Exit(1)
	}
	initialConfig := common.Config{
		SelectedLanguages: selectedLanguages,
	}
	fmt.Printf("\n")
	return common.WithConfig(ctx, initialConfig)
}

// Helper function to check if a slice contains a string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// Function to remove an item from a slice
func removeItem(slice []string, item string) []string {
	for i, s := range slice {
		if s == item {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}

func multiSelect(label string, options []string) ([]string, error) {
	var selectedLangs []string
	availableOptions := options

	for {
		// Remove "All" and "None" if all programming languages are selected
		if len(selectedLangs) == len(options)-3 { // Exclude "All", "None", "Done"
			availableOptions = removeItem(availableOptions, "All")
			availableOptions = removeItem(availableOptions, "None")
		}

		// Create a prompt for selecting languages
		prompt := promptui.Select{
			Label: label,
			Items: availableOptions,
		}

		// Show the prompt to the user
		_, result, err := prompt.Run()
		if err != nil {
			return nil, fmt.Errorf("prompt failed: %w", err)
		}

		// Handle special cases
		switch result {
		case "All":
			// Add all languages to selectedLangs
			selectedLangs = append(
				selectedLangs,
				availableOptions[3:]...)
			return selectedLangs, nil
		case "None":
			return []string{}, nil
		case "Done":
			return selectedLangs, nil
		default:
			// Add the selected language to the list of selected languages
			if !contains(selectedLangs, result) {
				selectedLangs = append(selectedLangs, result)
			}
			// Remove the selected item from the available options
			availableOptions = removeItem(availableOptions, result)
		}
	}
}
