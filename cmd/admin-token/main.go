package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"open-transit-rt/internal/auth"
)

func main() {
	subject := flag.String("sub", "", "admin JWT subject")
	agencyID := flag.String("agency-id", "", "admin JWT agency_id")
	ttlRaw := flag.String("ttl", "", "token lifetime; defaults to ADMIN_JWT_TTL or 8h")
	flag.Parse()

	cfg, err := auth.JWTConfigFromEnv()
	if err != nil {
		log.Fatal(err)
	}
	ttl := cfg.TTL
	if *ttlRaw != "" {
		parsed, err := time.ParseDuration(*ttlRaw)
		if err != nil {
			log.Fatal(err)
		}
		ttl = parsed
	}
	signer, err := auth.NewSigner(cfg)
	if err != nil {
		log.Fatal(err)
	}
	token, claims, err := signer.Sign(*subject, *agencyID, ttl)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("token=%s\n", token)
	fmt.Printf("sub=%s\nagency_id=%s\nexpires_at=%s\n", claims.Subject, claims.AgencyID, time.Unix(claims.Expires, 0).UTC().Format(time.RFC3339))
}
