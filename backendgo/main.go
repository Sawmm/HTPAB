package main

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sahilm/fuzzy"
)

type Performance struct {
	Date    string
	Name    string
	Label   string
	Time    string
	Floor   string
	Closing bool
	Year    int
}

type App struct {
	Performances []Performance
}

func NewApp(filePath string) (*App, error) {
	fmt.Printf("Attempting to load CSV from file: %s\n", filePath)
	performances, err := loadCSV(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to load CSV: %w", err)
	}
	fmt.Printf("Successfully loaded %d performances from CSV\n", len(performances))
	return &App{Performances: performances}, nil
}

func loadCSV(filePath string) ([]Performance, error) {
	fmt.Printf("Opening CSV file: %s\n", filePath)
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open CSV file: %w", err)
	}
	defer file.Close()

	fmt.Println("CSV file opened successfully, creating reader")
	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1 // Allow variable number of fields

	fmt.Println("Reading all records from CSV")
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV: %w", err)
	}

	fmt.Printf("Read %d records from CSV\n", len(records))

	var performances []Performance
	for i, record := range records {
		if i == 0 || len(record) < 7 { // Skip header and malformed records
			continue
		}

		year, _ := strconv.Atoi(record[6])
		closing, _ := strconv.ParseBool(record[5])

		perf := Performance{
			Date:    record[0],
			Name:    strings.TrimSpace(record[1]),
			Label:   record[2],
			Time:    record[3],
			Floor:   record[4],
			Closing: closing,
			Year:    year,
		}
		performances = append(performances, perf)
	}

	fmt.Printf("Parsed %d performances from CSV\n", len(performances))
	return performances, nil
}

func (app *App) findDJ(name string) ([]Performance, bool) {
	name = strings.ToLower(strings.TrimSpace(name))
	var matches []Performance
	for _, perf := range app.Performances {
		if strings.ToLower(strings.TrimSpace(perf.Name)) == name {
			matches = append(matches, perf)
		}
	}
	return matches, len(matches) > 0
}

func (app *App) suggestDJ(name string) string {
	names := make([]string, len(app.Performances))
	for i, perf := range app.Performances {
		names[i] = perf.Name
	}
	matches := fuzzy.Find(name, names)
	if len(matches) > 0 {
		return matches[0].Str
	}
	return ""
}

func main() {
	fmt.Println("Starting application...")

	app, err := NewApp("berghain_lineup.csv")
	if err != nil {
		fmt.Println("Error initializing application:", err)
		return
	}

	if app == nil {
		fmt.Println("App is nil after initialization")
		return
	}

	fmt.Printf("Loaded %d performances from CSV\n", len(app.Performances))

	if len(app.Performances) == 0 {
		fmt.Println("No performances loaded from CSV")
		// Maybe print the current working directory and list files
		pwd, _ := os.Getwd()
		fmt.Printf("Current working directory: %s\n", pwd)
		files, _ := os.ReadDir(".")
		fmt.Println("Files in current directory:")
		for _, file := range files {
			fmt.Println(file.Name())
		}
	} else {
		fmt.Println("First 5 entries in the CSV:")
		for i, perf := range app.Performances[:min(5, len(app.Performances))] {
			fmt.Printf("%d: %+v\n", i+1, perf)
		}
	}

	r := gin.Default()

	// Configure CORS
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:8080"} // Adjust this to match your frontend URL
	r.Use(cors.New(config))

	r.POST("/check-dj", func(c *gin.Context) {
		var json struct {
			DJName string `json:"djName"`
		}
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		fmt.Printf("Received request for DJ: %s\n", json.DJName)

		matches, found := app.findDJ(json.DJName)
		if found {
			sort.Slice(matches, func(i, j int) bool {
				return matches[i].Year > matches[j].Year
			})
			latestPerf := matches[0]
			c.JSON(http.StatusOK, gin.H{
				"message":      fmt.Sprintf("%s has played at Berghain %d time(s)!", latestPerf.Name, len(matches)),
				"details":      fmt.Sprintf("Last played on %s at %s on %s.", latestPerf.Date, latestPerf.Time, latestPerf.Floor),
				"performances": matches,
			})
		} else {
			suggestion := app.suggestDJ(json.DJName)
			if suggestion != "" {
				suggestedMatches, _ := app.findDJ(suggestion)
				sort.Slice(suggestedMatches, func(i, j int) bool {
					return suggestedMatches[i].Year > suggestedMatches[j].Year
				})
				latestPerf := suggestedMatches[0]
				c.JSON(http.StatusOK, gin.H{
					"message":      fmt.Sprintf("%s might have played at Berghain. Did you mean %s?", json.DJName, suggestion),
					"suggestion":   suggestion,
					"details":      fmt.Sprintf("%s has played %d time(s), last on %s at %s on %s.", suggestion, len(suggestedMatches), latestPerf.Date, latestPerf.Time, latestPerf.Floor),
					"performances": suggestedMatches,
				})
			} else {
				c.JSON(http.StatusOK, gin.H{
					"message": fmt.Sprintf("%s has not played at Berghain (according to our records).", json.DJName),
				})
			}
		}
	})

	fmt.Println("Starting server on :3001...")
	r.Run(":3001")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
