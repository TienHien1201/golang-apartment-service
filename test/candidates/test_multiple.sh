#!/bin/bash

# API endpoint
API_URL="http://localhost:8080/api/v1/candidates/multiple"

# Test case 1: Create multiple candidates without files
echo "Test case 1: Create multiple candidates without files"
curl -X POST \
  -H "Content-Type: multipart/form-data" \
  -F "candidates=@multiple_candidates.json" \
  $API_URL

echo -e "\n\n"

# Test case 2: Create multiple candidates with files
echo "Test case 2: Create multiple candidates with files"
curl -X POST \
  -H "Content-Type: multipart/form-data" \
  -F "candidates=@multiple_candidates.json" \
  -F "uploads=@test_resume1.pdf" \
  -F "uploads=@test_resume2.pdf" \
  -F "uploads=@test_resume3.pdf" \
  $API_URL

echo -e "\n\n"

# Test case 3: Create multiple candidates with compressed data
echo "Test case 3: Create multiple candidates with compressed data"
curl -X POST \
  -H "Content-Type: multipart/form-data" \
  -F "candidates=@multiple_candidates.json" \
  -F "compressed=true" \
  $API_URL 