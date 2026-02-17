package server

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/evanqhuang/resume-cli/generator"
	"github.com/evanqhuang/resume-cli/matching"
	"github.com/evanqhuang/resume-cli/resume"
	"github.com/go-chi/chi/v5"
)

// resumeCache caches the loaded resume with modification time tracking
type resumeCache struct {
	mu       sync.RWMutex
	resume   *resume.Resume
	modTime  time.Time
	filePath string
}

var cache = &resumeCache{}

func registerRoutes(r chi.Router, s *Server) {
	// Initialize cache with resume path
	cache.filePath = s.resumePath

	r.Route("/api", func(r chi.Router) {
		r.Get("/health", handleHealth)
		r.Get("/resume", s.handleGetResume)
		r.Post("/resume/reload", s.handleReloadResume)
		r.Post("/job/analyze", s.handleAnalyzeJob)
		r.Post("/generate", s.handleGenerate)
	})
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}

// loadResume loads the resume from cache or disk
func loadResume(forceReload bool) (*resume.Resume, error) {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	// Check file modification time
	stat, err := os.Stat(cache.filePath)
	if err != nil {
		return nil, err
	}

	// Return cached version if still valid
	if !forceReload && cache.resume != nil && !stat.ModTime().After(cache.modTime) {
		return cache.resume, nil
	}

	// Load from disk
	log.Printf("Loading resume from %s", cache.filePath)
	r, err := resume.LoadResume(cache.filePath)
	if err != nil {
		return nil, err
	}

	// Update cache
	cache.resume = r
	cache.modTime = stat.ModTime()

	return r, nil
}

func (s *Server) handleGetResume(w http.ResponseWriter, r *http.Request) {
	res, err := loadResume(false)
	if err != nil {
		log.Printf("Error loading resume: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	transformed := TransformResume(res)
	if err := json.NewEncoder(w).Encode(transformed); err != nil {
		log.Printf("Error encoding response: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (s *Server) handleReloadResume(w http.ResponseWriter, r *http.Request) {
	res, err := loadResume(true)
	if err != nil {
		log.Printf("Error reloading resume: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	transformed := TransformResume(res)
	if err := json.NewEncoder(w).Encode(transformed); err != nil {
		log.Printf("Error encoding response: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// JobAnalysisRequest represents the request body for job analysis
type JobAnalysisRequest struct {
	JobTitle    string `json:"job_title"`
	Company     string `json:"company"`
	Description string `json:"description"`
}

func (s *Server) handleAnalyzeJob(w http.ResponseWriter, r *http.Request) {
	var req JobAnalysisRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding request: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid request body"})
		return
	}

	if req.Description == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "description is required"})
		return
	}

	res, err := loadResume(false)
	if err != nil {
		log.Printf("Error loading resume: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	result, err := matching.AnalyzeJobForAPI(res, req.JobTitle, req.Company, req.Description)
	if err != nil {
		log.Printf("Error analyzing job: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	if err := json.NewEncoder(w).Encode(result); err != nil {
		log.Printf("Error encoding response: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// GenerateRequest represents the request body for PDF generation
type GenerateRequest struct {
	Selections map[string][]string `json:"selections"`
	Template   string              `json:"template"`
}

func (s *Server) handleGenerate(w http.ResponseWriter, r *http.Request) {
	var req GenerateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding request: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid request body"})
		return
	}

	res, err := loadResume(false)
	if err != nil {
		log.Printf("Error loading resume: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	// Convert selections to map[string]bool
	selectedIDs := make(map[string]bool)
	for _, ids := range req.Selections {
		for _, id := range ids {
			selectedIDs[id] = true
		}
	}

	// Generate PDF
	pdfBytes, err := generator.GeneratePDF(res, selectedIDs)
	if err != nil {
		log.Printf("Error generating PDF: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	// Return PDF
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename=resume.pdf")
	w.WriteHeader(http.StatusOK)
	w.Write(pdfBytes)
}
