package handlers

import (
	"go-api/models"
	"go-api/repo"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CustomerHandler struct {
	CustomerRepo *repo.RepoFactory[models.Customer]
}

func NewCustomerHandler() *CustomerHandler {
	return &CustomerHandler{
		CustomerRepo: repo.NewRepoFactory(
			repo.NewPostgresRepo[models.Customer]("customers"),
			repo.NewMongoRepo[models.Customer]("customers"),
		),
	}
}

func (h *CustomerHandler) Create(c *gin.Context) {
	storage := repo.GetReadSource(c)
	var body struct {
		Email string `json:"email" binding:"required,email"`
		Phone string `json:"phone" binding:"required"`
		Name  string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	customer, err := models.NewCustomer(body.Name, body.Email, body.Phone)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create customer"})
		return
	}

	if err := h.CustomerRepo.Create(customer, storage); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create customer"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Customer created successfully",
		"data":    customer,
	})
}

func (h *CustomerHandler) Get(c *gin.Context) {
	readSource := repo.GetReadSource(c)

	customers, err := h.CustomerRepo.Get(readSource)

	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())

	}

	c.JSON(http.StatusOK, gin.H{"data": customers})
}
