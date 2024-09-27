package handlers

import (
    "fmt"
    "io/ioutil"
    "net/http"
    "path/filepath"
    "github.com/gin-gonic/gin"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
    "your-backend-module/models"
    "context"
)

// DocumentContentHandler dosyanın içeriğini döner
func DocumentContentHandler(documentCollection *mongo.Collection) gin.HandlerFunc {
    return func(c *gin.Context) {
        documentID := c.Param("docID") // Parametreden docID alınıyor
        fmt.Println("Fetching document with ID:", documentID)

        // ObjectID'yi oluşturuyoruz
        docObjectID, err := primitive.ObjectIDFromHex(documentID)
        if err != nil {
            fmt.Println("Invalid document ID:", err)
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid document ID"})
            return
        }

        // Document'ı MongoDB'den bulma
        var document models.Document
        err = documentCollection.FindOne(context.TODO(), bson.M{"_id": docObjectID}).Decode(&document)
        if err != nil {
            fmt.Println("Document not found:", err)
            c.JSON(http.StatusNotFound, gin.H{"error": "Document not found"})
            return
        }

        // Dosya yolunu oluşturma
        documentPath := document.Path // Database'deki path'i alıyoruz
        fmt.Println("Document Path:", documentPath)

        // Dosyanın içeriğini okuma
        fileContent, err := ioutil.ReadFile(documentPath)
        if err != nil {
            fmt.Println("Error reading file:", err)
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not read file"})
            return
        }

        // Dosya türüne göre geri döndürme
        if filepath.Ext(documentPath) == ".csv" {
            c.Data(http.StatusOK, "text/csv", fileContent)
        } else if filepath.Ext(documentPath) == ".xls" || filepath.Ext(documentPath) == ".xlsx" {
            c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", fileContent)
        } else {
            c.JSON(http.StatusUnsupportedMediaType, gin.H{"error": "Unsupported file type"})
        }
    }
}
