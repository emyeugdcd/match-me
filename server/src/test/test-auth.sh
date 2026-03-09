#!/bin/bash

# Test script to verify authentication endpoints
BASE_URL="http://localhost:8080"

echo "Testing Authentication Endpoints..."
echo "=================================="

# Test 1: Signup
echo -e "\n1. Testing signup..."
SIGNUP_RESPONSE=$(curl -s -X POST "$BASE_URL/signup" \
  -H "Content-Type: application/json" \
  -d '{"email": "test@example.com", "password": "testpassword123"}' \
  -w "\nHTTP_CODE:%{http_code}")

echo "Signup Response:"
echo "$SIGNUP_RESPONSE"

# Extract token from signup response
TOKEN=$(echo "$SIGNUP_RESPONSE" | grep -o '"token":"[^"]*' | cut -d'"' -f4)
echo "Extracted Token: ${TOKEN:0:50}..."

# Test 2: Login with same credentials
echo -e "\n2. Testing login..."
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/login" \
  -H "Content-Type: application/json" \
  -d '{"email": "test@example.com", "password": "testpassword123"}' \
  -w "\nHTTP_CODE:%{http_code}")

echo "Login Response:"
echo "$LOGIN_RESPONSE"

# Test 3: Access protected endpoint with token
echo -e "\n3. Testing protected endpoint with token..."
PROTECTED_RESPONSE=$(curl -s -X GET "$BASE_URL/me" \
  -H "Authorization: Bearer $TOKEN" \
  -w "\nHTTP_CODE:%{http_code}")

echo "Protected endpoint response:"
echo "$PROTECTED_RESPONSE"

# Test 4: Access protected endpoint without token
echo -e "\n4. Testing protected endpoint without token..."
UNAUTH_RESPONSE=$(curl -s -X GET "$BASE_URL/me" \
  -w "\nHTTP_CODE:%{http_code}")

echo "Unauthorized response:"
echo "$UNAUTH_RESPONSE"

# Test 5: Logout
echo -e "\n5. Testing logout..."
LOGOUT_RESPONSE=$(curl -s -X POST "$BASE_URL/logout" \
  -H "Authorization: Bearer $TOKEN" \
  -w "\nHTTP_CODE:%{http_code}")

echo "Logout response:"
echo "$LOGOUT_RESPONSE"

echo -e "\n=================================="
echo "Authentication tests completed!"