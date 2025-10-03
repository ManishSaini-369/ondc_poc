
package registry

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"ondc-poc/internal/database"
	"ondc-poc/internal/models"
	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

// FindParticipants searches the registry for participants matching the given criteria.
func FindParticipants(query models.LookupRequest) ([]models.Participant, error) {
	// 1. Create a cache key from the LookupRequest.
	cacheKey, err := json.Marshal(query)
	if err != nil {
		return nil, fmt.Errorf("error creating cache key: %w", err)
	}

	// 2. Try to get the data from Redis.
	val, err := database.RDB.Get(ctx, string(cacheKey)).Result()
	if err == nil {
		var participants []models.Participant
		if err := json.Unmarshal([]byte(val), &participants); err == nil {
			// Data found in cache
			return participants, nil
		}
	}
	if err != redis.Nil {
		// Some other error with Redis
		return nil, fmt.Errorf("error getting from redis: %w", err)
	}

	// 3. If not in Redis, query the database.
	participants, err := findParticipantsInDB(query)
	if err != nil {
		return nil, fmt.Errorf("error getting from database: %w", err)
	}

	// 4. Cache the result in Redis.
	serialized, err := json.Marshal(participants)
	if err == nil {
		// Cache for 1 hour
		database.RDB.Set(ctx, string(cacheKey), serialized, time.Hour)
	}

	return participants, nil
}

func findParticipantsInDB(query models.LookupRequest) ([]models.Participant, error) {
	sql := `SELECT subscriber_id, status, ukid, subscriber_url, country, domain, 
                   valid_from, valid_until, type, signing_public_key, encr_public_key, 
                   created, updated, br_id, city 
            FROM ondc.participants WHERE 1=1`

	var args []interface{}
	argId := 1

	if query.SubscriberID != "" {
		sql += fmt.Sprintf(" AND subscriber_id = $%d", argId)
		args = append(args, query.SubscriberID)
		argId++
	}
	if query.Domain != "" {
		sql += fmt.Sprintf(" AND domain = $%d", argId)
		args = append(args, query.Domain)
		argId++
	}
	if query.UkID != "" {
		sql += fmt.Sprintf(" AND ukid = $%d", argId)
		args = append(args, query.UkID)
		argId++
	}
	if query.Country != "" {
		sql += fmt.Sprintf(" AND country = $%d", argId)
		args = append(args, query.Country)
		argId++
	}
	if query.City != "" {
		sql += fmt.Sprintf(" AND city = $%d", argId)
		args = append(args, query.City)
		argId++
	}
	if query.Type != "" {
		sql += fmt.Sprintf(" AND type = $%d", argId)
		args = append(args, query.Type)
		argId++
	}

	sqlDB, err := database.DB.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB: %w", err)
	}

	rows, err := sqlDB.Query(sql, args...)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	var participants []models.Participant
	for rows.Next() {
		var p models.Participant
		if err := rows.Scan(
			&p.SubscriberID,
			&p.Status,
			&p.UkID,
			&p.SubscriberURL,
			&p.Country,
			&p.Domain,
			&p.ValidFrom,
			&p.ValidUntil,
			&p.Type,
			&p.SigningPublicKey,
			&p.EncrPublicKey,
			&p.Created,
			&p.Updated,
			&p.BrID,
			&p.City,
		); err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}
		participants = append(participants, p)
	}

	return participants, nil
}

func GetParticipantBySubscriberID(subscriberID string) (*models.Participant, error) {
	query := models.LookupRequest{
		SubscriberID: subscriberID,
	}
	participants, err := findParticipantsInDB(query)
	if err != nil {
		return nil, err
	}
	if len(participants) == 0 {
		return nil, fmt.Errorf("participant with subscriber ID '%s' not found", subscriberID)
	}
	return &participants[0], nil
}
