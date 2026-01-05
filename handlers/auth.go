package handlers

import (
	"go-api/models"
	"go-api/repo"
	"go-api/utils"
	"net/http"
	"os"
	"regexp"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	SessionRepo *repo.RepoFactory[models.Session]
}

var customerRepo = NewCustomerHandler()

func NewSessionHandler() *AuthHandler {
	return &AuthHandler{
		SessionRepo: repo.NewRepoFactory(
			repo.NewPostgresRepo[models.Session]("sessions"),
			repo.NewMongoRepo[models.Session]("sessions"),
		),
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var body struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6,max=20"`
		Phone    string `json:"phone" binding:"required"`
		Name     string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(body.Password) < 6 || len(body.Password) > 20 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password must be between 6 and 20 characters"})
		return
	}

	if !regexp.MustCompile(`\d`).MatchString(body.Password) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password must contain at least one number"})
		return
	}

	if !regexp.MustCompile(`[A-Z]`).MatchString(body.Password) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password must contain at least one uppercase letter"})
		return
	}

	customer, err := models.NewCustomer(body.Name, body.Email, body.Phone, body.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create customer"})
		return
	}

	if err := customerRepo.CustomerRepo.Create(customer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create customer"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Customer created successfully",
		"data":    customer,
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var body struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6,max=20"`
	}
	var domain string

	if os.Getenv("GIN_MODE") != "production" {
		domain = "localhost"
	} else {
		domain = os.Getenv("DOMAIN_URL")
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	customer, err := customerRepo.CustomerRepo.GetByEmail(body.Email, repo.ReadSource("sql"))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid email or password"})
		return
	}

	if !utils.VerifyPassword(customer.HashedPassword, body.Password) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email or password"})
		return
	}

	token, err := utils.GenerateToken(customer.ID, customer.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	public_customer := utils.CustomerPublic{
		Name:      customer.Name,
		Email:     customer.Email,
		Phone:     customer.Phone,
		CreatedAt: customer.CreatedAt,
		UpdatedAt: customer.UpdatedAt,
	}
	session := models.NewSession(customer.ID, customer.Email, token)
	h.SessionRepo.Create(session)

	c.SetCookie("token", token, 3600, "/", domain, false, true)
	c.JSON(http.StatusOK, gin.H{
		"user": public_customer,
	})
}

func (h *AuthHandler) GetProfile(c *gin.Context) {
	// Get user ID from JWT middleware context
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Get customer
	customer, err := customerRepo.CustomerRepo.GetOne(userID, repo.ReadSource("sql"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
		return
	}

	public_customer := utils.CustomerPublic{
		Name:      customer.Name,
		Email:     customer.Email,
		Phone:     customer.Phone,
		CreatedAt: customer.CreatedAt,
		UpdatedAt: customer.UpdatedAt,
	}

	c.JSON(http.StatusOK, gin.H{
		"user": public_customer,
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	session, err := h.SessionRepo.GetByEmail(c.GetString("email"), repo.ReadSource("sql"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	h.SessionRepo.HardDelete(session.ID)

	c.SetCookie(
		"token",
		"",
		-1,
		"/",
		"",
		false,
		true,
	)

	c.JSON(http.StatusOK, gin.H{"message": "Log out successful"})
}
