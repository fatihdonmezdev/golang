package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	_ "github.com/lib/pq"
)

type URI struct {
	ID       int    `json:"id"`
	ServerIP string `json:"server_ip"`
	URI      string `json:"uri"`
	Resnum   string `json:"res_num"`
	Resphead  string `json:"res_head"`
}

type Conversation struct {
	Name  string `json:"name"`
	IP    string `json:"ip"`	
	URIs  []URI  `json:"uris"`
}


func main() {
	// Read JSON data from file
	jsonData, err := ioutil.ReadFile("./tutorial/data.json")
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	// Define a struct to hold the conversations
	var data struct {
		Conversations []Conversation `json:"conversations"`
	}

	// Unmarshal JSON data into the struct
	if err := json.Unmarshal(jsonData, &data); err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}
	db, err := sql.Open("postgres", "postgres://postgres:password1@localhost/mydb?sslmode=disable")
	if err != nil {
		log.Fatal("Error connecting to the database:", err)
	}
	defer db.Close()
	// Create the uris table if it doesn't exist
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS conversations (
		id SERIAL PRIMARY KEY,
		ip VARCHAR(255) ,
		name VARCHAR(255)
	)`)
	if err != nil {
		log.Fatal("Error creating conversations table:", err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS uris (
		id SERIAL PRIMARY KEY,
		conversation_id INT,
		server_ip VARCHAR(255),
		uri VARCHAR(255),
		res_num VARCHAR(255),
		res_date VARCHAR(255)
	)`)
	if err != nil {
		log.Fatal("Error creating uris table:", err)
	}

	for _, conv := range data.Conversations {
    _, err = db.Exec("INSERT INTO conversations (name, ip) VALUES ($1, $2)", conv.Name, conv.IP)
    if err != nil {
        log.Fatal("Error inserting conversation data:", err)
    }

    var convID int
    err := db.QueryRow("SELECT id FROM conversations ORDER BY id DESC LIMIT 1").Scan(&convID)
    if err != nil {
        log.Fatal("Error querying conversation ID:", err)
    }
	
	for _, uri := range conv.URIs {
		// Extract date from res_head
		resDate := extractDateFromResHead(uri.Resphead)
		println(resDate) // Print the extracted date
		
		// Insert URI data into the database
		_, err = db.Exec("INSERT INTO uris (conversation_id, server_ip, uri, res_num, res_date) VALUES ($1, $2, $3, $4, $5)",
			convID, uri.ServerIP, uri.URI, uri.Resnum, resDate)
		if err != nil {
			log.Fatal("Error inserting URI data:", err)
		}
	}
	
	}


		fmt.Println("Data inserted successfully!")

}