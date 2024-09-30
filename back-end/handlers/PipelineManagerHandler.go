package handlers

import (
    "fmt"
    "log"
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/streadway/amqp"
)

// QuestionRequest holds the structure for incoming question data
type QuestionRequest struct {
    QuestionID   int    `json:"question_id"`
    Answer       string `json:"answer"`
    DocumentID   string `json:"document_id"`
    DocumentPath string `json:"document_path"`
}

// sendToRabbitMQ sends a message to the RabbitMQ queue
func sendToRabbitMQ(message string) error {
    // Connect to RabbitMQ server
    conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
    if err != nil {
        return fmt.Errorf("Failed to connect to RabbitMQ: %v", err)
    }
    defer conn.Close()

    // Create a channel
    ch, err := conn.Channel()
    if err != nil {
        return fmt.Errorf("Failed to open a channel: %v", err)
    }
    defer ch.Close()

    // Declare a queue
    q, err := ch.QueueDeclare(
        "pipeline_queue", // name
        false,            // durable
        false,            // delete when unused
        false,            // exclusive
        false,            // no-wait
        nil,              // arguments
    )
    if err != nil {
        return fmt.Errorf("Failed to declare a queue: %v", err)
    }

    // Publish the message to the queue
    err = ch.Publish(
        "",     // exchange
        q.Name, // routing key
        false,  // mandatory
        false,  // immediate
        amqp.Publishing{
            ContentType: "text/plain",
            Body:        []byte(message),
        })
    if err != nil {
        return fmt.Errorf("Failed to publish a message: %v", err)
    }

    log.Printf(" [x] Sent %s\n", message)
    return nil
}

// PipelineManagerHandler processes each question and sends it to RabbitMQ
func PipelineManagerHandler(c *gin.Context) {
    startTime := time.Now() // Start time to measure request processing time

    // Log request metadata
    fmt.Printf("\n---- Incoming Request ----\n")
    fmt.Printf("Time: %v\n", startTime.Format("2006-01-02 15:04:05"))
    fmt.Printf("Method: %s\n", c.Request.Method)
    fmt.Printf("Path: %s\n", c.Request.URL.Path)
    fmt.Printf("Client IP: %s\n", c.ClientIP())

    var req QuestionRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        // Log error details
        fmt.Printf("Error: Failed to bind JSON - %v\n", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Log parsed question request details
    fmt.Printf("Parsed Question Request:\n")
    fmt.Printf("  QuestionID: %d\n", req.QuestionID)
    fmt.Printf("  Answer: %s\n", req.Answer)
    fmt.Printf("  DocumentID: %s\n", req.DocumentID)
    fmt.Printf("  DocumentPath: %s\n", req.DocumentPath)

    // Create message to send to RabbitMQ
    message := fmt.Sprintf("QuestionID: %d, Answer: %s, DocumentID: %s, DocumentPath: %s",
        req.QuestionID, req.Answer, req.DocumentID, req.DocumentPath)

    // Send message to RabbitMQ
    if err := sendToRabbitMQ(message); err != nil {
        log.Printf("Failed to send message to RabbitMQ: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send message to RabbitMQ"})
        return
    }

    // Simulate pipeline manager logic (replace with real logic)
    fmt.Println("Processing question in the pipeline manager...")

    // Simulate some processing time
    time.Sleep(1 * time.Second) // Simulate delay

    // Log successful processing
    fmt.Printf("Processed successfully - QuestionID: %d, Answer: %s\n", req.QuestionID, req.Answer)

    // Calculate total processing time
    elapsedTime := time.Since(startTime)
    fmt.Printf("Processing Time: %v\n", elapsedTime)

    // Log response before sending it to the client
    fmt.Printf("Sending response: Status: success, QuestionID: %d, Answer: %s, DocumentID: %s, DocumentPath: %s\n",
        req.QuestionID, req.Answer, req.DocumentID, req.DocumentPath)

    // Return a JSON response
    c.JSON(http.StatusOK, gin.H{
        "status":           "success",
        "processed_question": req.QuestionID,
        "answer":           req.Answer,
        "document_id":      req.DocumentID,
        "document_path":    req.DocumentPath,
    })

    // Finalize log for request completion
    fmt.Printf("---- End of Request ----\n\n")
}
