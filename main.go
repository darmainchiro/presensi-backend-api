package main

import (
	"net/http"
	"strconv"

	"backend-api/repository"
	"backend-api/service"
	"backend-api/config"
	"backend-api/utils"
	"backend-api/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	db := config.InitDB()
	defer db.Close()

	userRepo := repository.NewPostgresUserRepository(db)
	userService := service.NewUserService(userRepo)
	attendanceRepo := repository.NewPostgresAttendanceRepository(db)
	attendanceService := service.NewAttendanceService(attendanceRepo)

	r := gin.Default()

	r.POST("/api/login", func(c *gin.Context){
		email := c.PostForm("email")

		if email == "aji@example.com" {
			token, err := utils.GenerateToken(1)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat token"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"token": token})
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Email tidak terdaftar"})
		}
	})

	protected := r.Group("/api")
	protected.Use(middleware.ReqiuireAuth())
	{
		protected.POST("/check-in", func(c *gin.Context){
			userID := c.MustGet("user_id").(int64)
		
			err := attendanceService.ProcessCheckIn(c.Request.Context(), userID)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"message": "Check-in berhasil"})
		})

		protected.POST("/check-out", func(c *gin.Context){
			userID := c.MustGet("user_id").(int64)

			err := attendanceService.ProcessCheckOut(c.Request.Context(), userID)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"message": "Check-out berhasil"})
		})
	}

	r.GET("/api/users/:id", func(c *gin.Context){
		idStr := c.Param("id")
		userID, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil{
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
			return
		}

		user, err := userService.GetUserDetails(c.Request.Context(), userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if user == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User tidak ditemukan"})
			return
		}
		c.JSON(http.StatusOK, user)
	})

	r.Run(":8080")
	

}
