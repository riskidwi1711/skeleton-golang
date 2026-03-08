package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type Asset struct {
	ID             string    `json:"id"`
	Name           string    `json:"name"`
	Category       string    `json:"category"`
	Type           string    `json:"type"` // hardware or software
	Status         string    `json:"status"`
	AssignedTo     string    `json:"assigned_to"`
	Department     string    `json:"department"`
	Location       string    `json:"location"`
	SerialNumber   string    `json:"serial_number,omitempty"`
	LicenseKey     string    `json:"license_key,omitempty"`
	PurchaseDate   string    `json:"purchase_date"`
	WarrantyExpiry string    `json:"warranty_expiry,omitempty"`
	ExpiryDate     string    `json:"expiry_date,omitempty"`
	Value          string    `json:"value"`
	TenantID       string    `json:"tenant_id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type CreateAssetRequest struct {
	Name           string `json:"name"`
	Category       string `json:"category"`
	Type           string `json:"type"`
	Department     string `json:"department"`
	Location       string `json:"location"`
	SerialNumber   string `json:"serial_number,omitempty"`
	LicenseKey     string `json:"license_key,omitempty"`
	PurchaseDate   string `json:"purchase_date"`
	WarrantyExpiry string `json:"warranty_expiry,omitempty"`
	ExpiryDate     string `json:"expiry_date,omitempty"`
	Value          string `json:"value"`
}

type UpdateAssetRequest struct {
	Status     string `json:"status,omitempty"`
	AssignedTo string `json:"assigned_to,omitempty"`
	Location   string `json:"location,omitempty"`
	Department string `json:"department,omitempty"`
}

type AssetStats struct {
	Total    int `json:"total"`
	Hardware int `json:"hardware"`
	Software int `json:"software"`
	Active   int `json:"active"`
	Inactive int `json:"inactive"`
}

var db *sql.DB
var useMemory bool
var memoryAssets []Asset

func initDB() {
	var err error
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://mailhop_user:Mailhop2024!@localhost:5433/mailhop?sslmode=disable"
	}

	db, err = sql.Open("postgres", dbURL)
	if err != nil {
		log.Printf("DB unavailable, using in-memory mode: %v", err)
		useMemory = true
		return
	}
	if err := db.Ping(); err != nil {
		log.Printf("DB unavailable, using in-memory mode: %v", err)
		useMemory = true
		return
	}

	// Create assets table if not exists
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS assets (
		id VARCHAR(50) PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		category VARCHAR(100),
		type VARCHAR(50),
		status VARCHAR(50) DEFAULT 'active',
		assigned_to VARCHAR(100),
		department VARCHAR(100),
		location VARCHAR(200),
		serial_number VARCHAR(100),
		license_key VARCHAR(200),
		purchase_date VARCHAR(50),
		warranty_expiry VARCHAR(50),
		expiry_date VARCHAR(50),
		value VARCHAR(50),
		tenant_id VARCHAR(50),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Printf("Warning: Could not create table: %v", err)
	}

	// Create indexes for performance
	db.Exec("CREATE INDEX IF NOT EXISTS idx_assets_tenant ON assets(tenant_id)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_assets_type ON assets(type)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_assets_status ON assets(status)")
}

func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			// Allow without auth for development
			r.Header.Set("X-Tenant-ID", "default")
			r.Header.Set("X-User-Email", "dev@itms.com")
			next(w, r)
			return
		}

		// Parse JWT token if provided
		if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
			tokenString = tokenString[7:]
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		r.Header.Set("X-Tenant-ID", claims["tenant_id"].(string))
		r.Header.Set("X-User-Email", claims["email"].(string))
		next(w, r)
	}
}

func createAsset(w http.ResponseWriter, r *http.Request) {
	var req CreateAssetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	tenantID := r.Header.Get("X-Tenant-ID")

	// Generate asset ID
	assetCount := 0
	if useMemory {
		for _, a := range memoryAssets {
			if a.TenantID == tenantID {
				assetCount++
			}
		}
	} else {
		db.QueryRow("SELECT COUNT(*) FROM assets WHERE tenant_id = $1", tenantID).Scan(&assetCount)
	}
	assetID := fmt.Sprintf("AST-2024-%03d", assetCount+1)

	asset := Asset{
		ID:             assetID,
		Name:           req.Name,
		Category:       req.Category,
		Type:           req.Type,
		Status:         "active",
		AssignedTo:     "Unassigned",
		Department:     req.Department,
		Location:       req.Location,
		SerialNumber:   req.SerialNumber,
		LicenseKey:     req.LicenseKey,
		PurchaseDate:   req.PurchaseDate,
		WarrantyExpiry: req.WarrantyExpiry,
		ExpiryDate:     req.ExpiryDate,
		Value:          req.Value,
		TenantID:       tenantID,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if useMemory {
		memoryAssets = append([]Asset{asset}, memoryAssets...)
	} else {
		_, err := db.Exec(`
			INSERT INTO assets (id, name, category, type, status, assigned_to, department, location, 
			                   serial_number, license_key, purchase_date, warranty_expiry, expiry_date, value, tenant_id)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)`,
			asset.ID, asset.Name, asset.Category, asset.Type, asset.Status, asset.AssignedTo,
			asset.Department, asset.Location, asset.SerialNumber, asset.LicenseKey,
			asset.PurchaseDate, asset.WarrantyExpiry, asset.ExpiryDate, asset.Value, asset.TenantID)
		if err != nil {
			http.Error(w, "Failed to create asset", http.StatusInternalServerError)
			log.Printf("Error creating asset: %v", err)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(asset)
}

func getAssets(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Header.Get("X-Tenant-ID")

	// Get query parameters for filtering
	assetType := r.URL.Query().Get("type")
	status := r.URL.Query().Get("status")
	category := r.URL.Query().Get("category")
	department := r.URL.Query().Get("department")

	assets := []Asset{}
	if useMemory {
		for _, asset := range memoryAssets {
			if asset.TenantID != tenantID {
				continue
			}
			if assetType != "" && assetType != "all" && asset.Type != assetType {
				continue
			}
			if status != "" && status != "all" && asset.Status != status {
				continue
			}
			if category != "" && category != "all" && asset.Category != category {
				continue
			}
			if department != "" && department != "all" && asset.Department != department {
				continue
			}
			assets = append(assets, asset)
		}
	} else {
		query := `SELECT id, name, category, type, status, assigned_to, department, location, 
		          serial_number, license_key, purchase_date, warranty_expiry, expiry_date, value, created_at, updated_at
		          FROM assets WHERE tenant_id = $1`
		args := []interface{}{tenantID}
		argCount := 1
		if assetType != "" && assetType != "all" {
			argCount++
			query += fmt.Sprintf(" AND type = $%d", argCount)
			args = append(args, assetType)
		}
		if status != "" && status != "all" {
			argCount++
			query += fmt.Sprintf(" AND status = $%d", argCount)
			args = append(args, status)
		}
		if category != "" && category != "all" {
			argCount++
			query += fmt.Sprintf(" AND category = $%d", argCount)
			args = append(args, category)
		}
		if department != "" && department != "all" {
			argCount++
			query += fmt.Sprintf(" AND department = $%d", argCount)
			args = append(args, department)
		}
		query += " ORDER BY created_at DESC"
		rows, err := db.Query(query, args...)
		if err != nil {
			http.Error(w, "Failed to fetch assets", http.StatusInternalServerError)
			log.Printf("Error fetching assets: %v", err)
			return
		}
		defer rows.Close()
		for rows.Next() {
			var asset Asset
			var serialNumber, licenseKey, warrantyExpiry, expiryDate sql.NullString
			err := rows.Scan(&asset.ID, &asset.Name, &asset.Category, &asset.Type,
				&asset.Status, &asset.AssignedTo, &asset.Department, &asset.Location,
				&serialNumber, &licenseKey, &asset.PurchaseDate, &warrantyExpiry,
				&expiryDate, &asset.Value, &asset.CreatedAt, &asset.UpdatedAt)
			if err != nil {
				log.Printf("Error scanning asset: %v", err)
				continue
			}
			asset.TenantID = tenantID
			if serialNumber.Valid {
				asset.SerialNumber = serialNumber.String
			}
			if licenseKey.Valid {
				asset.LicenseKey = licenseKey.String
			}
			if warrantyExpiry.Valid {
				asset.WarrantyExpiry = warrantyExpiry.String
			}
			if expiryDate.Valid {
				asset.ExpiryDate = expiryDate.String
			}
			assets = append(assets, asset)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(assets)
}

func getAsset(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	assetID := vars["id"]
	tenantID := r.Header.Get("X-Tenant-ID")

	if useMemory {
		for _, a := range memoryAssets {
			if a.ID == assetID && a.TenantID == tenantID {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(a)
				return
			}
		}
		http.Error(w, "Asset not found", http.StatusNotFound)
		return
	}

	var asset Asset
	var serialNumber, licenseKey, warrantyExpiry, expiryDate sql.NullString
	err := db.QueryRow(`
		SELECT id, name, category, type, status, assigned_to, department, location,
		       serial_number, license_key, purchase_date, warranty_expiry, expiry_date, value, created_at, updated_at
		FROM assets WHERE id = $1 AND tenant_id = $2`,
		assetID, tenantID).Scan(
		&asset.ID, &asset.Name, &asset.Category, &asset.Type,
		&asset.Status, &asset.AssignedTo, &asset.Department, &asset.Location,
		&serialNumber, &licenseKey, &asset.PurchaseDate, &warrantyExpiry,
		&expiryDate, &asset.Value, &asset.CreatedAt, &asset.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Asset not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to fetch asset", http.StatusInternalServerError)
		}
		return
	}
	asset.TenantID = tenantID
	if serialNumber.Valid {
		asset.SerialNumber = serialNumber.String
	}
	if licenseKey.Valid {
		asset.LicenseKey = licenseKey.String
	}
	if warrantyExpiry.Valid {
		asset.WarrantyExpiry = warrantyExpiry.String
	}
	if expiryDate.Valid {
		asset.ExpiryDate = expiryDate.String
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(asset)
}

func updateAsset(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	assetID := vars["id"]
	tenantID := r.Header.Get("X-Tenant-ID")

	var req UpdateAssetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if useMemory {
		for i, a := range memoryAssets {
			if a.ID == assetID && a.TenantID == tenantID {
				if req.Status != "" {
					memoryAssets[i].Status = req.Status
				}
				if req.AssignedTo != "" {
					memoryAssets[i].AssignedTo = req.AssignedTo
				}
				if req.Location != "" {
					memoryAssets[i].Location = req.Location
				}
				if req.Department != "" {
					memoryAssets[i].Department = req.Department
				}
				memoryAssets[i].UpdatedAt = time.Now()
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(map[string]string{"message": "Asset updated successfully"})
				return
			}
		}
		http.Error(w, "Asset not found", http.StatusNotFound)
		return
	}

	updateFields := []string{}
	args := []interface{}{}
	argCount := 0
	if req.Status != "" {
		argCount++
		updateFields = append(updateFields, fmt.Sprintf("status = $%d", argCount))
		args = append(args, req.Status)
	}
	if req.AssignedTo != "" {
		argCount++
		updateFields = append(updateFields, fmt.Sprintf("assigned_to = $%d", argCount))
		args = append(args, req.AssignedTo)
	}
	if req.Location != "" {
		argCount++
		updateFields = append(updateFields, fmt.Sprintf("location = $%d", argCount))
		args = append(args, req.Location)
	}
	if req.Department != "" {
		argCount++
		updateFields = append(updateFields, fmt.Sprintf("department = $%d", argCount))
		args = append(args, req.Department)
	}
	if len(updateFields) == 0 {
		http.Error(w, "No fields to update", http.StatusBadRequest)
		return
	}
	argCount++
	updateFields = append(updateFields, fmt.Sprintf("updated_at = $%d", argCount))
	args = append(args, time.Now())
	argCount++
	args = append(args, assetID)
	argCount++
	args = append(args, tenantID)
	query := fmt.Sprintf("UPDATE assets SET %s WHERE id = $%d AND tenant_id = $%d", stringJoin(updateFields, ", "), argCount-1, argCount)
	result, err := db.Exec(query, args...)
	if err != nil {
		http.Error(w, "Failed to update asset", http.StatusInternalServerError)
		log.Printf("Error updating asset: %v", err)
		return
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Asset not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Asset updated successfully"})
}

func getAssetStats(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Header.Get("X-Tenant-ID")
	var stats AssetStats
	if useMemory {
		for _, a := range memoryAssets {
			if a.TenantID != tenantID {
				continue
			}
			stats.Total++
			if a.Type == "hardware" {
				stats.Hardware++
			}
			if a.Type == "software" {
				stats.Software++
			}
			if a.Status == "active" {
				stats.Active++
			}
			if a.Status == "inactive" {
				stats.Inactive++
			}
		}
	} else {
		db.QueryRow("SELECT COUNT(*) FROM assets WHERE tenant_id = $1", tenantID).Scan(&stats.Total)
		db.QueryRow("SELECT COUNT(*) FROM assets WHERE tenant_id = $1 AND type = 'hardware'", tenantID).Scan(&stats.Hardware)
		db.QueryRow("SELECT COUNT(*) FROM assets WHERE tenant_id = $1 AND type = 'software'", tenantID).Scan(&stats.Software)
		db.QueryRow("SELECT COUNT(*) FROM assets WHERE tenant_id = $1 AND status = 'active'", tenantID).Scan(&stats.Active)
		db.QueryRow("SELECT COUNT(*) FROM assets WHERE tenant_id = $1 AND status = 'inactive'", tenantID).Scan(&stats.Inactive)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

func deleteAsset(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	assetID := vars["id"]
	tenantID := r.Header.Get("X-Tenant-ID")

	if useMemory {
		for i, a := range memoryAssets {
			if a.ID == assetID && a.TenantID == tenantID {
				memoryAssets = append(memoryAssets[:i], memoryAssets[i+1:]...)
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(map[string]string{"message": "Asset deleted successfully"})
				return
			}
		}
		http.Error(w, "Asset not found", http.StatusNotFound)
		return
	}

	result, err := db.Exec("DELETE FROM assets WHERE id = $1 AND tenant_id = $2", assetID, tenantID)
	if err != nil {
		http.Error(w, "Failed to delete asset", http.StatusInternalServerError)
		return
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Asset not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Asset deleted successfully"})
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "healthy",
		"service": "asset-service",
	})
}

func stringJoin(strs []string, sep string) string {
	result := ""
	for i, s := range strs {
		if i > 0 {
			result += sep
		}
		result += s
	}
	return result
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8084"
	}

	initDB()
	defer db.Close()

	router := mux.NewRouter()

	// Health check
	router.HandleFunc("/health", healthCheck).Methods("GET")

	// Asset endpoints
	router.HandleFunc("/api/v1/assets", authMiddleware(createAsset)).Methods("POST")
	router.HandleFunc("/api/v1/assets", authMiddleware(getAssets)).Methods("GET")
	router.HandleFunc("/api/v1/assets/stats", authMiddleware(getAssetStats)).Methods("GET")
	router.HandleFunc("/api/v1/assets/{id}", authMiddleware(getAsset)).Methods("GET")
	router.HandleFunc("/api/v1/assets/{id}", authMiddleware(updateAsset)).Methods("PUT")
	router.HandleFunc("/api/v1/assets/{id}", authMiddleware(deleteAsset)).Methods("DELETE")

	// CORS middleware
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	})

	log.Printf("Asset Service running on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}