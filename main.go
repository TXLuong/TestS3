package main

import (
	"bytes"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"log"
	"net/http"
)

var (
	s3Client *s3.Client
	bucket   = "luongbuckettest"
	region   = "ap-southeast-1"
)

func main() {
	// Initialize AWS SDK configuration
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
	)
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	// Create an S3 service client
	s3Client = s3.NewFromConfig(cfg)

	// Define HTTP routes
	http.HandleFunc("/upload", handleUpload)

	// Start the server
	port := ":8082"
	fmt.Printf("Starting server on port %s...\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

// handleUpload handles the HTTP request to upload a file to S3
func handleUpload(w http.ResponseWriter, r *http.Request) {
	// Parse the form data
	err := r.ParseMultipartForm(10 << 20) // Max size 10MB
	if err != nil {
		http.Error(w, "Unable to parse form data", http.StatusBadRequest)
		return
	}

	// Prepare parameters for PutObject operation
	str := getJsonData()
	file := bytes.NewReader([]byte(str))
	params := &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(getObjectKey(r)),
		Body:   file,
	}

	// Upload the file to S3
	resp, err := s3Client.PutObject(context.TODO(), params)
	if err != nil {
		http.Error(w, fmt.Sprintf("Unable to upload file to S3: %v", err), http.StatusInternalServerError)
		return
	}

	// Return success response
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "File uploaded successfully to S3. ETag: %s", aws.ToString(resp.ETag))
}

func getJsonData() string {
	return "{\"order\": {\"id\": \"ZOiVezPL12iPbm1kjXkQ0RO6RjdZY\",\"location_id\": \"90JZWBT42N19G\",\"line_items\": [{\"uid\": \"513a0e77-f722-1b2e-23bd-6fd734b547c6\",\"quantity\": \"1\",\"base_price_money\": {\"amount\": 15000,\"currency\": \"USD\"},\"note\": \"Donation\",\"gross_sales_money\": {\"amount\": 15000,\"currency\": \"USD\"},\"total_tax_money\": {\"amount\": 0,\"currency\": \"USD\"},\"total_discount_money\": {\"amount\": 0,\"currency\": \"USD\"},\"total_money\": {\"amount\": 15000,\"currency\": \"USD\"},\"variation_total_price_money\": {\"amount\": 15000,\"currency\": \"USD\"},\"item_type\": \"CUSTOM_AMOUNT\",\"pricing_blocklists\": {},\"total_service_charge_money\": {\"amount\": 0,\"currency\": \"USD\"}}],\"created_at\": \"2024-06-07T11:58:46.533Z\",\"updated_at\": \"2024-06-07T11:58:48.000Z\",\"state\": \"COMPLETED\",\"version\": 4,\"total_tax_money\": {\"amount\": 0,\"currency\": \"USD\"},\"total_discount_money\": {\"amount\": 0,\"currency\": \"USD\"},\"total_tip_money\": {\"amount\": 0,\"currency\": \"USD\"},\"total_money\": {\"amount\": 15000,\"currency\": \"USD\"},\"closed_at\": \"2024-06-07T11:58:47.750Z\",\"tenders\": [{\"id\": \"haFECeXHmRy1LXDjJIBdC8nQlJDZY\",\"location_id\": \"90JZWBT42N19G\",\"transaction_id\": \"ZOiVezPL12iPbm1kjXkQ0RO6RjdZY\",\"created_at\": \"2024-06-07T11:58:47Z\",\"note\": \"Donation\",\"amount_money\": {\"amount\": 15000,\"currency\": \"USD\"},\"type\": \"OTHER\",\"payment_id\": \"haFECeXHmRy1LXDjJIBdC8nQlJDZY\"}],\"total_service_charge_money\": {\"amount\": 0,\"currency\": \"USD\"},\"net_amounts\": {\"total_money\": {\"amount\": 15000,\"currency\": \"USD\"},\"tax_money\": {\"amount\": 0,\"currency\": \"USD\"},\"discount_money\": {\"amount\": 0,\"currency\": \"USD\"},\"tip_money\": {\"amount\": 0,\"currency\": \"USD\"},\"service_charge_money\": {\"amount\": 0,\"currency\": \"USD\"}},\"source\": {},\"pricing_options\": {\"auto_apply_discounts\": true,\"auto_apply_taxes\": true},\"net_amount_due_money\": {\"amount\": 0,\"currency\": \"USD\"}}}"
}

func getObjectKey(r *http.Request) string {
	str := r.Header.Get("object-key")
	if str != "" {
		return str
	}
	return "object-key"
}
