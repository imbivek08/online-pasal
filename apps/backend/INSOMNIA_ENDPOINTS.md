# Nepify API Endpoints - Insomnia Collection

## Base URL
```
http://localhost:8080
```

## All Available Endpoints

### 1. Health Check
**GET** `/health`
- **Description**: Check if the server is running
- **Authentication**: None
- **Response**: 200 OK
```json
{
  "status": "ok",
  "message": "server is running"
}
```

---

### 2. Clerk Webhook
**POST** `/api/v1/webhooks/clerk`
- **Description**: Clerk authentication webhook for user sync
- **Authentication**: Clerk webhook signature
- **Headers**:
  - `svix-id`: Webhook message ID
  - `svix-timestamp`: Webhook timestamp
  - `svix-signature`: Webhook signature
- **Body**: Clerk webhook payload
- **Note**: This is automatically called by Clerk, not for manual testing

---

## User Endpoints (Protected - Requires Authentication)

### 3. Get User Profile
**GET** `/api/v1/users/profile`
- **Description**: Get the authenticated user's profile
- **Authentication**: Required (Bearer Token)
- **Headers**:
  ```
  Authorization: Bearer YOUR_CLERK_TOKEN
  ```
- **Response**: 200 OK
```json
{
  "data": {
    "id": "uuid",
    "clerk_id": "user_xxx",
    "email": "user@example.com",
    "username": "username",
    "first_name": "John",
    "last_name": "Doe",
    "phone": "+977-9841234567",
    "role": "customer",
    "is_active": true,
    "created_at": "2026-01-17T10:00:00Z",
    "updated_at": "2026-01-17T10:00:00Z"
  }
}
```

### 4. Update User Profile
**PUT** `/api/v1/users/profile`
- **Description**: Update the authenticated user's profile
- **Authentication**: Required (Bearer Token)
- **Headers**:
  ```
  Authorization: Bearer YOUR_CLERK_TOKEN
  Content-Type: application/json
  ```
- **Body**:
```json
{
  "username": "newusername",
  "first_name": "John",
  "last_name": "Doe",
  "phone": "+977-9841234567"
}
```
- **Response**: 200 OK

### 5. Delete User Account
**DELETE** `/api/v1/users/account`
- **Description**: Delete the authenticated user's account
- **Authentication**: Required (Bearer Token)
- **Headers**:
  ```
  Authorization: Bearer YOUR_CLERK_TOKEN
  ```
- **Response**: 200 OK

### 6. Get User By ID
**GET** `/api/v1/users/:id`
- **Description**: Get a specific user by ID
- **Authentication**: Required (Bearer Token)
- **Headers**:
  ```
  Authorization: Bearer YOUR_CLERK_TOKEN
  ```
- **Example**: `/api/v1/users/123e4567-e89b-12d3-a456-426614174000`
- **Response**: 200 OK

---

## Product Endpoints

### 7. Get All Products (Public)
**GET** `/api/v1/products`
- **Description**: Get all products with optional filtering
- **Authentication**: None
- **Query Parameters**:
  - `search` (optional): Search by product name or description
    - Example: `/api/v1/products?search=laptop`
  - `category` (optional): Filter by category slug
    - Example: `/api/v1/products?category=electronics`
  - `shop_id` (optional): Filter by shop ID
  - `min_price` (optional): Minimum price filter
  - `max_price` (optional): Maximum price filter
  - `is_featured` (optional): Filter featured products (true/false)
  - `page` (optional): Page number (default: 1)
  - `limit` (optional): Items per page (default: 20)
- **Response**: 200 OK
```json
{
  "data": [
    {
      "id": "uuid",
      "shop_id": "uuid",
      "name": "Dell XPS 15 Laptop",
      "slug": "dell-xps-15-laptop",
      "description": "High-performance laptop...",
      "short_description": "Premium laptop for professionals",
      "sku": "DELL-XPS15-001",
      "price": 189999.00,
      "compare_at_price": 209999.00,
      "stock_quantity": 15,
      "is_active": true,
      "is_featured": true,
      "created_at": "2026-01-17T10:00:00Z"
    }
  ]
}
```

### 8. Get Product By ID (Public)
**GET** `/api/v1/products/:id`
- **Description**: Get a single product by ID
- **Authentication**: None
- **Example**: `/api/v1/products/123e4567-e89b-12d3-a456-426614174000`
- **Response**: 200 OK

### 9. Create Product (Protected - Vendor Only)
**POST** `/api/v1/products`
- **Description**: Create a new product (vendors only)
- **Authentication**: Required (Bearer Token - must be vendor role)
- **Headers**:
  ```
  Authorization: Bearer YOUR_CLERK_TOKEN
  Content-Type: application/json
  ```
- **Body**:
```json
{
  "name": "New Product Name",
  "description": "Detailed product description",
  "short_description": "Brief description",
  "sku": "PROD-001",
  "price": 1999.99,
  "compare_at_price": 2499.99,
  "stock_quantity": 50,
  "category_id": "category-uuid",
  "weight": 0.5,
  "weight_unit": "kg",
  "is_active": true,
  "is_featured": false
}
```
- **Response**: 201 Created

### 10. Update Product (Protected - Vendor Only)
**PUT** `/api/v1/products/:id`
- **Description**: Update an existing product (only owner can update)
- **Authentication**: Required (Bearer Token - must be vendor/owner)
- **Headers**:
  ```
  Authorization: Bearer YOUR_CLERK_TOKEN
  Content-Type: application/json
  ```
- **Example**: `/api/v1/products/123e4567-e89b-12d3-a456-426614174000`
- **Body**:
```json
{
  "name": "Updated Product Name",
  "price": 1899.99,
  "stock_quantity": 45,
  "is_active": true
}
```
- **Response**: 200 OK

### 11. Delete Product (Protected - Vendor Only)
**DELETE** `/api/v1/products/:id`
- **Description**: Delete a product (only owner can delete)
- **Authentication**: Required (Bearer Token - must be vendor/owner)
- **Headers**:
  ```
  Authorization: Bearer YOUR_CLERK_TOKEN
  ```
- **Example**: `/api/v1/products/123e4567-e89b-12d3-a456-426614174000`
- **Response**: 200 OK

### 12. Get Vendor's Products (Protected - Vendor Only)
**GET** `/api/v1/vendor/products`
- **Description**: Get all products belonging to the authenticated vendor
- **Authentication**: Required (Bearer Token - must be vendor role)
- **Headers**:
  ```
  Authorization: Bearer YOUR_CLERK_TOKEN
  ```
- **Response**: 200 OK
```json
{
  "data": [
    {
      "id": "uuid",
      "name": "My Product 1",
      "price": 1999.99,
      "stock_quantity": 50
    }
  ]
}
```

---

## Testing with Seed Data

You can test the following with seeded data:

### Test Accounts (Mock Clerk IDs)
- **Customer**: `user_test_customer1` (john.doe@example.com)
- **Vendor 1**: `user_test_vendor1` (vendor1@nepify.com) - Tech Store Nepal
- **Vendor 2**: `user_test_vendor2` (vendor2@nepify.com) - Fashion Hub
- **Vendor 3**: `user_test_vendor3` (vendor3@nepify.com) - Book World

### Sample Product Searches
```
GET /api/v1/products?search=laptop
GET /api/v1/products?search=book
GET /api/v1/products?category=electronics
GET /api/v1/products?is_featured=true
GET /api/v1/products?min_price=1000&max_price=5000
```

---

## Authentication Notes

### Getting a Clerk Token
1. Sign in through your frontend application
2. Use browser DevTools to capture the token from:
   - Local storage
   - Network requests (Authorization header)
3. Copy the token (starts with `ey...`)
4. Use it in Insomnia as: `Bearer YOUR_TOKEN_HERE`

### Testing Without Real Clerk Tokens
Since seed data uses mock `clerk_id` values, you'll need real Clerk authentication for protected routes. Alternatively, you can:
1. Disable auth middleware temporarily for testing
2. Mock the auth middleware in development
3. Use actual Clerk sign-in flow to get real tokens

---

## Common Response Codes

- **200 OK**: Request successful
- **201 Created**: Resource created successfully
- **400 Bad Request**: Invalid request data
- **401 Unauthorized**: Missing or invalid authentication
- **403 Forbidden**: Insufficient permissions
- **404 Not Found**: Resource not found
- **500 Internal Server Error**: Server error

---

## Example Insomnia Requests

### 1. Get All Products
```
Method: GET
URL: http://localhost:8080/api/v1/products
```

### 2. Search Products
```
Method: GET
URL: http://localhost:8080/api/v1/products?search=laptop&category=electronics
```

### 3. Get Product by ID
```
Method: GET
URL: http://localhost:8080/api/v1/products/{product_id}
```

### 4. Create Product (Vendor)
```
Method: POST
URL: http://localhost:8080/api/v1/products
Headers:
  Authorization: Bearer {your_clerk_token}
  Content-Type: application/json
Body (JSON):
{
  "name": "Test Product",
  "description": "This is a test product",
  "price": 999.99,
  "stock_quantity": 10
}
```

### 5. Get Vendor Products
```
Method: GET
URL: http://localhost:8080/api/v1/vendor/products
Headers:
  Authorization: Bearer {your_clerk_token}
```

### 6. Get User Profile
```
Method: GET
URL: http://localhost:8080/api/v1/users/profile
Headers:
  Authorization: Bearer {your_clerk_token}
```

---

## Environment Variables for Insomnia

Create an environment in Insomnia with these variables:

```json
{
  "base_url": "http://localhost:8080",
  "api_v1": "http://localhost:8080/api/v1",
  "clerk_token": "YOUR_CLERK_TOKEN_HERE"
}
```

Then use them in requests:
- URL: `{{ _.base_url }}/api/v1/products`
- Header: `Bearer {{ _.clerk_token }}`
