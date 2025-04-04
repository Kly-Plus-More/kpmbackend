package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/mailjet/mailjet-apiv3-go"
)

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	router := gin.Default()

	// Add CORS middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"POST", "GET", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	/////////////////// POST ROUTES /////////////////////
	router.POST("/training", HandleFormationForm)

	router.GET("/api/", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{
			"message": "Welcome to the API",
		})
	})

	err = router.Run(":8080")
	if err != nil {
		panic(err)
	}
	
}

///////////////////////////////// MAILJET FUNCTION ////////////////////////////////////////////////////////////////////////////////
func sendFormationEmail(name, email, message, phoneNumber string) error {
	// Get Mailjet API keys from environment variables
	apiKey := os.Getenv("MAILJET_API_KEY")
	apiSecret := os.Getenv("MAILJET_API_SECRET")
	if apiKey == "" || apiSecret == "" {
		return fmt.Errorf("Mailjet API keys not set")
	}

	// Create the Mailjet client
	client := mailjet.NewMailjetClient(apiKey, apiSecret)

	// Define the email content
	emailData := &mailjet.InfoSendMail{
		FromEmail: "klyplusandmore@gmail.com", // Replace with your verified sender email
		FromName:  fmt.Sprintf("%s (via HOME ABOMO LAW FIRM WEBSITE)", name), // User's name + website info
		Subject:   "New Training Request",
		TextPart: fmt.Sprintf("Name: %s\nPhone: %s\nEmail: %s\nMessage: %s",
			name, phoneNumber, email, message),
		HTMLPart: fmt.Sprintf("<strong>Name:</strong> %s<br><strong>Phone:</strong> %s<br><strong>Email:</strong> %s<br><strong>Message:</strong> %s",
			name, phoneNumber, email, message),
		Recipients: []mailjet.Recipient{
			{Email: "benazizsangare2@gmail.com"}, // Replace with your client's email
			{Email: "kerly.fenn@gmail.com"}, // Additional recipient
		},
		Headers: map[string]string{
			"Reply-To": fmt.Sprintf("%s <%s>", name, email), // Set Reply-To to the user's email
		},
	}

	// Send the email
	res, err := client.SendMail(emailData)
	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	log.Printf("Email sent! Response: %+v\n", res)
	return nil
}

type FormationRequest struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	Message     string `json:"message"`
	PhoneNumber string `json:"phonenumber"`
}

// Handler for contact form submission
func HandleFormationForm(context *gin.Context) {
	var req FormationRequest

	// Bind JSON request to the struct
	if err := context.ShouldBindJSON(&req); err != nil {
		log.Println("Invalid input:", err)
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Send the email
	err := sendFormationEmail(req.Name, req.Email, req.Message, req.PhoneNumber)
	if err != nil {
		log.Println("Failed to send email:", err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send email"})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"message": "Training Email sent successfully",
	})
}
