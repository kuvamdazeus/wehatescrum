package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

type CommitStat struct {
	CommitMsg    string
	TotalChanges int
	Hash         string
	FilesChanged []string
}

type Summary map[string]map[string][]CommitStat

func generate_summary() {
	files, err := os.ReadDir("./tmp")
	if err != nil {
		fmt.Println("BRUH, not working...")
		return
	}

	summary := make(Summary)

	for _, file := range files {
		if !file.IsDir() {
			continue
		}

		repoPath := fmt.Sprintf("./tmp/%s", file.Name())

		repo, err := git.PlainOpen(repoPath)
		if err != nil {
			fmt.Println("ERR: tmp dir contains non-git dir:", repoPath)
			continue
		}

		worktree, err := repo.Worktree()
		if err != nil {
			fmt.Println("Worktree couldn't be resolved", repoPath)
			return
		}

		// refName := plumbing.NewReferenceFromStrings("refs/heads/dev_stable", "dev_stable")
		// worktree.Checkout(&git.CheckoutOptions{
		// 	Branch: refName.Name(),
		// })
		worktree.Pull(&git.PullOptions{
			RemoteName: "origin",
		})

		authorStats := make(map[string][]CommitStat)

		commits, _ := repo.CommitObjects()
		commits.ForEach(func(commit *object.Commit) error {
			if time.Now().Unix()-commit.Author.When.Unix() <= 24*60*60 {
				stats, _ := commit.Stats()

				author := fmt.Sprintf("%s <%s>", commit.Author.Name, commit.Author.Email)
				commitMsg := commit.Message
				totalChanges := 0
				filesChanged := []string{}

				for _, stat := range strings.Split(strings.Trim(stats.String(), "\n"), "\n") {
					tmp := strings.Split(stat, " | ")
					nChanges, _ := strconv.Atoi(strings.Split(tmp[1], " ")[0])

					filesChanged = append(filesChanged, strings.Trim(tmp[0], " "))
					totalChanges += nChanges
				}

				commitStat := CommitStat{
					CommitMsg:    strings.Trim(commitMsg, "\n"),
					TotalChanges: totalChanges,
					Hash:         commit.Hash.String(),
					FilesChanged: filesChanged,
				}

				authorStats[author] = append(authorStats[author], commitStat)
			}

			return nil
		})

		summary[file.Name()] = authorStats
	}

	fileName := "summary.json"
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")

	err = encoder.Encode(summary)
	if err != nil {
		fmt.Println("Error encountered while encoding to JSON:", err)
		return
	}

	fmt.Println("Summary written to file:", fileName)
}
