# Nepify

Modern e-commerce platform built with Go and React.

## Tech Stack

- **Backend**: Go, PostgreSQL
- **Frontend**: React, TypeScript, Tailwind CSS
- **Auth**: Clerk

## Getting Started

### Backend
```bash
cd apps/backend
make run
```

### Frontend
```bash
cd apps/frontend
npm install
npm run dev
```

### Docker
```bash
docker-compose up
```

## Features

- 🛒 Product browsing and cart management
- 👤 User authentication and profiles
- 🏪 Vendor shops and product management
- 📦 Order processing
- 💳 Checkout system

## Environment Variables

Create `.env` files in backend and frontend directories. See respective folders for required variables.

## API

Backend runs on `http://localhost:8080`  
Frontend runs on `http://localhost:5173`
