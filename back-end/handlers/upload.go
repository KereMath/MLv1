package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// UploadHandler handles file uploads
func UploadHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("UploadHandler started")

		// Get the JWT token from the Authorization header
		tokenString := c.GetHeader("Authorization")
		fmt.Println("Authorization header received:", tokenString)

		if len(tokenString) < 7 || tokenString[:7] != "Bearer " {
			fmt.Println("Authorization token missing or incorrect format")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token is required"})
			return
		}

		// Remove the "Bearer " prefix
		tokenString = tokenString[7:]
		fmt.Println("Token without Bearer prefix:", tokenString)

		// Parse the token and extract claims
		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			fmt.Println("Parsing token")
			return jwtKey, nil
		})

		if err != nil {
			fmt.Println("Error parsing token:", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}

		if !token.Valid {
			fmt.Println("Invalid token")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		// Get the user ID from the claims
		userID, ok := claims["user_id"].(string)
		if !ok {
			fmt.Println("Error extracting user ID from token claims")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not extract user ID from token"})
			return
		}
		fmt.Println("User ID extracted from token:", userID)

		// Create user's folder: data/{user_id}
		userFolder := filepath.Join("..", "data", userID)
		fmt.Println("User folder path:", userFolder)

		if _, err := os.Stat(userFolder); os.IsNotExist(err) {
			fmt.Println("User folder does not exist, creating folder:", userFolder)
			err = os.MkdirAll(userFolder, os.ModePerm) // os.MkdirAll ensures all intermediate directories are created
			if err != nil {
				fmt.Println("Error creating user folder:", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create user folder"})
				return
			}
		} else {
			fmt.Println("User folder already exists:", userFolder)
		}

		// Get the uploaded file
		file, err := c.FormFile("file")
		if err != nil {
			fmt.Println("No file uploaded or error retrieving file:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
			return
		}
		fmt.Println("File uploaded:", file.Filename)

		// Save the file to the user's folder
		filePath := filepath.Join(userFolder, file.Filename)
		fmt.Println("Saving file to:", filePath)

		if err := c.SaveUploadedFile(file, filePath); err != nil {
			fmt.Println("Error saving file:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not save file"})
			return
		}

		fmt.Println("File saved successfully to", filePath)
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("File %s uploaded successfully to %s", file.Filename, userFolder),
		})
	}
}
