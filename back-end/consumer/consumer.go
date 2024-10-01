package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

// Struct for sending data to Django API
type FilePaths struct {
	InputPath  string `json:"input_path"`
	OutputPath string `json:"output_path"`
}

func callDjangoAPI(inputPath, outputPath string) error {
	// Prepare the data to be sent to the Django API
	data := FilePaths{
		InputPath:  inputPath,
		OutputPath: outputPath,
	}

	// Convert data to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %v", err)
	}

	// Send POST request to Django API
	resp, err := http.Post("http://127.0.0.1:8000/app/fake-deletion/", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to call Django API: %v", err)
	}
	defer resp.Body.Close()

	// Check for a successful response
	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("error from Django API: %v, response: %s", resp.Status, body)
	}

	log.Printf("Successfully called Django API for input: %s and output: %s", inputPath, outputPath)
	return nil
}
// callDjangoAPI'yi imputation işlemi için güncelleyelim
func callDjangoAPIimp(inputPath, outputPath string, answer string) error {
	// Prepare the data to be sent to the Django API
	data := map[string]interface{}{
		"input_path":  inputPath,
		"output_path": outputPath,
		"answer":      answer,  // Veri doldurma işlemi için answer bilgisi
	}

	// Convert data to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %v", err)
	}

	// Send POST request to Django API for imputation
	resp, err := http.Post("http://127.0.0.1:8000/app/data-imputation/", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to call Django API: %v", err)
	}
	defer resp.Body.Close()

	// Check for a successful response
	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("error from Django API: %v, response: %s", resp.Status, body)
	}

	log.Printf("Successfully called Django API for input: %s and output: %s", inputPath, outputPath)
	return nil
}

func deleteAllTempFiles(logFilePath string) error {
	// Log dosyasını oku
	data, err := ioutil.ReadFile(logFilePath)
	if err != nil {
		return fmt.Errorf("failed to read log file: %v", err)
	}

	// Tüm temp dosya path'lerini bul ve sil
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		parts := strings.Split(line, ",")
		if len(parts) < 4 {
			continue
		}
		tempFilePath := strings.TrimSpace(parts[3]) // Temp dosya yolunu al

		// Dosyayı sil
		if err := os.Remove(tempFilePath); err != nil && !os.IsNotExist(err) {
			log.Printf("Failed to delete temp file: %s, error: %v", tempFilePath, err)
		} else {
			log.Printf("Deleted temp file: %s", tempFilePath)
		}
	}

	// Log dosyasını temizle
	err = ioutil.WriteFile(logFilePath, []byte(""), 0644)
	if err != nil {
		return fmt.Errorf("failed to clear log file: %v", err)
	}

	return nil
}

func main() {
	// RabbitMQ'ya bağlan
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	// Kanal oluştur
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	// Kuyruğa bağlan
	q, err := ch.QueueDeclare(
		"pipeline_queue", // Kuyruk adı
		false,            // Kalıcı
		false,            // Silinsin mi?
		false,            // Exclusive mi?
		false,            // No-wait
		nil,              // Ekstra argümanlar
	)
	failOnError(err, "Failed to declare a queue")

	// Kuyruğu dinle
	msgs, err := ch.Consume(
		q.Name, // Kuyruk adı
		"",     // Consumer adı
		true,   // Auto-ack
		false,  // Exclusive
		false,  // No-local
		false,  // No-wait
		nil,    // Args
	)
	failOnError(err, "Failed to register a consumer")

	// Çalıştırma dizininin bir üst dizinini al
	baseDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get current directory: %v", err)
	}
	// Üst dizine git
	parentDir := filepath.Dir(baseDir)

	// 'scripts/tempdata' klasörünü oluştur
	scriptsDir := filepath.Join(parentDir, "scripts")
	tempDataDir := filepath.Join(scriptsDir, "tempdata")

	// Klasör yoksa oluştur
	err = os.MkdirAll(tempDataDir, os.ModePerm)
	if err != nil {
		log.Fatalf("Failed to create tempdata directory: %v", err)
	}

	// Her documentID için tempPath saklayan map yapısı
	tempPaths := make(map[string]string)
	var mu sync.Mutex // Mutex kullanarak map yapısını eş zamanlı erişime karşı koruyacağız

	// Sonsuza kadar dinle
	forever := make(chan bool)

	// Mesajları dinle ve işleyelim
	go func() {
		for d := range msgs {
			fmt.Printf("Received a message: %s\n", d.Body)

			// Mesajı split ile ayıralım
			messageParts := strings.Split(string(d.Body), ", ")
			if len(messageParts) != 4 {
				log.Printf("Invalid message format: %s", d.Body)
				continue
			}

			// Mesaj parçalarını ayrıştıralım
			questionID := strings.Split(messageParts[0], ": ")[1]
			answer := strings.Split(messageParts[1], ": ")[1]
			documentID := strings.Split(messageParts[2], ": ")[1]
			documentPath := strings.Split(messageParts[3], ": ")[1]

			// Yolu temizleyelim ve baseDir'in bir üstüne ekleyelim
			fmt.Printf("Input path: %s\n", documentPath)

			cleanDocumentPath := filepath.Clean(documentPath)
			fmt.Printf("Clean document path: %s\n", cleanDocumentPath)

			fullDocumentPath := filepath.Join(parentDir, cleanDocumentPath)
			fmt.Printf("Full document path: %s\n", fullDocumentPath)

			// Kontrol et: Bu bir dosya mı?
			info, err := os.Stat(fullDocumentPath)
			if os.IsNotExist(err) {
				log.Printf("Error: The file does not exist at path: %s", fullDocumentPath)
				continue
			}
			if info.IsDir() {
				log.Printf("Error: The path is a directory, not a file: %s", fullDocumentPath)
				continue
			}
			logFilePath := filepath.Join(parentDir, "logs", fmt.Sprintf("rabbitlog%s.txt", documentID))
			err = os.MkdirAll(filepath.Dir(logFilePath), os.ModePerm)
			if err != nil {
				log.Fatalf("Failed to create logs directory: %v", err)
			}
			if questionID == "1" {
				tempPaths[documentID]=fullDocumentPath
				log.Printf("First question received, deleting all temp files for DocumentID: %s", documentID)
				err := deleteAllTempFiles(logFilePath)
				if err != nil {
					log.Printf("Error deleting temp files: %v", err)
					continue
				}
			}
			// Log dosyasının yolu


			// Eğer ilk soruysa, tüm eski temp dosyaları sil


			// Mutex kullanarak map'e güvenli erişim
			mu.Lock()
			tempFilePath, exists := tempPaths[documentID]
			mu.Unlock()

			var newTempFilePath string

			// Eğer tempPath yoksa, yani bu ilk soruysa, yeni bir temp dosya oluştur
			if !exists {
				fileContent, err := ioutil.ReadFile(fullDocumentPath)
				if err != nil {
					log.Printf("Error reading file %s: %v", fullDocumentPath, err)
					continue
				}

				// İlk sorunun temp dosyasını oluştur
				originalFileName := filepath.Base(documentPath)
				ext := filepath.Ext(originalFileName) // .csv uzantısını al
				baseFileName := originalFileName[:len(originalFileName)-len(ext)] // uzantısız dosya adı

				newTempFilePath = fmt.Sprintf("%s_{%s_%s}%s", baseFileName, questionID, answer, ext) // uzantıyı en sona ekle
				newTempFilePath = filepath.Join(tempDataDir, newTempFilePath) // scripts/tempdata içine kaydet
				tempPaths[documentID]=newTempFilePath

				err = ioutil.WriteFile(newTempFilePath, fileContent, 0644)
				if err != nil {
					log.Printf("Error creating temp file %s: %v", newTempFilePath, err)
					continue
				}

				// tempPath'i map'e kaydet
				mu.Lock()
				tempPaths[documentID] = newTempFilePath
				mu.Unlock()
			} else {
				// Temp dosya zaten varsa, ona ekleme yap
				fileContent, err := ioutil.ReadFile(tempFilePath)
				if err != nil {
					log.Printf("Error reading temp file %s: %v", tempFilePath, err)
					continue
				}

				// Mevcut temp dosyasına sorunun cevabını ekle
				originalFileName := filepath.Base(tempFilePath)
				ext := filepath.Ext(originalFileName) // .csv uzantısını al
				baseFileName := originalFileName[:len(originalFileName)-len(ext)] // uzantısız dosya adı

				newTempFilePath = fmt.Sprintf("%s_{%s_%s}%s", baseFileName, questionID, answer, ext) // uzantıyı en sona ekle
				newTempFilePath = filepath.Join(tempDataDir, newTempFilePath)

				// Eğer ikinci sorunun cevabı "Option 1" ise Django API'yi çağır
				if questionID == "2" && answer == "Yes" {
					log.Printf("Calling Django API for second question Option 1")

					// Call Django API to process the file
					err = callDjangoAPI(tempFilePath, newTempFilePath)
					if err != nil {
						log.Printf("Error calling Django API: %v", err)
						continue
					}
				} else {
					// Option 2 veya Option 3 için işlemi yapmadan temp dosyayı kaydet
					err = ioutil.WriteFile(newTempFilePath, fileContent, 0644)
					if err != nil {
						log.Printf("Error creating updated temp file %s: %v", newTempFilePath, err)
						continue
					}
				}
// 3. soru geldiğinde Django API'yi imputation metoduyla çağır
if questionID == "3" {
    log.Printf("Third question received, calling Django API for imputation")

    // Call Django API for imputation
    err = callDjangoAPIimp(tempFilePath, newTempFilePath,answer)
    if err != nil {
        log.Printf("Error calling Django API for imputation: %v", err)
        continue
    }
}


				// Temp path'i güncelle
				mu.Lock()
				tempPaths[documentID] = newTempFilePath
				mu.Unlock()
			}

			// Log dosyasına yazarken güncellenmiş dosya adını kullan
			logEntry := fmt.Sprintf("%s,%s,%s,%s\n", questionID, answer, documentID, newTempFilePath) // Güncellenmiş temp dosya yolu
			f, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				log.Printf("Error opening log file %s: %v", logFilePath, err)
				continue
			}
			defer f.Close()
			if _, err := f.WriteString(logEntry); err != nil {
				log.Printf("Error writing to log file %s: %v", logFilePath, err)
				continue
			}

			// Log dosyasına başarıyla yazıldığını belirtelim
			log.Printf("Successfully processed and logged for DocumentID: %s\n", documentID)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
