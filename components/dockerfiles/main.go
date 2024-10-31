package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"
)

type DockerImage struct {
	Name    string
	Tag     string
	Path    string
	Latest  string
	Updated time.Time
}

type DockerHubResponse struct {
	Results []struct {
		Name        string    `json:"name"`
		LastUpdated time.Time `json:"last_updated"`
	} `json:"results"`
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: docker-base-checker <directory>")
		os.Exit(1)
	}

	rootDir := os.Args[1]
	images := findDockerfiles(rootDir)
	checkUpdates(images)
}

func findDockerfiles(root string) []DockerImage {
	var images []DockerImage
	fromRegex := regexp.MustCompile(`(?i)^FROM\s+([^/\s]+/[^/\s]+|[^/\s]+)\s*:?\s*([^\s]*)\s*`)

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() && (d.Name() == "node_modules" || d.Name() == ".git") {
			return filepath.SkipDir
		}

		if strings.ToLower(d.Name()) == "dockerfile" {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				line := strings.TrimSpace(scanner.Text())
				matches := fromRegex.FindStringSubmatch(line)
				if len(matches) >= 3 {
					name := matches[1]
					tag := matches[2]
					if tag == "" {
						tag = "latest"
					}

					// Handle official images
					if !strings.Contains(name, "/") {
						name = "library/" + name
					}

					images = append(images, DockerImage{
						Name: name,
						Tag:  tag,
						Path: path,
					})
					break // Only process the first FROM statement (ignore multi-stage builds for now)
				}
			}
		}
		return nil
	})

	if err != nil {
		fmt.Printf("Error walking directory: %v\n", err)
		os.Exit(1)
	}

	return images
}

func checkUpdates(images []DockerImage) {
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, 5) // Limit concurrent requests

	for i := range images {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			semaphore <- struct{}{}        // Acquire semaphore
			defer func() { <-semaphore }() // Release semaphore

			img := &images[i]
			latest, updated, err := getLatestVersion(img.Name)
			if err != nil {
				fmt.Printf("Error checking %s: %v\n", img.Name, err)
				return
			}

			img.Latest = latest
			img.Updated = updated

			// Print results
			if img.Tag != latest {
				fmt.Printf("\nUpdate available for %s\n", img.Path)
				fmt.Printf("Current: %s:%s\n", img.Name, img.Tag)
				fmt.Printf("Latest:  %s:%s (updated %s)\n", img.Name, latest, updated.Format("2006-01-02"))
			}
		}(i)
	}

	wg.Wait()
}

func getLatestVersion(imageName string) (string, time.Time, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	url := fmt.Sprintf("https://hub.docker.com/v2/repositories/%s/tags/", imageName)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", time.Time{}, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", time.Time{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", time.Time{}, fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	var response DockerHubResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", time.Time{}, err
	}

	if len(response.Results) == 0 {
		return "", time.Time{}, fmt.Errorf("no tags found")
	}

	// Find the latest non-RC/beta tag
	for _, result := range response.Results {
		tag := result.Name
		if !strings.Contains(tag, "rc") && !strings.Contains(tag, "beta") && !strings.Contains(tag, "alpha") {
			return tag, result.LastUpdated, nil
		}
	}

	// If no stable version found, return the first tag
	return response.Results[0].Name, response.Results[0].LastUpdated, nil
}
