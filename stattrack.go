package main

import (
	"fmt"
	"os/exec"
	"strings"
)

func main() {
	cmd := exec.Command("git", "commit", "-m", ":3")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error running git commit: %v\n", err)
		return
	}

	outputStr := string(output)
	fmt.Println(outputStr)
	lines := strings.Split(outputStr, "\n")

	if len(lines) > 1 {
		parts := strings.Split(lines[1], ",")
		fmt.Println(parts)

		filesChanged := ""
		insertions := ""
		deletions := ""

		if len(parts) >= 1 {
			filesChanged = strings.Split(parts[0], " ")[1]
		}
		if len(parts) >= 2 {
			secondPart := strings.TrimSpace(parts[1])
			if strings.Contains(secondPart, "insertion") {
				insertions = strings.Split(secondPart, " ")[0]
			} else if strings.Contains(secondPart, "deletion") {
				deletions = strings.Split(secondPart, " ")[0]
			}
		}
		if len(parts) >= 3 {
			deletions = strings.Split(parts[2], " ")[1]
		}

		fmt.Printf("Files changed: %s\n", filesChanged)
		if insertions != "" {
			fmt.Printf("Insertions: %s\n", insertions)
		}
		if deletions != "" {
			fmt.Printf("Deletions: %s\n", deletions)
		}
	} else {
		fmt.Println("No second line found in output.")
	}
}
