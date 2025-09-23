package main

import (
    "context"
    "crypto/ed25519"
    "database/sql"
    "encoding/base64"
    "encoding/json"
    "fmt"
    "log"
    "os"

    _ "github.com/lib/pq"
)

type Participant struct {
    ID string
}

func main() {
    host := "localhost"
    port := 5432
    user := "postgres"
    password := "admin"
    dbname := "registry_db"

    psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
        host, port, user, password, dbname)

    db, err := sql.Open("postgres", psqlInfo)
    if err != nil {
        log.Fatalf("Error opening database: %v\n", err)
    }
    defer db.Close()

    ctx := context.Background()

    rows, err := db.QueryContext(ctx, "SELECT id FROM participants")
    if err != nil {
        log.Fatalf("Error querying participants: %v\n", err)
    }
    defer rows.Close()

    participants := []Participant{}
    for rows.Next() {
        var p Participant
        if err := rows.Scan(&p.ID); err != nil {
            log.Fatalf("Error scanning participant id: %v\n", err)
        }
        participants = append(participants, p)
    }
    if err = rows.Err(); err != nil {
        log.Fatalf("Row error: %v\n", err)
    }

    updateStmt, err := db.PrepareContext(ctx, "UPDATE participants SET signing_public_key=$1 WHERE id=$2")
    if err != nil {
        log.Fatalf("Error preparing update: %v\n", err)
    }
    defer updateStmt.Close()

    keyMap := make(map[string]map[string]string)

    for _, p := range participants {
        pubKey, privKey, err := ed25519.GenerateKey(nil)
        if err != nil {
            log.Printf("Key generation error for %s: %v", p.ID, err)
            continue
        }

        // Store raw 32-byte public key
        pubBase64 := base64.StdEncoding.EncodeToString(pubKey)
        privBase64 := base64.StdEncoding.EncodeToString(privKey)

        _, err = updateStmt.ExecContext(ctx, pubBase64, p.ID)
        if err != nil {
            log.Printf("Error updating public key for %s: %v", p.ID, err)
            continue
        }

        keyMap[p.ID] = map[string]string{
            "public":  pubBase64,
            "private": privBase64,
        }

        fmt.Printf("✅ Updated participant %s with raw Ed25519 keys\n", p.ID)
    }

    file, err := os.Create("participant_keys.json")
    if err != nil {
        log.Fatalf("Error creating JSON file: %v", err)
    }
    defer file.Close()

    encoder := json.NewEncoder(file)
    encoder.SetIndent("", "  ")
    if err := encoder.Encode(keyMap); err != nil {
        log.Fatalf("Error writing JSON file: %v", err)
    }

    fmt.Println("✅ All key pairs saved to participant_keys.json")
}
