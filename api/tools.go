package api

import (
	"GoForPT/pkg/cfg"
	"GoForPT/pkg/tools"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"os"
	"path/filepath"
)

func UploadStaticFile(c *gin.Context) {
	// Retrieve the file from the request
	file, fileheader, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(400, gin.H{"error": "Failed to get file"})
		return
	}
	defer file.Close()

	// Read the file content
	fileContent, err := io.ReadAll(file)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to read file"})
		return
	}

	// Generate a salted hash of the file content
	salt := tools.GenerateToken()
	hash := md5.New()
	hash.Write([]byte(salt))
	hash.Write(fileContent)
	hashedFilename := hex.EncodeToString(hash.Sum(nil))
	hashedFilename = fmt.Sprintf("%s%s", hex.EncodeToString(hash.Sum(nil)), filepath.Ext(fileheader.Filename))
	// Create the static folder if it doesn't exist
	staticFolder := cfg.Cfg.Site.StaticPath
	if _, err := os.Stat(staticFolder); os.IsNotExist(err) {
		err = os.Mkdir(staticFolder, os.ModePerm)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to create static folder"})
			return
		}
	}

	// Save the file with the hashed filename
	filePath := filepath.Join(staticFolder, hashedFilename)
	out, err := os.Create(filePath)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to create file"})
		return
	}
	defer out.Close()

	_, err = out.Write(fileContent)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to save file"})
		return
	}

	c.JSON(200, gin.H{"message": "File uploaded successfully", "filename": hashedFilename})
}
