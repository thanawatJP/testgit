package main

import (
	"authenservice/database"
	"authenservice/handler"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// เชื่อมต่อกับฐานข้อมูล
	database.Connect()

	// สร้างตัวแปร Gin router
	r := gin.Default()

	// ตั้งค่า CORS middleware
	r.Use(cors.Default()) // ใช้ค่าเริ่มต้นของ CORS

	// หรือถ้าต้องการตั้งค่าที่เฉพาะเจาะจง
	// r.Use(cors.New(cors.Config{
	//     AllowOrigins:     []string{"http://localhost:3000"},  // ตัวอย่าง: เฉพาะโดเมนนี้สามารถเข้าถึงได้
	//     AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
	//     AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
	//     AllowCredentials: true,
	// }))

	userGroup := r.Group("/user")
	{
		userGroup.GET("/", handler.GetAllUserHandler)
		userGroup.POST("/", handler.CreateUserHandler)
		userGroup.GET("/:id", handler.GetOneUserHandler)
	}
	roleGroup := r.Group("/role")
	{
		roleGroup.POST("/", handler.CreateRoleHandler)
		roleGroup.GET("/", handler.GetAllRoleHandler)
	}
	authGroup := r.Group("/auth")
	{
		authGroup.GET("/google", handler.GoogleCallbackHandler)
		authGroup.GET("/google/start", handler.GoogleStartHandler)
		authGroup.POST("/login", handler.NormalAuthHandler)
	}

	// เริ่ม server ที่ port 8080
	r.Run(":8080")
}
