package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
)

// ساختار پاسخ JSON
type Response struct {
	Result string `json:"result"`
	Error  string `json:"error"`
}

// ساختار سرور
type Server struct {
	port string
}

// سازنده سرور
func NewServer(port string) *Server {
	return &Server{port: port}
}

// هندلر اصلی برای مسیرها
func (s *Server) handleOperation(w http.ResponseWriter, r *http.Request, op string) {
	w.Header().Set("Content-Type", "application/json")

	query := r.URL.Query().Get("numbers")
	if query == "" {
		resp := Response{"", "'numbers' parameter missing"}
		json.NewEncoder(w).Encode(resp)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	parts := strings.Split(query, ",")
	if len(parts) == 0 {
		resp := Response{"", "'numbers' parameter missing"}
		json.NewEncoder(w).Encode(resp)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var numbers []int64
	for _, p := range parts {
		n, err := strconv.ParseInt(strings.TrimSpace(p), 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			resp := Response{"", "invalid number format"}
			json.NewEncoder(w).Encode(resp)
			return
		}
		numbers = append(numbers, n)
	}

	if len(numbers) == 0 {
		resp := Response{"", "'numbers' parameter missing"}
		json.NewEncoder(w).Encode(resp)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var result int64
	result = numbers[0]

	for i := 1; i < len(numbers); i++ {
		prev := result
		if op == "add" {
			result += numbers[i]
		} else if op == "sub" {
			result -= numbers[i]
		}

		// بررسی overflow
		if (op == "add" && ((numbers[i] > 0 && result < prev) || (numbers[i] < 0 && result > prev))) ||
			(op == "sub" && ((numbers[i] < 0 && result < prev) || (numbers[i] > 0 && result > prev))) {
			w.WriteHeader(http.StatusBadRequest)
			resp := Response{"", "Overflow"}
			json.NewEncoder(w).Encode(resp)
			return
		}
	}

	// بررسی overflow در محدوده int64
	if result > math.MaxInt64 || result < math.MinInt64 {
		w.WriteHeader(http.StatusBadRequest)
		resp := Response{"", "Overflow"}
		json.NewEncoder(w).Encode(resp)
		return
	}

	// پاسخ موفق
	resp := Response{
		Result: fmt.Sprintf("The result of your query is: %d", result),
		Error:  "",
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// متد شروع سرور
func (s *Server) Start() {
	http.HandleFunc("/add", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(Response{"", "Invalid method"})
			return
		}
		s.handleOperation(w, r, "add")
	})

	http.HandleFunc("/sub", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(Response{"", "Invalid method"})
			return
		}
		s.handleOperation(w, r, "sub")
	})

	fmt.Printf("Server is running on port %s\n", s.port)
	log.Fatal(http.ListenAndServe(":"+s.port, nil))
}

// تابع main برای اجرا
func main() {
	s := NewServer("8080")
	s.Start()
}
