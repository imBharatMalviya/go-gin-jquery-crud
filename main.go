package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/imBharatMalviya/go-gin-jquery-crud/models"
)

func main() {
	err := models.ConnectDatabase()
	checkErr(err)
	r := gin.Default()
	router := r.Group("/api")
	{
		router.GET("/employee", GetEmployees)
		router.GET("/employee/:id", getEmployeeById)
		router.DELETE("/employee/:id", DeleteEmployee)
		router.POST("/employee", CreateOrUpdateEmployee)
	}
	r.StaticFile("/", "./webui/")
	r.StaticFile("/form.html", "./webui/form.html")
	r.Run() // listen and serve on 0.0.0.0:8080
}

func GetEmployees(c *gin.Context) {
	employee, err := models.GetEmployees(10)
	checkErr(err)

	if employee == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No Records Found"})
		return
	} else {
		c.JSON(http.StatusOK, employee)
	}
}

func getEmployeeById(c *gin.Context) {

	id := c.Param("id")

	employee, err := models.GetEmployee(id)
	checkErr(err)
	// if the id is blank we can assume nothing is found
	if employee.Id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No Records Found"})
		return
	} else {
		c.JSON(http.StatusOK, employee)
	}
}

func CreateOrUpdateEmployee(c *gin.Context) {

	var json models.Employee

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"binding_error": err.Error()})
		return
	}

	success, err := models.SaveOrUpdateEmployee(json)

	if success {
		c.JSON(http.StatusOK, gin.H{"message": "Success"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}
}

func DeleteEmployee(c *gin.Context) {

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
	}

	success, err := models.DeleteEmployee(id)

	if success {
		c.JSON(http.StatusOK, gin.H{"message": "Success"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
