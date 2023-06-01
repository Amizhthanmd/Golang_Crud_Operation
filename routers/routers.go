package routers

import (
	"fmt"

	"Crud_operation_go/crudfunction"

	"github.com/gin-gonic/gin"
)

func Router() {
	crudfunction.SetupDB()

	go crudfunction.UpdateEmails()

	r := gin.Default()

	r.POST("/createusers", crudfunction.CreateUser)
	r.GET("/getusers/:id", crudfunction.GetUser)
	r.GET("/getallusers", crudfunction.GetAllUsers)
	r.PUT("/updateusers/:id", crudfunction.UpdateUser)
	r.DELETE("/deleteusers/:id", crudfunction.DeleteUser)

	if err := r.Run(":8080"); err != nil {
		fmt.Println("Failed to start the server:", err)
	}
}
