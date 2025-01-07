package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "dbuser"
	password = "password"
	dbname   = "mydb"
)

var db *sql.DB

func main() {
	// Connect to the PostgreSQL database
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	var err error
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer db.Close()

	// Test the database connection
	err = db.Ping()
	if err != nil {
		log.Fatalf("Unable to ping database: %v", err)
	}

	// Define API routes
	http.HandleFunc("/tables", listTablesHandler)
	http.HandleFunc("/table-schemas", listTableSchemasHandler)

	// Start the HTTP server
	fmt.Println("Server is running on http://localhost:8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}

// TableInfo represents the schema and table name
type TableInfo struct {
	Schema string
	Table  string
}

// ColumnInfo represents the column details of a table
type ColumnInfo struct {
	Schema     string
	Table      string
	ColumnName string
	DataType   string
	IsNullable string
}

// listTablesHandler handles the request to list all tables and schemas
func listTablesHandler(w http.ResponseWriter, r *http.Request) {
	// Query to fetch all tables and their schemas
	query := `
		SELECT table_schema, table_name
		FROM information_schema.tables
		WHERE table_schema NOT IN ('pg_catalog', 'information_schema')
		ORDER BY table_schema, table_name;
	`

	rows, err := db.Query(query)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error querying database: %v", err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Store the results in a slice
	var tables []TableInfo
	for rows.Next() {
		var schema, table string
		if err := rows.Scan(&schema, &table); err != nil {
			http.Error(w, fmt.Sprintf("Error scanning row: %v", err), http.StatusInternalServerError)
			return
		}
		tables = append(tables, TableInfo{Schema: schema, Table: table})
	}

	// Check for errors after iteration
	if err := rows.Err(); err != nil {
		http.Error(w, fmt.Sprintf("Error after iterating rows: %v", err), http.StatusInternalServerError)
		return
	}

	// Render the HTML template
	renderHTML(w, tables, "tables")
}

// listTableSchemasHandler handles the request to list all table schemas with fields
func listTableSchemasHandler(w http.ResponseWriter, r *http.Request) {
	// Query to fetch all columns and their details
	query := `
		SELECT table_schema, table_name, column_name, data_type, is_nullable
		FROM information_schema.columns
		WHERE table_schema NOT IN ('pg_catalog', 'information_schema')
		ORDER BY table_schema, table_name, ordinal_position;
	`

	rows, err := db.Query(query)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error querying database: %v", err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Store the results in a slice
	var columns []ColumnInfo
	for rows.Next() {
		var schema, table, columnName, dataType, isNullable string
		if err := rows.Scan(&schema, &table, &columnName, &dataType, &isNullable); err != nil {
			http.Error(w, fmt.Sprintf("Error scanning row: %v", err), http.StatusInternalServerError)
			return
		}
		columns = append(columns, ColumnInfo{Schema: schema, Table: table, ColumnName: columnName, DataType: dataType, IsNullable: isNullable})
	}

	// Check for errors after iteration
	if err := rows.Err(); err != nil {
		http.Error(w, fmt.Sprintf("Error after iterating rows: %v", err), http.StatusInternalServerError)
		return
	}

	// Render the HTML template
	renderHTML(w, columns, "table-schemas")
}

// renderHTML renders the data in an HTML table
func renderHTML(w http.ResponseWriter, data interface{}, templateName string) {
	var htmlTemplate string
	switch templateName {
	case "tables":
		htmlTemplate = `
		<!DOCTYPE html>
		<html>
		<head>
			<title>Database Tables</title>
			<style>
				table {
					width: 50%;
					border-collapse: collapse;
					margin: 20px auto;
				}
				th, td {
					border: 1px solid #ddd;
					padding: 8px;
					text-align: left;
				}
				th {
					background-color: #f2f2f2;
				}
			</style>
		</head>
		<body>
			<h1 style="text-align: center;">Database Tables</h1>
			<table>
				<tr>
					<th>Schema</th>
					<th>Table</th>
				</tr>
				{{range .}}
				<tr>
					<td>{{.Schema}}</td>
					<td>{{.Table}}</td>
				</tr>
				{{end}}
			</table>
		</body>
		</html>
		`
	case "table-schemas":
		htmlTemplate = `
		<!DOCTYPE html>
		<html>
		<head>
			<title>Table Schemas</title>
			<style>
				table {
					width: 80%;
					border-collapse: collapse;
					margin: 20px auto;
				}
				th, td {
					border: 1px solid #ddd;
					padding: 8px;
					text-align: left;
				}
				th {
					background-color: #f2f2f2;
				}
			</style>
		</head>
		<body>
			<h1 style="text-align: center;">Table Schemas</h1>
			<table>
				<tr>
					<th>Schema</th>
					<th>Table</th>
					<th>Column</th>
					<th>Data Type</th>
					<th>Nullable</th>
				</tr>
				{{range .}}
				<tr>
					<td>{{.Schema}}</td>
					<td>{{.Table}}</td>
					<td>{{.ColumnName}}</td>
					<td>{{.DataType}}</td>
					<td>{{.IsNullable}}</td>
				</tr>
				{{end}}
			</table>
		</body>
		</html>
		`
	default:
		http.Error(w, "Invalid template name", http.StatusInternalServerError)
		return
	}

	// Parse and execute the template
	tmpl, err := template.New(templateName).Parse(htmlTemplate)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error parsing template: %v", err), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error executing template: %v", err), http.StatusInternalServerError)
		return
	}
}
