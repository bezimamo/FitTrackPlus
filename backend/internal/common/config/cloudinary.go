package config

import (
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/google/uuid"
)

// CloudinaryConfig holds Cloudinary configuration
type CloudinaryConfig struct {
	CloudName string
	APIKey    string
	APISecret string
	Client    *cloudinary.Cloudinary
}

// MediaUploadResult represents the result of a media upload
type MediaUploadResult struct {
	PublicID  string  `json:"public_id"`
	URL       string  `json:"url"`
	SecureURL string  `json:"secure_url"`
	Format    string  `json:"format"`
	Width     int     `json:"width"`
	Height    int     `json:"height"`
	Duration  float64 `json:"duration,omitempty"`
	Bytes     int     `json:"bytes"`
}

// NewCloudinaryConfig creates a new Cloudinary configuration
func NewCloudinaryConfig() *CloudinaryConfig {
	cloudName := os.Getenv("CLOUDINARY_CLOUD_NAME")
	apiKey := os.Getenv("CLOUDINARY_API_KEY")
	apiSecret := os.Getenv("CLOUDINARY_API_SECRET")

	if cloudName == "" || apiKey == "" || apiSecret == "" {
		log.Println("⚠️  Cloudinary credentials not found. Media uploads will be disabled.")
		return nil
	}

	// Initialize Cloudinary
	cld, err := cloudinary.NewFromParams(cloudName, apiKey, apiSecret)
	if err != nil {
		log.Printf("❌ Failed to initialize Cloudinary: %v", err)
		return nil
	}

	log.Println("✅ Cloudinary initialized successfully")
	return &CloudinaryConfig{
		CloudName: cloudName,
		APIKey:    apiKey,
		APISecret: apiSecret,
		Client:    cld,
	}
}

// UploadImage uploads an image to Cloudinary with fitness-specific transformations
func (c *CloudinaryConfig) UploadImage(ctx context.Context, file *multipart.FileHeader, folder string) (*MediaUploadResult, error) {
	if c.Client == nil {
		return nil, fmt.Errorf("cloudinary not configured")
	}

	// Open the uploaded file
	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer src.Close()

	// Generate unique filename
	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%s_%s%s", folder, uuid.New().String(), ext)

	// Upload parameters for fitness images
	uploadParams := uploader.UploadParams{
		PublicID:       fmt.Sprintf("fittrackplus/%s/%s", folder, filename),
		Folder:         fmt.Sprintf("fittrackplus/%s", folder),
		ResourceType:   "image",
		Transformation: "f_auto,q_auto,w_800,h_600,c_fill,g_auto", // Auto-optimize for fitness app
		Tags:           []string{"fitness", "exercise", folder}, // api.CldAPIArray is []string
	}

	// Upload to Cloudinary
	result, err := c.Client.Upload.Upload(ctx, src, uploadParams)
	if err != nil {
		return nil, fmt.Errorf("failed to upload image: %v", err)
	}

	return &MediaUploadResult{
		PublicID:  result.PublicID,
		URL:       result.URL,
		SecureURL: result.SecureURL,
		Format:    result.Format,
		Width:     result.Width,
		Height:    result.Height,
		Bytes:     result.Bytes,
	}, nil
}

// UploadVideo uploads a video to Cloudinary with fitness-specific optimizations
func (c *CloudinaryConfig) UploadVideo(ctx context.Context, file *multipart.FileHeader, folder string) (*MediaUploadResult, error) {
	if c.Client == nil {
		return nil, fmt.Errorf("cloudinary not configured")
	}

	// Open the uploaded file
	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer src.Close()

	// Generate unique filename
	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%s_%s%s", folder, uuid.New().String(), ext)

	// Upload parameters for fitness videos
	uploadParams := uploader.UploadParams{
		PublicID:       fmt.Sprintf("fittrackplus/%s/%s", folder, filename),
		Folder:         fmt.Sprintf("fittrackplus/%s", folder),
		ResourceType:   "video",
		Transformation: "f_auto,q_auto,w_1280,h_720,c_fill", // HD quality, auto-format
		Tags:           []string{"fitness", "exercise", "video", folder}, // api.CldAPIArray is []string
		Eager:          "f_mp4,q_auto,w_640,h_360,f_webm,q_auto,w_640,h_360", // String as expected by API
		EagerAsync:     &[]bool{true}[0], // Pointer to bool instead of bool
	}

	// Upload to Cloudinary
	result, err := c.Client.Upload.Upload(ctx, src, uploadParams)
	if err != nil {
		return nil, fmt.Errorf("failed to upload video: %v", err)
	}

	// Duration field doesn't exist in UploadResult, so we'll set it to 0
	// In a real implementation, you might need to extract this from video metadata separately
	duration := 0.0

	return &MediaUploadResult{
		PublicID:  result.PublicID,
		URL:       result.URL,
		SecureURL: result.SecureURL,
		Format:    result.Format,
		Width:     result.Width,
		Height:    result.Height,
		Duration:  duration,
		Bytes:     result.Bytes,
	}, nil
}

// DeleteMedia deletes a media file from Cloudinary
func (c *CloudinaryConfig) DeleteMedia(ctx context.Context, publicID string, resourceType string) error {
	if c.Client == nil {
		return fmt.Errorf("cloudinary not configured")
	}

	// Delete from Cloudinary
	_, err := c.Client.Upload.Destroy(ctx, uploader.DestroyParams{
		PublicID:     publicID,
		ResourceType: resourceType,
	})

	if err != nil {
		return fmt.Errorf("failed to delete media: %v", err)
	}

	return nil
}

// GenerateThumbnail generates a thumbnail from a video
func (c *CloudinaryConfig) GenerateThumbnail(ctx context.Context, videoPublicID string, timeOffset string) (string, error) {
	if c.Client == nil {
		return "", fmt.Errorf("cloudinary not configured")
	}

	// Generate thumbnail URL with transformation using string concatenation
	// Cloudinary URLs follow the pattern: https://res.cloudinary.com/cloud_name/video/upload/transformation/public_id
	thumbnailURL := fmt.Sprintf("https://res.cloudinary.com/%s/video/upload/w_300,h_200,c_fill,so_%s/%s", 
		c.CloudName, timeOffset, videoPublicID)

	return thumbnailURL, nil
}

// GetOptimizedURL generates an optimized URL for different screen sizes
func (c *CloudinaryConfig) GetOptimizedURL(publicID string, width, height int, format string) string {
	if c.Client == nil {
		return ""
	}

	// Generate optimized URL using string concatenation
	// Cloudinary URLs follow the pattern: https://res.cloudinary.com/cloud_name/image/upload/transformation/public_id
	url := fmt.Sprintf("https://res.cloudinary.com/%s/image/upload/f_%s,q_auto,w_%d,h_%d,c_fill/%s", 
		c.CloudName, format, width, height, publicID)

	return url
}

// IsVideoFile checks if the uploaded file is a video
func IsVideoFile(filename string) bool {
	videoExtensions := []string{".mp4", ".avi", ".mov", ".wmv", ".flv", ".webm", ".mkv"}
	ext := strings.ToLower(filepath.Ext(filename))
	
	for _, videoExt := range videoExtensions {
		if ext == videoExt {
			return true
		}
	}
	return false
}

// IsImageFile checks if the uploaded file is an image
func IsImageFile(filename string) bool {
	imageExtensions := []string{".jpg", ".jpeg", ".png", ".gif", ".bmp", ".webp", ".svg"}
	ext := strings.ToLower(filepath.Ext(filename))
	
	for _, imageExt := range imageExtensions {
		if ext == imageExt {
			return true
		}
	}
	return false
}
