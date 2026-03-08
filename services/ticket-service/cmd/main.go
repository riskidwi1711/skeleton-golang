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

type Ticket struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	Priority    string    `json:"priority"`
	Category    string    `json:"category"`
	Assignee    string    `json:"assignee"`
	Requester   string    `json:"requester"`
	TenantID    string    `json:"tenant_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	SLAResponse string    `json:"sla_response"`
	SLAResolve  string    `json:"sla_resolve"`
}

type CreateTicketRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Category    string `json:"category"`
	Priority    string `json:"priority"`
}

type UpdateTicketRequest struct {
	Status   string `json:"status,omitempty"`
	Assignee string `json:"assignee,omitempty"`
	Priority string `json:"priority,omitempty"`
}

type TicketStats struct {
	Total      int `json:"total"`
	Open       int `json:"open"`
	InProgress int `json:"in_progress"`
	Resolved   int `json:"resolved"`
	Closed     int `json:"closed"`
}

var db *sql.DB
var useMemory bool
var memoryTickets []Ticket

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

	// Create tickets table if not exists
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS tickets (
		id VARCHAR(50) PRIMARY KEY,
		title VARCHAR(255) NOT NULL,
		description TEXT,
		status VARCHAR(50) DEFAULT 'open',
		priority VARCHAR(50) DEFAULT 'medium',
		category VARCHAR(100),
		assignee VARCHAR(100),
		requester VARCHAR(100),
		tenant_id VARCHAR(50),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		sla_response VARCHAR(20) DEFAULT '2h',
		sla_resolve VARCHAR(20) DEFAULT '8h'
	);`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Printf("Warning: Could not create table: %v", err)
	}

	// Create index for performance
	db.Exec("CREATE INDEX IF NOT EXISTS idx_tickets_tenant ON tickets(tenant_id)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_tickets_status ON tickets(status)")
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

func createTicket(w http.ResponseWriter, r *http.Request) {
	var req CreateTicketRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	tenantID := r.Header.Get("X-Tenant-ID")
	requester := r.Header.Get("X-User-Email")

	// Generate ticket ID
	ticketCount := 0
	if useMemory {
		for _, t := range memoryTickets {
			if t.TenantID == tenantID {
				ticketCount++
			}
		}
	} else {
		db.QueryRow("SELECT COUNT(*) FROM tickets WHERE tenant_id = $1", tenantID).Scan(&ticketCount)
	}
	ticketID := fmt.Sprintf("INC-2024-%03d", ticketCount+1)

	// Set SLA based on priority
	slaResponse, slaResolve := "2h", "8h"
	switch req.Priority {
	case "high":
		slaResponse, slaResolve = "1h", "4h"
	case "low":
		slaResponse, slaResolve = "4h", "48h"
	}

	ticket := Ticket{
		ID:          ticketID,
		Title:       req.Title,
		Description: req.Description,
		Status:      "open",
		Priority:    req.Priority,
		Category:    req.Category,
		Assignee:    "Unassigned",
		Requester:   requester,
		TenantID:    tenantID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		SLAResponse: slaResponse,
		SLAResolve:  slaResolve,
	}

	if useMemory {
		memoryTickets = append([]Ticket{ticket}, memoryTickets...)
	} else {
		_, err := db.Exec(`
			INSERT INTO tickets (id, title, description, status, priority, category, assignee, requester, tenant_id, sla_response, sla_resolve)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
			ticket.ID, ticket.Title, ticket.Description, ticket.Status, ticket.Priority,
			ticket.Category, ticket.Assignee, ticket.Requester, ticket.TenantID,
			ticket.SLAResponse, ticket.SLAResolve)
		if err != nil {
			http.Error(w, "Failed to create ticket", http.StatusInternalServerError)
			log.Printf("Error creating ticket: %v", err)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ticket)
}

func getTickets(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Header.Get("X-Tenant-ID")
	
	// Get query parameters for filtering
	status := r.URL.Query().Get("status")
	priority := r.URL.Query().Get("priority")
	
	tickets := []Ticket{}
	if useMemory {
		for _, ticket := range memoryTickets {
			if ticket.TenantID != tenantID {
				continue
			}
			if status != "" && status != "all" && ticket.Status != status {
				continue
			}
			if priority != "" && priority != "all" && ticket.Priority != priority {
				continue
			}
			tickets = append(tickets, ticket)
		}
	} else {
		query := "SELECT id, title, description, status, priority, category, assignee, requester, created_at, updated_at, sla_response, sla_resolve FROM tickets WHERE tenant_id = $1"
		args := []interface{}{tenantID}
		argCount := 1
		if status != "" && status != "all" {
			argCount++
			query += fmt.Sprintf(" AND status = $%d", argCount)
			args = append(args, status)
		}
		if priority != "" && priority != "all" {
			argCount++
			query += fmt.Sprintf(" AND priority = $%d", argCount)
			args = append(args, priority)
		}
		query += " ORDER BY created_at DESC"
		rows, err := db.Query(query, args...)
		if err != nil {
			http.Error(w, "Failed to fetch tickets", http.StatusInternalServerError)
			log.Printf("Error fetching tickets: %v", err)
			return
		}
		defer rows.Close()
		for rows.Next() {
			var ticket Ticket
			err := rows.Scan(&ticket.ID, &ticket.Title, &ticket.Description, &ticket.Status,
				&ticket.Priority, &ticket.Category, &ticket.Assignee, &ticket.Requester,
				&ticket.CreatedAt, &ticket.UpdatedAt, &ticket.SLAResponse, &ticket.SLAResolve)
			if err != nil {
				log.Printf("Error scanning ticket: %v", err)
				continue
			}
			ticket.TenantID = tenantID
			tickets = append(tickets, ticket)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tickets)
}

func getTicket(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ticketID := vars["id"]
	tenantID := r.Header.Get("X-Tenant-ID")

	if useMemory {
		for _, t := range memoryTickets {
			if t.ID == ticketID && t.TenantID == tenantID {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(t)
				return
			}
		}
		http.Error(w, "Ticket not found", http.StatusNotFound)
		return
	}

	var ticket Ticket
	err := db.QueryRow(`
		SELECT id, title, description, status, priority, category, assignee, requester, created_at, updated_at, sla_response, sla_resolve
		FROM tickets WHERE id = $1 AND tenant_id = $2`,
		ticketID, tenantID).Scan(
		&ticket.ID, &ticket.Title, &ticket.Description, &ticket.Status,
		&ticket.Priority, &ticket.Category, &ticket.Assignee, &ticket.Requester,
		&ticket.CreatedAt, &ticket.UpdatedAt, &ticket.SLAResponse, &ticket.SLAResolve)

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Ticket not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to fetch ticket", http.StatusInternalServerError)
		}
		return
	}

	ticket.TenantID = tenantID
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ticket)
}

func updateTicket(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ticketID := vars["id"]
	tenantID := r.Header.Get("X-Tenant-ID")

	var req UpdateTicketRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if useMemory {
		for i, t := range memoryTickets {
			if t.ID == ticketID && t.TenantID == tenantID {
				if req.Status != "" {
					memoryTickets[i].Status = req.Status
				}
				if req.Assignee != "" {
					memoryTickets[i].Assignee = req.Assignee
				}
				if req.Priority != "" {
					memoryTickets[i].Priority = req.Priority
				}
				memoryTickets[i].UpdatedAt = time.Now()
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(map[string]string{"message": "Ticket updated successfully"})
				return
			}
		}
		http.Error(w, "Ticket not found", http.StatusNotFound)
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
	if req.Assignee != "" {
		argCount++
		updateFields = append(updateFields, fmt.Sprintf("assignee = $%d", argCount))
		args = append(args, req.Assignee)
	}
	if req.Priority != "" {
		argCount++
		updateFields = append(updateFields, fmt.Sprintf("priority = $%d", argCount))
		args = append(args, req.Priority)
	}
	if len(updateFields) == 0 {
		http.Error(w, "No fields to update", http.StatusBadRequest)
		return
	}
	argCount++
	updateFields = append(updateFields, fmt.Sprintf("updated_at = $%d", argCount))
	args = append(args, time.Now())
	argCount++
	args = append(args, ticketID)
	argCount++
	args = append(args, tenantID)
	query := fmt.Sprintf("UPDATE tickets SET %s WHERE id = $%d AND tenant_id = $%d", stringJoin(updateFields, ", "), argCount-1, argCount)
	result, err := db.Exec(query, args...)
	if err != nil {
		http.Error(w, "Failed to update ticket", http.StatusInternalServerError)
		log.Printf("Error updating ticket: %v", err)
		return
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Ticket not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Ticket updated successfully"})
}

func getTicketStats(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Header.Get("X-Tenant-ID")
	var stats TicketStats
	if useMemory {
		for _, t := range memoryTickets {
			if t.TenantID != tenantID {
				continue
			}
			stats.Total++
			switch t.Status {
			case "open":
				stats.Open++
			case "in-progress":
				stats.InProgress++
			case "resolved":
				stats.Resolved++
			case "closed":
				stats.Closed++
			}
		}
	} else {
		db.QueryRow("SELECT COUNT(*) FROM tickets WHERE tenant_id = $1", tenantID).Scan(&stats.Total)
		db.QueryRow("SELECT COUNT(*) FROM tickets WHERE tenant_id = $1 AND status = 'open'", tenantID).Scan(&stats.Open)
		db.QueryRow("SELECT COUNT(*) FROM tickets WHERE tenant_id = $1 AND status = 'in-progress'", tenantID).Scan(&stats.InProgress)
		db.QueryRow("SELECT COUNT(*) FROM tickets WHERE tenant_id = $1 AND status = 'resolved'", tenantID).Scan(&stats.Resolved)
		db.QueryRow("SELECT COUNT(*) FROM tickets WHERE tenant_id = $1 AND status = 'closed'", tenantID).Scan(&stats.Closed)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "healthy",
		"service": "ticket-service",
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
		port = "8083"
	}

	initDB()
	defer db.Close()

	router := mux.NewRouter()

	// Health check
	router.HandleFunc("/health", healthCheck).Methods("GET")

	// Ticket endpoints
	router.HandleFunc("/api/v1/tickets", authMiddleware(createTicket)).Methods("POST")
	router.HandleFunc("/api/v1/tickets", authMiddleware(getTickets)).Methods("GET")
	router.HandleFunc("/api/v1/tickets/stats", authMiddleware(getTicketStats)).Methods("GET")
	router.HandleFunc("/api/v1/tickets/{id}", authMiddleware(getTicket)).Methods("GET")
	router.HandleFunc("/api/v1/tickets/{id}", authMiddleware(updateTicket)).Methods("PUT")

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

	log.Printf("Ticket Service running on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}