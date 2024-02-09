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

type SummaryContents map[string]map[string][]CommitStat

type Summary struct {
	Contents SummaryContents
	GeneratedAt int64
}

func getSummary() (Summary, error) {
	summaryFile, err := os.ReadFile("summary.json")
  if err != nil {
		return Summary{
			GeneratedAt: 0,
		}, err
	}
  
  var summary Summary
  err = json.Unmarshal(summaryFile, &summary)
	if err != nil {
		return Summary{
			GeneratedAt: 0,
		}, err
	}

  return summary, nil
}

func generateSummary(force bool) {
	summary, _ := getSummary()
	fmt.Println(time.Now().Unix(), summary.GeneratedAt)
	if (time.Now().Unix() - summary.GeneratedAt <= 1 * 24 * 60 * 60) && !force {
		return
	}

	files, err := os.ReadDir("./tmp")
	if err != nil {
		fmt.Println("BRUH, not working...", err)
		return
	}

	summaryContents := make(SummaryContents)

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
			if time.Now().Unix()-commit.Author.When.Unix() <= 24 * 60 * 60 {
				stats, _ := commit.Stats()

				author := fmt.Sprintf("%s <%s>", commit.Author.Name, commit.Author.Email)
				commitMsg := commit.Message
				totalChanges := 0
				filesChanged := []string{}

				for _, stat := range strings.Split(strings.Trim(stats.String(), "\n"), "\n") {
					tmp := strings.Split(stat, " | ")
					if len(tmp) < 2 {
						continue
					}

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

		summaryContents[file.Name()] = authorStats
	}

	fileName := "summary.json"
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	summary = Summary {
		GeneratedAt: time.Now().Unix(),
		Contents: summaryContents,
	}

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
