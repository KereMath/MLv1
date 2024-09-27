package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"your-backend-module/models"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"

)

// UploadHandler handles file uploads and saves file metadata in MongoDB
func UploadHandler(documentCollection *mongo.Collection) gin.HandlerFunc {
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
		userIDString, ok := claims["user_id"].(string)
		if !ok {
			fmt.Println("Error extracting user ID from token claims")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not extract user ID from token"})
			return
		}

		userID, err := primitive.ObjectIDFromHex(userIDString) // ObjectID'ye çevirme
		if err != nil {
			fmt.Println("Invalid user ID format")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
			return
		}
		fmt.Println("User ID extracted from token:", userID)

		// Get the uploaded file
		file, err := c.FormFile("file")
		if err != nil {
			fmt.Println("No file uploaded or error retrieving file:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
			return
		}
		fmt.Println("File uploaded:", file.Filename)

		// Create a new document entry in MongoDB and get the document ID
		document := models.Document{
			Owner:        userID,
			OriginalName: file.Filename, // Orijinal dosya adını saklıyoruz
		}

		insertResult, err := documentCollection.InsertOne(c, document)
		if err != nil {
			fmt.Println("Error saving document to MongoDB:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not save document metadata"})
			return
		}

		documentID := insertResult.InsertedID.(primitive.ObjectID) // Belge ID'sini alıyoruz
		fmt.Println("Document ID generated:", documentID.Hex())

		// Create user's folder: data/{user_id}
		userFolder := filepath.Join("..", "data", userID.Hex())
		fmt.Println("User folder path:", userFolder)

		if _, err := os.Stat(userFolder); os.IsNotExist(err) {
			fmt.Println("User folder does not exist, creating folder:", userFolder)
			err = os.MkdirAll(userFolder, os.ModePerm)
			if err != nil {
				fmt.Println("Error creating user folder:", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create user folder"})
				return
			}
		} else {
			fmt.Println("User folder already exists:", userFolder)
		}

		// Save the file using the document ID as filename
		filePath := filepath.Join(userFolder, documentID.Hex()+filepath.Ext(file.Filename)) // Dosyayı documentID ile kaydediyoruz
		fmt.Println("Saving file to:", filePath)

		if err := c.SaveUploadedFile(file, filePath); err != nil {
			fmt.Println("Error saving file:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not save file"})
			return
		}

		// Update the document with the saved file path
		filter := bson.M{"_id": documentID}
		update := bson.M{"$set": bson.M{"path": filePath}}

		_, err = documentCollection.UpdateOne(c, filter, update)
		if err != nil {
			fmt.Println("Error updating document with file path:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update document with file path"})
			return
		}

		fmt.Println("File and document metadata saved successfully")
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("File %s uploaded successfully to %s", file.Filename, userFolder),
		})
	}
}
