package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Document defines the structure of a document in the database
type Document struct {
    ID           primitive.ObjectID `bson:"_id,omitempty"`  // MongoDB ObjectID
    Owner        primitive.ObjectID `bson:"owner"`          // User ID (owner)
    Path         string             `bson:"path"`           // Dosyanın kaydedildiği path
    OriginalName string             `bson:"original_name"`  // Orijinal dosya adı
}
