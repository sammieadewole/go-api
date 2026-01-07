package handlers

import (
	"go-api/middleware"
	"go-api/models"
	"go-api/repo"
	"net/http"
	"os"
	"regexp"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	SessionRepo *repo.RepoFactory[models.Session]
}

type AdminHandler struct {
	AdminRepo *repo.RepoFactory[models.AdminUser]
}

type CombinedAuthHandler struct {
	*AuthHandler
	*AdminHandler
}

func NewCombinedHandker() *CombinedAuthHandler {
	return &CombinedAuthHandler{
		AuthHandler: &AuthHandler{
			SessionRepo: repo.NewRepoFactory(
				repo.NewPostgresRepo[models.Session]("sessions"), // Remove .(repo.Repository[models.Session])
				repo.NewMongoRepo[models.Session]("sessions"),
			),
		},
		AdminHandler: &AdminHandler{
			AdminRepo: repo.NewRepoFactory(
				repo.NewPostgresRepo[models.AdminUser]("admin_users"), // Already correct
				repo.NewMongoRepo[models.AdminUser]("admin_users"),
			),
		},
	}
}

func (h *CombinedAuthHandler) Register(c *gin.Context) {
	var body struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
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

	user := models.NewAdminUser(body.Email, body.Password)
	storage := repo.GetReadSource(c)
	if err := h.AdminRepo.Create(user, storage); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create session"})
		return
	}
}

func (h *CombinedAuthHandler) Login(c *gin.Context) {
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

	storage := repo.GetReadSource(c)

	user, err := h.AdminHandler.AdminRepo.GetByEmail(body.Email, storage)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid email or password"})
		return
	}

	token, err := middleware.GenerateToken(user.ID, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}
	session := models.NewSession(user.ID, user.Email, token)
	h.SessionRepo.Create(session, storage)

	c.SetCookie("token", token, 3600, "/", domain, false, true)
	c.JSON(http.StatusOK, gin.H{
		"message": "Login Successful",
	})

}

func (h *AuthHandler) Logout(c *gin.Context) {
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
