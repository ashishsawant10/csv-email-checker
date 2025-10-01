package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sync"

	"github.com/google/uuid"
)

func main() {

	// Job struct represents a processing job with its ID, file path, and status
	type Job struct {
		ID       string
		FilePath string
		Done     bool
	}

	jobs := make(map[string]*Job)
	var mu sync.Mutex

	// Simple email regex for validation
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

	// Directory to store processed CSVs
	storageDir := "uploads"
	os.MkdirAll(storageDir, os.ModePerm)

	// Endpoint 1: /API/upload
	http.HandleFunc("/API/upload", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, `{"error":"only POST allowed"}`, http.StatusBadRequest)
			return
		}

		// Get uploaded file
		file, _, err := r.FormFile("file")
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"%v"}`, err), http.StatusBadRequest)
			return
		}
		defer file.Close()

		// Generate unique job ID
		id := uuid.New().String()
		outputPath := filepath.Join(storageDir, id+".csv")
		outFile, err := os.Create(outputPath)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"%v"}`, err), http.StatusInternalServerError)
			return
		}
		defer outFile.Close()

		reader := csv.NewReader(file)
		writer := csv.NewWriter(outFile)

		// Read header row and add "flag" column
		header, err := reader.Read()
		if err != nil {
			http.Error(w, `{"error":"invalid CSV"}`, http.StatusBadRequest)
			return
		}
		header = append(header, "flag")
		writer.Write(header)

		// Process each row
		for {
			record, err := reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				http.Error(w, `{"error":"failed to parse CSV"}`, http.StatusBadRequest)
				return
			}
			if len(record) == 0 {
				continue
			}

			// Default flag = false, set to true if email found
			flag := "false"
			for _, field := range record {
				if emailRegex.MatchString(field) {
					flag = "true"
					break
				}
			}
			record = append(record, flag)
			writer.Write(record)
		}
		writer.Flush()

		mu.Lock()
		jobs[id] = &Job{ID: id, FilePath: outputPath, Done: true}
		mu.Unlock()

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"id": id})
	})

	// Endpoint 2: /API/download/{id}
	http.HandleFunc("/API/download/", func(w http.ResponseWriter, r *http.Request) {
		id := filepath.Base(r.URL.Path)

		// Check if job exists
		mu.Lock()
		job, ok := jobs[id]
		mu.Unlock()
		if !ok {
			http.Error(w, `{"error":"invalid id"}`, http.StatusBadRequest)
			return
		}

		// If job not completed
		if !job.Done {
			http.Error(w, `{"error":"job in progress"}`, http.StatusLocked)
			return
		}

		// Serve processed file
		http.ServeFile(w, r, job.FilePath)
	})

	// Start server on port 8080
	fmt.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
