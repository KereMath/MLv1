package handlers

import (
	"fmt"
	"net/http"

	"your-backend-module/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// GetDocumentsHandler retrieves documents from the MongoDB collection based on the user ID
func GetDocumentsHandler(documentCollection *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("GetDocumentsHandler started")

		// JWT Token'dan kullanıcı ID'sini alıyoruz
		tokenString := c.GetHeader("Authorization")
		fmt.Println("Authorization header received:", tokenString)

		if len(tokenString) < 7 || tokenString[:7] != "Bearer " {
			fmt.Println("Authorization token missing or incorrect format")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token is required"})
			return
		}

		// Remove "Bearer " prefix
		tokenString = tokenString[7:]
		fmt.Println("Token without Bearer prefix:", tokenString)

		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			fmt.Println("Parsing token")
			return jwtKey, nil
		})

		if err != nil {
			fmt.Println("Error parsing token:", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		if !token.Valid {
			fmt.Println("Token is invalid")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		// Token'dan userID'yi alıyoruz
		userIDString, ok := claims["user_id"].(string)
		if !ok {
			fmt.Println("Could not extract user ID from token claims")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not extract user ID from token"})
			return
		}

		userID, err := primitive.ObjectIDFromHex(userIDString)
		if err != nil {
			fmt.Println("Invalid user ID format")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
			return
		}

		// Kullanıcının sahip olduğu dökümanları MongoDB'den sorguluyoruz
		filter := bson.M{"owner": userID}
		findOptions := options.Find()

		cursor, err := documentCollection.Find(c, filter, findOptions)
		if err != nil {
			fmt.Println("Error finding documents for user:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve documents"})
			return
		}

		var documents []models.Document
		if err = cursor.All(c, &documents); err != nil {
			fmt.Println("Error decoding documents:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error decoding documents"})
			return
		}

		// Kullanıcıya dökümanları JSON formatında döndürüyoruz
		fmt.Println("Returning documents as JSON")
		c.JSON(http.StatusOK, gin.H{
			"documents": documents,
		})
	}
}
