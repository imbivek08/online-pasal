#!/bin/bash

# Migration helper script for goose

MIGRATION_DIR="apps/backend/internal/database/migratons"

case "$1" in
  create)
    if [ -z "$2" ]; then
      echo "Usage: ./migrate.sh create <migration_name>"
      exit 1
    fi
    goose -dir "$MIGRATION_DIR" create "$2" sql
    echo "Migration created in $MIGRATION_DIR"
    ;;
  
  up)
    echo "Run migrations through the application using 'go run cmd/nepify/main.go migrate up'"
    echo "Or manually: goose -dir $MIGRATION_DIR postgres \$DATABASE_URL up"
    ;;
  
  down)
    echo "Rollback through the application using 'go run cmd/nepify/main.go migrate down'"
    echo "Or manually: goose -dir $MIGRATION_DIR postgres \$DATABASE_URL down"
    ;;
  
  status)
    echo "Check migration status through the application using 'go run cmd/nepify/main.go migrate status'"
    echo "Or manually: goose -dir $MIGRATION_DIR postgres \$DATABASE_URL status"
    ;;
  
  *)
    echo "Usage: ./migrate.sh {create|up|down|status} [args]"
    echo ""
    echo "Commands:"
    echo "  create <name>  - Create a new migration file"
    echo "  up            - Run all pending migrations"
    echo "  down          - Rollback the last migration"
    echo "  status        - Show migration status"
    exit 1
    ;;
esac
