package main

import (
	"crypto/sha256"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"sync"

	"github.com/gin-gonic/gin"
)

type HashData struct {
	Input     string `json:"input"`
	AppName   string `json:"appName"`
	ChartPath string `json:"chartPath"`
	Hash      string `json:"hash"`
}

var (
	hashList []HashData
	mutex    sync.Mutex
)

func generateHash(input string) string {
	hash := sha256.Sum256([]byte(input))
	return fmt.Sprintf("%x", hash)
}

func generateHashHandler(c *gin.Context) {
	input := c.Query("input")
	appName := c.Query("appName")
	chartPath := c.Query("chartPath")
	if input == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input is required"})
		return
	}

	hashValue := generateHash(input)

	mutex.Lock()
	hashList = append(hashList, HashData{Input: input, AppName: appName, ChartPath: chartPath, Hash: hashValue})
	mutex.Unlock()

	c.JSON(http.StatusOK, gin.H{"hash": hashValue})

	// Call the function to run Helm install command with the newly generated hash value
	runHelmInstall(appName, chartPath, hashValue)
}

func getAllHashesHandler(c *gin.Context) {
	mutex.Lock()
	defer mutex.Unlock()

	c.JSON(http.StatusOK, hashList)
}

func runHelmInstall(appName, chartPath, hashValue string) {
	cmd := exec.Command("helm", "install", appName, chartPath, fmt.Sprintf("--set=hash=%s", hashValue))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error running Helm install command: %v\n", err)
	}
}

func main() {
	router := gin.Default()

	router.POST("/generate", generateHashHandler)
	router.GET("/hashes", getAllHashesHandler)

	router.Run(":8081")
}
