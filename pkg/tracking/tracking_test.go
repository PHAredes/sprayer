package tracking

import (
	"testing"
)

func TestEmailObfuscation(t *testing.T) {
	email := "test@example.com"
	obfuscated := ObfuscatedEmail(email)

	// Should be a valid MD5 hash (32 characters hex)
	if len(obfuscated) != 32 {
		t.Errorf("Obfuscated email should be 32 characters, got %d", len(obfuscated))
	}

	// Same email should produce same hash
	obfuscated2 := ObfuscatedEmail(email)
	if obfuscated != obfuscated2 {
		t.Errorf("Same email should produce same obfuscation hash")
	}

	// Different emails should produce different hashes
	differentEmail := "different@example.com"
	obfuscated3 := ObfuscatedEmail(differentEmail)
	if obfuscated == obfuscated3 {
		t.Errorf("Different emails should produce different obfuscation hashes")
	}
}

func TestTrackingURLs(t *testing.T) {
	jobID := "job-123"
	emailHash := "testhash123"
	originalURL := "https://example.com/jobs/123"

	// Test pixel URL
	pixelURL := GenerateTrackingPixelURL(jobID, emailHash)
	expectedPixelURL := "/track/open/job-123/testhash123.gif"
	if pixelURL != expectedPixelURL {
		t.Errorf("GenerateTrackingPixelURL() = %v, want %v", pixelURL, expectedPixelURL)
	}

	// Test link URL
	linkURL := GenerateTrackingLink(jobID, emailHash, originalURL)
	expectedLinkURL := "/track/click/job-123/testhash123?url=https://example.com/jobs/123"
	if linkURL != expectedLinkURL {
		t.Errorf("GenerateTrackingLink() = %v, want %v", linkURL, expectedLinkURL)
	}
}

func TestAnalyticsCalculation(t *testing.T) {
	events := []TrackingEvent{
		{Opened: true, Clicked: true, ApplicationSent: true},
		{Opened: true, Clicked: false, ApplicationSent: false},
		{Opened: false, Clicked: false, ApplicationSent: false},
	}

	analytics := CalculateAnalytics(events)

	if analytics.TotalEmailsSent != 3 {
		t.Errorf("TotalEmailsSent should be 3, got %d", analytics.TotalEmailsSent)
	}
	if analytics.EmailsOpened != 2 {
		t.Errorf("EmailsOpened should be 2, got %d", analytics.EmailsOpened)
	}
	if analytics.LinksClicked != 1 {
		t.Errorf("LinksClicked should be 1, got %d", analytics.LinksClicked)
	}
	if analytics.ApplicationsSent != 1 {
		t.Errorf("ApplicationsSent should be 1, got %d", analytics.ApplicationsSent)
	}

	// Check rates
	expectedOpenRate := 66.66666666666666
	if analytics.OpenRate != expectedOpenRate {
		t.Errorf("OpenRate should be %f, got %f", expectedOpenRate, analytics.OpenRate)
	}
}