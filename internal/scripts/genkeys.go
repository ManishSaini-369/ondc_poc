// run: go run internal/scripts/genkeys.go
package main

import (
	"crypto/ed25519"
	"database/sql"
	"encoding/base64"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	// ⚡ Static DB connection string
	connStr := "host=localhost port=5432 user=postgres password=admin dbname=registry_db sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("❌ Failed to connect DB:", err)
	}
	defer db.Close()

	rows, err := db.Query(`SELECT ukid FROM participants`)
	if err != nil {
		log.Fatal("❌ Failed to fetch UKIDs:", err)
	}
	defer rows.Close()

	// create a dir for private keys
	_ = os.MkdirAll("keys", 0755)

	for rows.Next() {
		var ukid string
		if err := rows.Scan(&ukid); err != nil {
			log.Println("❌ Scan error:", err)
			continue
		}

		// generate keypair
		pub, priv, _ := ed25519.GenerateKey(nil)
		pubB64 := base64.StdEncoding.EncodeToString(pub)
		privB64 := base64.StdEncoding.EncodeToString(priv)

		// update table
		_, err = db.Exec(`UPDATE participants SET signing_public_key=$1 WHERE ukid=$2`, pubB64, ukid)
		if err != nil {
			log.Println("❌ Failed to update for ukid:", ukid, err)
			continue
		}

		// save private key in file
		privPath := fmt.Sprintf("keys/%s_priv.key", ukid)
		err = os.WriteFile(privPath, []byte(privB64), 0600)
		if err != nil {
			log.Println("❌ Failed to write private key for ukid:", ukid, err)
			continue
		}

		log.Println("✅ Updated UKID:", ukid)
	}

	log.Println("🎉 Done updating all participants")
}
