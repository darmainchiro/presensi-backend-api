package main

import (
	"net/http"
	"strconv"

	"backend-api/config"
	"backend-api/middleware"
	"backend-api/repository"
	"backend-api/service"
	"backend-api/utils"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "backend-api/docs"
)

// 1. Buat Struct AppHandler untuk menampung dependency service
type AppHandler struct {
	userService       service.UserService
	attendanceService service.AttendanceService
}

// 2. Ekstrak Fungsi Login
// @Summary Login Pengguna
// @Description Endpoint untuk mendapatkan token JWT dengan email.
// @Tags Auth
// @Accept application/x-www-form-urlencoded
// @Produce json
// @Param email formData string true "Email pengguna (misal: aji@example.com)"
// @Success 200 {object} map[string]interface{} "Token berhasil dibuat"
// @Failure 401 {object} map[string]interface{} "Email tidak terdaftar"
// @Failure 500 {object} map[string]interface{} "Gagal membuat token"
// @Router /login [post]
func (h *AppHandler) LoginAPI(c *gin.Context) {
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
}

// 3. Ekstrak Fungsi Check-In
// @Summary Melakukan presensi masuk
// @Description Endpoint untuk mencatat waktu check-in pengguna.
// @Tags Attendance
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "Check-in berhasil"
// @Failure 400 {object} map[string]interface{} "Bad Request / Gagal Check-in"
// @Router /check-in [post]
func (h *AppHandler) CheckInAPI(c *gin.Context) {
	userID := c.MustGet("user_id").(int64)

	err := h.attendanceService.ProcessCheckIn(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Check-in berhasil"})
}

// 4. Ekstrak Fungsi Check-Out
// @Summary Melakukan presensi pulang
// @Description Endpoint untuk mencatat waktu check-out pengguna.
// @Tags Attendance
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "Check-out berhasil"
// @Failure 400 {object} map[string]interface{} "Bad Request / Gagal Check-out"
// @Router /check-out [post]
func (h *AppHandler) CheckOutAPI(c *gin.Context) {
	userID := c.MustGet("user_id").(int64)

	err := h.attendanceService.ProcessCheckOut(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Check-out berhasil"})
}

// 5. Ekstrak Fungsi Get User
// @Summary Mendapatkan detail user
// @Description Endpoint untuk mengambil data user berdasarkan ID.
// @Tags Users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} map[string]interface{} "Data user"
// @Failure 400 {object} map[string]interface{} "ID tidak valid"
// @Failure 404 {object} map[string]interface{} "User tidak ditemukan"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /users/{id} [get]
func (h *AppHandler) GetUserAPI(c *gin.Context) {
	idStr := c.Param("id")
	userID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	user, err := h.userService.GetUserDetails(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, user)
}

// =====================================================================
// Anotasi Global Swagger
// =====================================================================

// @title API Presensi Karyawan
// @version 1.0
// @description Layanan backend untuk sistem check-in dan check-out terintegrasi Kafka.
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email aji@example.com
// @host localhost:8080
// @BasePath /api
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
	db := config.InitDB()
	defer db.Close()

	userRepo := repository.NewPostgresUserRepository(db)
	userService := service.NewUserService(userRepo)
	attendanceRepo := repository.NewPostgresAttendanceRepository(db)
	attendanceService := service.NewAttendanceService(attendanceRepo)

	// 6. Inisialisasi Handler Struct dengan memasukkan service ke dalamnya
	appHandler := &AppHandler{
		userService:       userService,
		attendanceService: attendanceService,
	}

	r := gin.Default()

	// Rute UI Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 7. Gunakan method dari appHandler untuk rute-rute API
	r.POST("/api/login", appHandler.LoginAPI)

	protected := r.Group("/api")
	protected.Use(middleware.ReqiuireAuth()) // Pastikan nama fungsi middleware Anda (ReqiuireAuth) sesuai
	{
		protected.POST("/check-in", appHandler.CheckInAPI)
		protected.POST("/check-out", appHandler.CheckOutAPI)
	}

	r.GET("/api/users/:id", appHandler.GetUserAPI)

	go utils.StartKafkaConsumer()

	r.Run(":8080")
}