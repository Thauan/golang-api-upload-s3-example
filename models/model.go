package models

import (
	"database/sql"
)

// Create an exported global variable to hold the database connection pool.
var DB *sql.DB
