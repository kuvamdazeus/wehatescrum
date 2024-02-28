package main

import (
	"context"
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
	ISO string
}

type SummaryContents map[string]map[string][]CommitStat

type SummaryOpts struct {
	force bool
	date time.Time
	duration time.Duration
}

type Summary struct {
	Contents SummaryContents
	GeneratedAt int64
}

func getRedisSummary(opts SummaryOpts) (Summary, error) {
	ctx := context.Background()
	
	// get initialized redis client
	redisClient, err := getRedisClient()
	if err != nil {
		return Summary{}, err
	}

	// fetch summary data from redis
	strSummary, err := redisClient.Get(ctx, getRedisSummaryKey(opts)).Result()
	if err != nil {
		return Summary{}, err
	}

	var summary Summary
	err = json.Unmarshal([]byte(strSummary), &summary)
	if err != nil {
		return Summary{}, err
	}
	
	return summary, nil
}

func generateSummary(opts SummaryOpts) (Summary, error) {
	summary, err := getRedisSummary(opts)
	if err == nil && !opts.force {
		fmt.Println("CACHE HIT")
		return summary, nil
	}

	fmt.Println(time.Now().Unix(), summary.GeneratedAt)
	fmt.Println("Generating summary...")

	files, err := os.ReadDir("./repos")
	if err != nil {
		fmt.Println("BRUH, not working...", err)
		return summary, fmt.Errorf("error occured %s", err)
	}

	summaryContents := make(SummaryContents)

	for _, file := range files {
		if !file.IsDir() {
			continue
		}

		repoPath := fmt.Sprintf("./repos/%s", file.Name())

		repo, err := git.PlainOpen(repoPath)
		if err != nil {
			fmt.Println("ERR: repos dir contains non-git dir:", repoPath)
			continue
		}

		worktree, err := repo.Worktree()
		if err != nil {
			fmt.Println("Worktree couldn't be resolved", repoPath)
			return summary, fmt.Errorf("error occured %s", err)
		}

		// refName := plumbing.NewReferenceFromStrings("refs/heads/dev_stable", "dev_stable")
		// worktree.Checkout(&git.CheckoutOptions{
		// 	Branch: refName.Name(),
		// })
		err = worktree.Pull(&git.PullOptions{
			RemoteName: "origin",
		})
		if err != nil {
			fmt.Println("Error pulling repo:", err)
		}

		authorStats := make(map[string][]CommitStat)

		commits, _ := repo.CommitObjects()
		commits.ForEach(func(commit *object.Commit) error {
			ll := opts.date.Add(-opts.duration).Unix()
			ul := opts.date.Unix()

			commitTime := commit.Author.When.Unix()
			
			if commitTime >= ll && commitTime <= ul {
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
					ISO: commit.Author.When.Format(getConstants().timeLayout),
				}

				authorStats[author] = append(authorStats[author], commitStat)
			}

			return nil
		})

		summaryContents[file.Name()] = authorStats
	}

	summary = Summary {
		GeneratedAt: time.Now().Unix(),
		Contents: summaryContents,
	}
	
	bSummaryJson, err := json.Marshal(summary)
	if err != nil {
		fmt.Println("Error occured while json encoding", err)
		return summary, fmt.Errorf("error occured %s", err)
	}

	ctx := context.Background()
	redisClient, err := getRedisClient()
	if err != nil {
		fmt.Println("Error occured while fetching redis client", err)
		return summary, fmt.Errorf("error occured %s", err)
	}

	redisKey := getRedisSummaryKey(opts)
	result, err := redisClient.Set(ctx, redisKey, string(bSummaryJson), time.Hour * 24).Result()
	if err != nil {
		fmt.Println("Error setting redis data", err)
	}

	fmt.Println("Redis SET command result: ", result)

	return summary, nil
}
