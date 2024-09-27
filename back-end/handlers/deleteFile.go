package handlers

import (
	"fmt"
	"net/http"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/gin-gonic/gin"
)

// DeleteFileHandler handles the deletion of a file and its metadata from MongoDB
func DeleteFileHandler(documentCollection *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the document ID from the URL
		documentIDString := c.Param("id")
		fmt.Println(documentIDString)
		documentID, err := primitive.ObjectIDFromHex(documentIDString)
		if err != nil {
			fmt.Println("Invalid document ID:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid document ID"})
			return
		}

		// Find the document in MongoDB
		var document struct {
			Path string `bson:"path"`
		}
		err = documentCollection.FindOne(c, bson.M{"_id": documentID}).Decode(&document)
		if err != nil {
			fmt.Println("Document not found:", err)
			c.JSON(http.StatusNotFound, gin.H{"error": "Document not found"})
			return
		}

		// Delete the file from the file system
		err = os.Remove(document.Path)
		if err != nil {
			fmt.Println("Error deleting file:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not delete file"})
			return
		}

		// Delete the document metadata from MongoDB
		_, err = documentCollection.DeleteOne(c, bson.M{"_id": documentID})
		if err != nil {
			fmt.Println("Error deleting document metadata:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not delete document metadata"})
			return
		}

		fmt.Println("File and document metadata deleted successfully")
		c.JSON(http.StatusOK, gin.H{"message": "File deleted successfully"})
	}
}
