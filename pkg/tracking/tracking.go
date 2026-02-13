package tracking

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// TrackingEvent represents a tracking event for email opens
type TrackingEvent struct {
	ID          string    `json:"id"`
	JobID       string    `json:"job_id"`
	Recipient   string    `json:"recipient"`
	EmailHash   string    `json:"email_hash"`
	IPAddress   string    `json:"ip_address"`
	UserAgent   string    `json:"user_agent"`
	Timestamp   time.Time `json:"timestamp"`
	Opened      bool      `json:"opened"`
	Clicked     bool      `json:"clicked"`
	ApplicationSent bool  `json:"application_sent"`
}

// ObfuscatedEmail creates an obfuscated version of an email for tracking
func ObfuscatedEmail(email string) string {
	// Hash the email for privacy while maintaining tracking capability
	hash := md5.Sum([]byte(strings.ToLower(strings.TrimSpace(email))))
	return hex.EncodeToString(hash[:])
}

// GenerateTrackingPixelURL creates a URL for tracking email opens
func GenerateTrackingPixelURL(jobID, obfuscatedEmail string) string {
	// Create a tracking pixel URL that will be embedded in emails
	// This URL will be called when the email is opened
	return fmt.Sprintf("/track/open/%s/%s.gif", jobID, obfuscatedEmail)
}

// GenerateTrackingLink creates a tracking link for CV clicks
func GenerateTrackingLink(jobID, obfuscatedEmail, originalURL string) string {
	// Create a tracking link that redirects to the original URL
	return fmt.Sprintf("/track/click/%s/%s?url=%s", jobID, obfuscatedEmail, originalURL)
}

// TrackEmailOpen handles tracking when an email is opened
func TrackEmailOpen(w http.ResponseWriter, r *http.Request, jobID, emailHash string) {
	// Create a 1x1 transparent GIF for tracking pixel
	gifData := []byte{
		0x47, 0x49, 0x46, 0x38, 0x39, 0x61, // GIF89a
		0x01, 0x00, 0x01, 0x00, // 1x1 pixels
		0x80, 0x00, 0x00, // Background color
		0x00, 0x00, 0x00, // Aspect ratio
		0x2c, 0x00, 0x00, 0x00, 0x00, // Image descriptor
		0x01, 0x00, 0x01, 0x00, 0x00, // Local color table
		0x02, 0x02, 0x44, 0x01, 0x00, // Image data
		0x3b, // Trailer
	}

	// Set headers for GIF
	w.Header().Set("Content-Type", "image/gif")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(gifData)))
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	// Log the tracking event (in production, this would be stored in a database)
	ip := getClientIP(r)
	userAgent := r.UserAgent()
	
	fmt.Printf("Email opened - JobID: %s, EmailHash: %s, IP: %s, UserAgent: %s\n", 
		jobID, emailHash, ip, userAgent)

	// Write the GIF data
	w.Write(gifData)
}

// TrackLinkClick handles tracking when a link in the email is clicked
func TrackLinkClick(w http.ResponseWriter, r *http.Request, jobID, emailHash string) {
	// Get the original URL from query parameters
	originalURL := r.URL.Query().Get("url")
	if originalURL == "" {
		http.Error(w, "Missing URL parameter", http.StatusBadRequest)
		return
	}

	// Log the click event
	ip := getClientIP(r)
	userAgent := r.UserAgent()
	
	fmt.Printf("Link clicked - JobID: %s, EmailHash: %s, IP: %s, UserAgent: %s, URL: %s\n", 
		jobID, emailHash, ip, userAgent, originalURL)

	// Redirect to the original URL
	http.Redirect(w, r, originalURL, http.StatusFound)
}

// getClientIP extracts the client IP address from the request
func getClientIP(r *http.Request) string {
	// Check for X-Forwarded-For header (common in proxy setups)
	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		ips := strings.Split(forwarded, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	// Check for X-Real-IP header
	if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
		return realIP
	}

	// Fall back to remote address
	return strings.Split(r.RemoteAddr, ":")[0]
}

// AnalyticsData represents aggregated tracking data for reporting
type AnalyticsData struct {
	TotalEmailsSent    int     `json:"total_emails_sent"`
	EmailsOpened       int     `json:"emails_opened"`
	LinksClicked       int     `json:"links_clicked"`
	ApplicationsSent   int     `json:"applications_sent"`
	OpenRate          float64 `json:"open_rate"`
	ClickRate         float64 `json:"click_rate"`
	ResponseRate      float64 `json:"response_rate"`
}

// CalculateAnalytics calculates analytics from tracking events
func CalculateAnalytics(events []TrackingEvent) AnalyticsData {
	var analytics AnalyticsData
	
	if len(events) == 0 {
		return analytics
	}

	analytics.TotalEmailsSent = len(events)
	
	for _, event := range events {
		if event.Opened {
			analytics.EmailsOpened++
		}
		if event.Clicked {
			analytics.LinksClicked++
		}
		if event.ApplicationSent {
			analytics.ApplicationsSent++
		}
	}

	if analytics.TotalEmailsSent > 0 {
		analytics.OpenRate = float64(analytics.EmailsOpened) / float64(analytics.TotalEmailsSent) * 100
		analytics.ClickRate = float64(analytics.LinksClicked) / float64(analytics.TotalEmailsSent) * 100
		analytics.ResponseRate = float64(analytics.ApplicationsSent) / float64(analytics.TotalEmailsSent) * 100
	}

	return analytics
}