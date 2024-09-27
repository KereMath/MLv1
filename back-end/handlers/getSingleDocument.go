package handlers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"context"
)

// GetSingleDocumentHandler dökümanı ID'ye göre döndürür
func GetSingleDocumentHandler(documentCollection *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		docID := c.Param("docID")
		objID, err := primitive.ObjectIDFromHex(docID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid document ID"})
			return
		}

		var document bson.M
		err = documentCollection.FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&document)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Document not found"})
			return
		}

		c.JSON(http.StatusOK, document)
	}
}
