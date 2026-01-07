package repo

import (
	"go-api/models"

	"github.com/gin-gonic/gin"
)

// StripSensitiveData removes ID/IsDeleted and masks email/phone
func StripSensitiveData(customer *models.Customer) models.CustomerPublic {
	return models.CustomerPublic{
		Name:      customer.Name,
		Email:     maskString(customer.Email),
		Phone:     maskString(customer.Phone),
		CreatedAt: customer.CreatedAt,
		UpdatedAt: customer.UpdatedAt,
	}
}

// StripSensitiveDataList processes multiple customers
func StripSensitiveDataList(customers []*models.Customer) []models.CustomerPublic {
	result := make([]models.CustomerPublic, len(customers))
	for i, customer := range customers {
		result[i] = StripSensitiveData(customer)
	}
	return result
}

// maskString safely masks the middle of a string
func maskString(s string) string {
	if len(s) <= 6 {
		return "***"
	}
	return s[:3] + "..." + s[len(s)-3:]
}

func GetReadSource(c *gin.Context) ReadSource {
	storage := c.Query("storage")
	if storage == "sql" {
		return Postgres
	}
	return Mongo
}
