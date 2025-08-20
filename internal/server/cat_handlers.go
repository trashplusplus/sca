package server

import (
	"fmt"
	"sca/internal/entity/cat"
	"sca/internal/validator"
	"strconv"

	"github.com/gin-gonic/gin"
)

func createCatHandler(repo *cat.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var cat cat.Cat

		if err := c.ShouldBindJSON(&cat); err != nil {
			c.JSON(400, gin.H{"error": "invalid request body"})
			return
		}

		//breed validation
		isValidBreed, err := validator.ValidateBreed(cat.Breed)
		if err != nil {
			c.JSON(500, gin.H{"error": "failed to validate breed"})
			return
		}

		if !isValidBreed {
			c.JSON(404, gin.H{"error": "invalid breed"})
			return
		}

		if err := repo.Create(&cat); err != nil {
			c.JSON(500, gin.H{"error": "failed to create cat"})
			return
		}

		c.JSON(200, cat)
	}
}

func listCatHandler(repo *cat.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		cats, err := repo.GetAll()
		if err != nil {
			c.JSON(500, gin.H{"error": "failed to get cats"})
			return
		}
		c.JSON(200, cats)
	}
}

func updateCatHandler(repo *cat.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {

		var cat cat.Cat
		//валідуємо реквест баді
		if err := c.ShouldBindJSON(&cat); err != nil {
			c.JSON(400, gin.H{"error": "invalid request body"})
			return
		}
		//якщо ок то шукаємо кота за айді
		targetCat, err := repo.GetByID(cat.Id)
		if err != nil {
			c.JSON(404, gin.H{"error": "failed to get cat"})
			return
		}

		targetCat.Name = cat.Name
		targetCat.Experience = cat.Experience
		targetCat.Breed = cat.Breed
		targetCat.Salary = cat.Salary

		if err := repo.Update(targetCat); err != nil {
			c.JSON(500, gin.H{"error": "failed to update cat"})
			return
		}

		c.JSON(200, gin.H{"message": fmt.Sprintf("Cat ID: [%d] updated succesfully", targetCat.Id)})

	}
}

func deleteCatHandler(repo *cat.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		id_string := c.Param("id")
		//конвертуємо id string в id int
		id, err := strconv.Atoi(id_string)
		if err != nil {
			c.JSON(400, gin.H{"error": "invalid request body"})
			return
		}
		//шукаємо кота, щоб було що видаляти
		targetCat, err := repo.GetByID(id)
		if err != nil {
			c.JSON(404, gin.H{"error": "failed to get cat"})
			return
		}

		if err := repo.Delete(targetCat.Id); err != nil {
			c.JSON(500, gin.H{"error": "failed to delete cat"})
			return
		}
		c.JSON(200, gin.H{"message": fmt.Sprintf("Cat ID: [%d] deleted succesfully", id)})
	}
}

func getCatHandler(repo *cat.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		id_string := c.Param("id")
		id, err := strconv.Atoi(id_string)
		if err != nil {
			c.JSON(400, gin.H{"error": "wrong param"})
			return
		}
		targetCat, err := repo.GetByID(id)
		if err != nil {
			c.JSON(404, gin.H{"error": "failed to get cat"})
			return
		}
		c.JSON(200, targetCat)
	}
}
