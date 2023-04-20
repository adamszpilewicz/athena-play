package router

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"s3-interaction/loggerCdc"
	"s3-interaction/s3Client"
)

// Router is the main router of the application
type Router struct {
	S3Client *s3Client.S3Client
	Logger   loggerCdc.Logger
}

// NewRouter creates a new router
func NewRouter(s3Client *s3Client.S3Client) (*Router, error) {
	customLogger, err := loggerCdc.NewLogger("router", true)
	if err != nil {
		return nil, err
	}

	return &Router{
		S3Client: s3Client,
		Logger:   customLogger,
	}, nil
}

// SetupRoutes sets up the routes
func (r *Router) SetupRoutes() {
	http.HandleFunc("/downloadToFile", r.HandleDownloadToFile)
	http.HandleFunc("/upload", r.HandleUpload)
}

// HandleDownloadToFile handles the download request
func (r *Router) HandleDownloadToFile(w http.ResponseWriter, req *http.Request) {
	bucket, key, err := getBucketAndKey(req)
	r.Logger.Info(fmt.Sprintf("Downloading file from S3: %s", filepath.Join(bucket, key)))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	localPath := req.URL.Query().Get("path")
	if localPath == "" {
		http.Error(w, "Missing 'path' query parameter", http.StatusBadRequest)
		return
	}

	fileContent, err := r.S3Client.DownloadFile(bucket, key)
	if err != nil {
		http.Error(w, "Failed to download the file", http.StatusInternalServerError)
		r.Logger.Error(fmt.Sprintf("Failed to download the file: %s", err.Error()))
		return
	}

	saveFileToLocalPath(bucket, localPath, key, fileContent, w)
}

// HandleUpload handles the upload request
func (r *Router) HandleUpload(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, "Invalid method, only POST is allowed", http.StatusMethodNotAllowed)
		return
	}

	file, fileHeader, err := req.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to read the file from request", http.StatusBadRequest)
		return
	}
	r.Logger.Info(fmt.Sprintf("File name got from request: %s", fileHeader.Filename))
	defer file.Close()

	bucket, key, err := getBucketAndKeyFromRequest(req)
	if err != nil {
		http.Error(w, "Missing 'bucket' or 'key' form field", http.StatusBadRequest)
		return
	}

	key = buildKeyWithOriginalFilename(key, fileHeader.Filename)

	fileContent, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Failed to read the file content", http.StatusInternalServerError)
		return
	}

	r.Logger.Info(fmt.Sprintf("Uploading file to S3: %s", filepath.Join(bucket, key)))
	err = r.S3Client.UploadFile(bucket, fileContent, key)
	if err != nil {
		http.Error(w, "Failed to upload the file", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "File uploaded successfully to S3: %s", filepath.Join(bucket, key))
}

// getBucketAndKey gets the bucket and key from the request
func getBucketAndKey(req *http.Request) (string, string, error) {
	bucket := req.URL.Query().Get("bucket")
	key := req.URL.Query().Get("key")
	if bucket == "" || key == "" {
		return "", "", fmt.Errorf("Missing 'bucket' or 'key' query parameter")
	}
	return bucket, key, nil
}

// saveFileToLocalPath saves the file to the local path
func saveFileToLocalPath(bucket, localPath, key string, fileContent []byte, w http.ResponseWriter) {
	filename := filepath.Base(key)
	completeLocalPath := filepath.Join(localPath, filename)
	err := os.WriteFile(completeLocalPath, fileContent, os.ModePerm)
	if err != nil {
		http.Error(w, "Failed to save the file to the local path", http.StatusInternalServerError)
		return
	}
	s3Path := fmt.Sprintf("s3://%s/%s", bucket, key)
	fmt.Fprintf(w, "File successfully downloaded from S3: \n%s \nand saved to: \n%s", s3Path, completeLocalPath)
}

// buildKeyWithOriginalFilename builds the key with the original filename
func buildKeyWithOriginalFilename(key, originalFilename string) string {
	return filepath.Join(key, originalFilename)
}

// getBucketAndKeyFromRequest gets the bucket and key from the request
func getBucketAndKeyFromRequest(req *http.Request) (string, string, error) {
	bucket := req.FormValue("bucket")
	key := req.FormValue("key")
	if bucket == "" || key == "" {
		return "", "", errors.New("missing required parameters")
	}
	return bucket, key, nil
}
