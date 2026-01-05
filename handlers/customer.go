package handlers

import (
	"go-api/models"
	"go-api/repo"
	"go-api/utils"
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

func (h *CustomerHandler) Get(c *gin.Context) {
	readSource := utils.GetReadSource(c)

	customers, err := h.CustomerRepo.Get(readSource)

	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())

	}

	c.JSON(http.StatusOK, gin.H{"data": customers})
}

func (h *CustomerHandler) GetOne(c *gin.Context) {
	id := c.Param("id")
	readSource := utils.GetReadSource(c)

	customer, err := h.CustomerRepo.GetOne(id, readSource)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"data": customer})
}

func (h *CustomerHandler) GetByEmail(c *gin.Context) {
	email := c.Param("email")
	readSource := utils.GetReadSource(c)

	customer, err := h.CustomerRepo.GetByEmail(email, readSource)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"data": customer})
}

// update and soft delete and hard delete
func (h *CustomerHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var body struct {
		Storage  string `json:"storage"`
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Phone    string `json:"phone"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	customer, err := models.NewCustomer(body.Name, body.Email, body.Phone, body.Password)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.CustomerRepo.Update(id, customer)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}

func (h *CustomerHandler) SoftDelete(c *gin.Context) {
	id := c.Param("id")
	var err error

	err = h.CustomerRepo.SoftDelete(id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

func (h *CustomerHandler) HardDelete(c *gin.Context) {
	id := c.Param("id")
	var err error

	err = h.CustomerRepo.HardDelete(id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

func (h *CustomerHandler) AdminGet(c *gin.Context) {
	storage := utils.GetReadSource(c)
	var customers []*models.Customer
	var err error

	customers, err = h.CustomerRepo.AdminGet(storage)

	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
	}

	c.JSON(http.StatusOK, gin.H{"data": customers})
}

func (h *CustomerHandler) AdminGetOne(c *gin.Context) {
	id := c.Param("id")
	storage := utils.GetReadSource(c)
	var customer *models.Customer
	var err error

	customer, err = h.CustomerRepo.AdminGetOne(id, storage)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"data": customer})
}
