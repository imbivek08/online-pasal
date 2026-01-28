#!/bin/bash

# Seed Database Script for Nepify
# This script loads seed data into the PostgreSQL database

set -e

# Load environment variables
if [ -f .env ]; then
    export $(cat .env | grep -v '^#' | xargs)
fi

# Database connection details
DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-5432}"
DB_USER="${DB_USER:-${USERNAME:-nepify}}"
DB_PASSWORD="${DB_PASSWORD:-${PASSWORD:-nepify_password}}"
DB_NAME="${DB_NAME:-${DB_NAME:-nepify_db}}"

echo "============================================"
echo "Nepify Database Seeding"
echo "============================================"
echo "Database: $DB_NAME"
echo "Host: $DB_HOST:$DB_PORT"
echo "User: $DB_USER"
echo "============================================"
echo ""

# Check if psql is installed
if ! command -v psql &> /dev/null; then
    echo "Error: psql is not installed. Please install PostgreSQL client."
    exit 1
fi

# Run the seed script
echo "Loading seed data..."
PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -f internal/database/seed_data.sql

if [ $? -eq 0 ]; then
    echo ""
    echo "============================================"
    echo "✅ Seed data loaded successfully!"
    echo "============================================"
    echo ""
    echo "Test Accounts:"
    echo "  Customer 1: john.doe@example.com"
    echo "  Customer 2: jane.smith@example.com"
    echo "  Vendor 1: vendor1@nepify.com (Tech Store Nepal)"
    echo "  Vendor 2: vendor2@nepify.com (Fashion Hub)"
    echo "  Vendor 3: vendor3@nepify.com (Book World)"
    echo "  Admin: admin@nepify.com"
    echo ""
    echo "Note: These are test accounts with mock clerk_ids"
    echo "============================================"
else
    echo ""
    echo "❌ Failed to load seed data"
    exit 1
fi
