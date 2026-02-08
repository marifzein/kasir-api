package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"kasir-api/database"
	"kasir-api/handlers"
	"kasir-api/repositories"
	"kasir-api/services"

	"github.com/spf13/viper"
)

type Config struct {
	Port   string `mapstructure:"PORT"`
	DBConn string `mapstructure:"DB_CONN"`
}

func main() {
	// 1. Inisialisasi Viper
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigFile(".env")
		_ = viper.ReadInConfig()
	}

	// if err := viper.ReadInConfig(); err != nil {
	// 	log.Printf("Peringatan: Tidak ada file config, menggunakan env vars atau default: %v", err)
	// }

	config := Config{
		Port:   viper.GetString("PORT"),
		DBConn: viper.GetString("DB_CONN"),
	}

	// 2. Setup database
	db, err := database.InitDB(config.DBConn)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	// 3. Layering Architecture (Dependency Injection)
	productRepo := repositories.NewProductRepository(db)
	productService := services.NewProductService(productRepo)
	productHandler := handlers.NewProductHandler(productService)

	// Transaction
	transactionRepo := repositories.NewTransactionRepository(db)
	transactionService := services.NewTransactionService(transactionRepo)
	transactionHandler := handlers.NewTransactionHandler(transactionService)

	// Report
	reportRepo := repositories.NewReportRepository(db)
	reportService := services.NewReportService(reportRepo)
	reportHandler := handlers.NewReportHandler(reportService)


	// 4. Setup routes
	http.HandleFunc("/api/products", productHandler.HandleProducts)
	http.HandleFunc("/api/products/", productHandler.HandleProductByID)
	http.HandleFunc("/api/checkout", transactionHandler.HandleCheckout) // POST
	// Report
	http.HandleFunc("/api/report/hari-ini", reportHandler.HandleDailyReport)
	// Report Range (optional challenge)
	http.HandleFunc("/api/report", reportHandler.HandleReportRange)

	// 5. Jalankan Server
	// --- Jalankan Server Sesuai Petunjuk Course ---
	// addr := "0.0.0.0:" + config.Port
	// fmt.Println("Server running di", addr)

	// err = http.ListenAndServe(addr, nil)
	// if err != nil {
	// 	fmt.Println("gagal running server", err)
	// }

	port := config.Port
	if port == "" {
			port = "8080" // Default jika .env tidak terbaca
	}
	addr := "0.0.0.0:" + port
	fmt.Println("Server running di", addr)

	// Tambahkan baris ini untuk menjaga server tetap hidup
	err = http.ListenAndServe(addr, nil) 
	if err != nil {
		log.Fatal("Gagal menjalankan server:", err)
	}

}