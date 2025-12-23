# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**go-lightning** is a lightweight Go library for simplified database operations using generics. It provides type-safe wrappers around `database/sql` with automatic struct-to-table mapping, reducing boilerplate for common CRUD operations.

## Repository Structure

This is a Go workspace with multiple modules:

- **Root module** (`github.com/tracewayapp/go-lightning`): Core generic database utilities
- **`pg/` submodule**: PostgreSQL-specific implementations with UUID support
- **`mysql/` submodule**: MySQL-specific implementations (placeholder)
- **`clickhouse/` submodule**: ClickHouse-specific implementations (placeholder)

The workspace is configured via `go.work` which includes all modules.

## Core Architecture

### Registration System

The library uses a registration pattern where struct types must be registered before use:

```go
Register[MyStruct](NamingStrategy)
```

Registration builds metadata (`FieldMap`) containing:
- Column name mappings (struct field → database column via naming strategy)
- Pre-built INSERT/UPDATE query templates
- ID field detection (for auto-increment handling)

Unregistered types will panic with a helpful error message pointing to the missing registration.

### Generic CRUD Operations

The library provides two API styles:

1. **Generic operations** (PostgreSQL module): `SelectGeneric`, `InsertGeneric`, `UpdateGeneric`
   - Automatically map struct fields to columns using reflection
   - No manual row scanning needed
   - Type must be registered first

2. **Manual operations** (Core module): `SelectMultiple`, `SelectSingle`, `Insert`, `Update`, `Delete`
   - Require explicit `mapLine` function for row scanning
   - More control but more boilerplate

### Naming Strategy

`DbNamingStrategy` interface converts Go naming to database naming:
- `GetTableNameFromStructName`: CamelCase → snake_case + pluralization (e.g., `UserProfile` → `user_profiles`)
- `GetColumnNameFromStructName`: CamelCase → snake_case (e.g., `FirstName` → `first_name`)

Default implementation: `DefaultDbNamingStrategy`

### PostgreSQL-Specific Features

The `pg/` module adds:
- **UUID support**: `InsertGenericUuid` auto-generates UUIDs for string ID fields
- **Existing UUID handling**: `InsertGenericExistingUuid` for pre-set UUIDs
- **Parameter offset handling**: `UpdateGeneric` automatically adjusts `$N` placeholders in WHERE clauses
- All generic operations work with PostgreSQL transactions (`*sql.Tx`)

## Development Commands

### Running Tests

```bash
# Run all tests in root module
go test ./...

# Run tests in PostgreSQL module
go test ./pg/...

# Run tests with coverage
go test -cover ./...

# Run specific test
go test -run TestRegister ./...
```

### Working with Go Workspace

```bash
# Sync workspace dependencies
go work sync

# Run commands in specific module
cd pg && go test
```

### Building

This is a library, not an executable. Import it in your projects:

```go
import "github.com/tracewayapp/go-lightning"
import lightning "github.com/tracewayapp/go-lightning/pg"
```

## Key Implementation Details

### Parameter Placeholders

- PostgreSQL uses `$1, $2, $3` placeholders (implemented in `pg/` module)
- Helper functions: `JoinStringForIn` and `JoinForIn` for dynamic IN clauses
- `UpdateGeneric` automatically renumbers placeholders when WHERE clause has `$N` references

### ID Field Handling

- Integer `Id` fields: Auto-detected and set to `DEFAULT` in INSERT, `RETURNING id` clause added
- UUID `Id` fields (string): Use `InsertGenericUuid` for auto-generation
- Other types: No special handling, must be provided by caller

### Reflection-Based Mapping

The `pg/` module uses `getPointersForColumns` to build scan destinations:
- Validates all returned columns exist in the registered struct
- Creates pointers to struct fields in correct order
- Panics on column mismatch (indicates schema/struct mismatch)

## Testing Approach

- Pure functions (naming strategy, string builders) have comprehensive unit tests
- Database operations (requiring `*sql.Tx`) use `t.Skip()` with guidance for `sqlmock` integration
- Test structs: `TestUser` (int ID), `TestProduct` (string ID) cover both ID scenarios
