package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Get the directory of the executable
	execPath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	execDir := filepath.Dir(execPath)

	// Initialize the SQLite database in the same directory as the executable
	databasePath := filepath.Join(execDir, "git_stats.db")
	database, err := sql.Open("sqlite3", databasePath)
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	// Create the table if it doesn't exist
	createTableSQL := `CREATE TABLE IF NOT EXISTS git_stats (
		"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		"files_changed" TEXT,
		"insertions" TEXT,
		"deletions" TEXT
	);`
	_, err = database.Exec(createTableSQL)
	if err != nil {
		log.Fatal(err)
	}

	// Check if a command-line argument is passed
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "show":
			showCommits(database)
			return
		case "push":
			fmt.Println("ToDo")
			return
		default:
			fmt.Println("Invalid command")
			return
		}
	}

	// Run the git commit command
	cmd := exec.Command("git", "commit", "-m", ":3")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error running git commit: %v\n", err)
		return
	}

	outputStr := string(output)
	fmt.Println(outputStr)
	lines := strings.Split(outputStr, "\n")

	filesChanged := ""
	insertions := ""
	deletions := ""

	if len(lines) > 1 {
		parts := strings.Split(lines[1], ",")
		fmt.Println(parts)

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
	}

	fmt.Printf("Files changed: %s\n", filesChanged)
	if insertions != "" {
		fmt.Printf("Insertions: %s\n", insertions)
	}
	if deletions != "" {
		fmt.Printf("Deletions: %s\n", deletions)
	}

	// Insert the data into the table
	insertSQL := `INSERT INTO git_stats (files_changed, insertions, deletions) VALUES (?, ?, ?)`
	statement, err := database.Prepare(insertSQL)
	if err != nil {
		log.Fatal(err)
	}
	defer statement.Close()

	_, err = statement.Exec(filesChanged, insertions, deletions)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Data successfully inserted into the database")
}

func showCommits(database *sql.DB) {
	rows, err := database.Query("SELECT id, files_changed, insertions, deletions FROM git_stats")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	fmt.Println("Commit Stats:")
	for rows.Next() {
		var id int
		var filesChanged, insertions, deletions string
		err := rows.Scan(&id, &filesChanged, &insertions, &deletions)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("ID: %d, Files Changed: %s, Insertions: %s, Deletions: %s\n", id, filesChanged, insertions, deletions)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
}
