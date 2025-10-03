package main

import (
	"bufio"
	"bytes"
	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/blake2b"
	"log"
	"strings"
)

func NormalizeJSON(body []byte) ([]byte, error) {
	var buf bytes.Buffer
	if err := json.Compact(&buf, body); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func CreateDigest(body []byte) (string, error) {
	// Step 1: Normalize JSON
	normalized, err := NormalizeJSON(body)
	if err != nil {
		return "", err
	}

	// Step 2: Generate BLAKE2b-256 hash
	hash := blake2b.Sum256(normalized)

	// Step 3: Return Base64 encoded hash
	return base64.StdEncoding.EncodeToString(hash[:]), nil
}

// GenerateKeys creates a new ed25519 key pair.
func GenerateKeys() (ed25519.PublicKey, ed25519.PrivateKey, error) {
	return ed25519.GenerateKey(nil)
}

type NetworkParticipant struct {
	Country      string `json:"country"`
	Domain       string `json:"domain"`
	Type         string `json:"type"`
	SubscriberID string `json:"subscriber_id"`
}

func main() {
	data := "\"IND\"\t\"ONDC:FOD01\"\t\"buyerApp\"\t\"std:001\"\t\"idfcfirstbank.com\"
\"IND\"\t\"ONDC:FOD01\"\t\"sellerApp\"\t\"std:002\"\t\"icicibank.com\"
\"IND\"\t\"ONDC:FOD01\"\t\"buyerApp\"\t\"std:003\"\t\"hdfcbank.com\"
\"IND\"\t\"ONDC:FOD01\"\t\"sellerApp\"\t\"std:004\"\t\"sbinationalbank.com\"
\"IND\"\t\"ONDC:FOD01\"\t\"buyerApp\"\t\"std:005\"\t\"axisbank.com\"
\"IND\"\t\"ONDC:FOD01\"\t\"sellerApp\"\t\"std:006\"\t\"paytm.com\"
\"IND\"\t\"ONDC:FOD01\"\t\"buyerApp\"\t\"std:007\"\t\"phonepe.com\"
\"IND\"\t\"ONDC:FOD01\"\t\"sellerApp\"\t\"std:008\"\t\"amazonpay.com\"
\"IND\"\t\"ONDC:FOD01\"\t\"buyerApp\"\t\"std:009\"\t\"googlepay.com\"
\"IND\"\t\"ONDC:FOD01\"\t\"sellerApp\"\t\"std:010\"\t\"mobikwik.com\"
\"IND\"\t\"ONDC:FOD01\"\t\"buyerApp\"\t\"std:011\"\t\"reliancejio.com\"
\"IND\"\t\"ONDC:FOD01\"\t\"sellerApp\"\t\"std:012\"\t\"tata.com\"
\"IND\"\t\"ONDC:FOD01\"\t\"buyerApp\"\t\"std:013\"\t\"infosys.com\"
\"IND\"\t\"ONDC:FOD01\"\t\"sellerApp\"\t\"std:014\"\t\"wipro.com\"
\"IND\"\t\"ONDC:FOD01\"\t\"buyerApp\"\t\"std:015\"\t\"mahindraltd.com\"
\"IND\"\t\"ONDC:FOD01\"\t\"sellerApp\"\t\"std:016\"\t\"airtel.in\"
\"IND\"\t\"ONDC:FOD01\"\t\"buyerApp\"\t\"std:017\"\t\"vodafone.in\"
\"IND\"\t\"ONDC:FOD01\"\t\"sellerApp\"\t\"std:018\"\t\"flipkart.com\"
\"IND\"\t\"ONDC:FOD01\"\t\"buyerApp\"\t\"std:019\"\t\"snapdeal.com\"
\"IND\"\t\"ONDC:FOD01\"\t\"sellerApp\"\t\"std:020\"\t\"myntra.com\"
\"IND\"\t\"ONDC:FOD01\"\t\"buyerApp\"\t\"std:021\"\t\"ola.com\"
\"IND\"\t\"ONDC:FOD01\"\t\"sellerApp\"\t\"std:022\"\t\"uber.com\"
\"IND\"\t\"ONDC:FOD01\"\t\"buyerApp\"\t\"std:023\"\t\"zomato.com\"
\"IND\"\t\"ONDC:FOD01\"\t\"sellerApp\"\t\"std:024\"\t\"swiggy.com\"
\"IND\"\t\"ONDC:FOD01\"\t\"buyerApp\"\t\"std:025\"\t\"bigbasket.com\"
\"IND\"\t\"ONDC:FOD01\"\t\"sellerApp\"\t\"std:026\"\t\"dhl.com\"
\"IND\"\t\"ONDC:FOD01\"\t\"buyerApp\"\t\"std:027\"\t\"fedex.com\"
\"IND\"\t\"ONDC:FOD01\"\t\"sellerApp\"\t\"std:028\"\t\"ups.com\"
\"IND\"\t\"ONDC:FOD01\"\t\"buyerApp\"\t\"std:029\"\t\"bluedart.com\"
\"IND\"\t\"ONDC:FOD01\"\t\"sellerApp\"\t\"std:030\"\t\"delhivery.com\"
\"IND\"\t\"ONDC:FOD01\"\t\"buyerApp\"\t\"std:031\"\t\"makeinindia.com\"
\"IND\"\t\"ONDC:FOD01\"\t\"sellerApp\"\t\"std:032\"\t\"startupindia.gov.in\"
\"IND\"\t\"ONDC:FOD01\"\t\"buyerApp\"\t\"std:033\"\t\"nasscom.in\"
\"IND\"\t\"ONDC:FOD01\"\t\"sellerApp\"\t\"std:034\"\t\"cognizant.com\"
\"IND\"\t\"ONDC:FOD01\"\t\"buyerApp\"\t\"std:035\"\t\"tcs.com\"
\"IND\"\t\"ONDC:FOD01\"\t\"sellerApp\"\t\"std:036\"\t\"larsentoubro.com\"
\"IND\"\t\"ONDC:FOD01\"\t\"buyerApp\"\t\"std:037\"\t\"drreddys.com\"
\"IND\"\t\"ONDC:FOD01\"\t\"sellerApp\"\t\"std:038\"\t\"biocon.com\"
\"IND\"\t\"ONDC:FOD01\"\t\"buyerApp\"\t\"std:039\"\t\"sunpharma.com\"
\"IND\"\t\"ONDC:FOD01\"\t\"sellerApp\"\t\"std:040\"\t\"cipla.com\"
\"IND\"\t\"ONDC:FOD01\"\t\"buyerApp\"\t\"std:041\"\t\"bajajauto.com\"
\"IND\"\t\"ONDC:FOD01\"\t\"sellerApp\"\t\"std:042\"\t\"hero.com\"
\"IND\"\t\"ONDC:FOD01\"\t\"buyerApp\"\t\"std:043\"\t\"tvs.com\"
\"IND\"\t\"ONDC:FOD01\"\t\"sellerApp\"\t\"std:044\"\t\"royalenfield.com\"
\"IND\"\t\"ONDC:FOD01\"\t\"buyerApp\"\t\"std:045\"\t\"marutisuzuki.com\"
\"IND\"\t\"ONDC:FOD01\"\t\"sellerApp\"\t\"std:046\"\t\"nestle.in\"
\"IND\"\t\"ONDC:FOD01\"\t\"buyerApp\"\t\"std:047\"\t\"pepsico.com\"
\"IND\"\t\"ONDC:FOD01\"\t\"sellerApp\"\t\"std:048\"\t\"itcportal.com\"
\"IND\"\t\"ONDC:FOD01\"\t\"buyerApp\"\t\"std:049\"\t\"britannia.co.in\"
\"IND\"\t\"ONDC:FOD01\"\t\"sellerApp\"\t\"std:050\"\t\"godrej.com\"
\"IND\"\t\"ONDC:FOD01\"\t\"buyerApp\"\t\"std:011\"\t\"idfcfirstbank.com\"
\"IND\"\t\"ONDC:FOD01\"\t\"sellerApp\"\t\"std:033\"\t\"icicibank.com\"
\"IND\"\t\"ONDC:TRV01\"\t\"buyerApp\"\t\"std:0755\"\t\"acc.in\"
\"IND\"\t\"ONDC:FOD01\"\t\"gateway\"\t\"std:080\"\t\"kotak.com\"
\"IND\"\t\"ONDC:MED01\"\t\"gateway\"\t\"std:011\"\t\"phonepe.com\"
\"IND\"\t\"ONDC:TRV01\"\t\"sellerApp\"\t\"std:040\"\t\"acer.in\"
\"IND\"\t\"ONDC:RET10\"\t\"sellerApp\"\t\"std:080\"\t\"fbb.com\"
\"IND\"\t\"ONDC:RET10\"\t\"buyerApp\"\t\"std:044\"\t\"idfcfirstbank.com\"
\"IND\"\t\"ONDC:FOD01\"\t\"buyerApp\"\t\"std:0755\"\t\"icicibank.com\"
\"IND\"\t\"ONDC:RET10\"\t\"buyerApp\"\t\"std:033\"\t\"phonepe.com\"
\"IND\"\t\"ONDC:FOD01\"\t\"buyerApp\"\t\"std:011\"\t\"fbb.com\"
\"IND\"\t\"ONDC:FOD01\"\t\"buyerApp\"\t\"std:044\"\t\"nykaa.com\"
\"IND\"\t\"ONDC:FOD01\"\t\"sellerApp\"\t\"std:044\"\t\"paytm.com\"
\"IND\"\t\"ONDC:TRV01\"\t\"sellerApp\"\t\"std:044\"\t\"paytm.com\"
\"IND\"\t\"ONDC:FOD01\"\t\"buyerApp\"\t\"std:033\"\t\"tatadigital.com\"
\"IND\"\t\"ONDC:FOD01\"\t\"buyerApp\"\t\"std:022\"\t\"ajio.com\"
\"IND\"\t\"ONDC:TRV01\"\t\"gateway\"\t\"std:0755\"\t\"axisbank.com\"
\"IND\"\t\"ONDC:LOG12\"\t\"buyerApp\"\t\"std:040\"\t\"acc.in\"
\"IND\"\t\"ONDC:RET10\"\t\"sellerApp\"\t\"std:040\"\t\"jio.com\"
\"IND\"\t\"ONDC:TRV01\"\t\"gateway\"\t\"std:044\"\t\"zomato.com\"
\"IND\"\t\"ONDC:RET10\"\t\"sellerApp\"\t\"std:044\"\t\"acer.in\"
\"IND\"\t\"ONDC:RET10\"\t\"sellerApp\"\t\"std:040\"\t\"indusind.com\"
\"IND\"\t\"ONDC:RET10\"\t\"gateway\"\t\"std:040\"\t\"indusind.com\"
\"IND\"\t\"ONDC:LOG12\"\t\"gateway\"\t\"std:080\"\t\"kotak.com\"
\"IND\"\t\"ONDC:LOG12\"\t\"sellerApp\"\t\"std:080\"\t\"sbi.nic.in\"
\"IND\"\t\"ONDC:FOD01\"\t\"sellerApp\"\t\"std:044\"\t\"flipkart.com\"
\"IND\"\t\"ONDC:TRV01\"\t\"gateway\"\t\"std:0755\"\t\"reliancefresh.in\"
\"IND\"\t\"ONDC:LOG12\"\t\"buyerApp\"\t\"std:040\"\t\"amazon.in\"
\"IND\"\t\"ONDC:RET10\"\t\"sellerApp\"\t\"std:033\"\t\"bigbasket.com\"
\"IND\"\t\"ONDC:LOG12\"\t\"sellerApp\"\t\"std:022\"\t\"flipkart.com\"
\"IND\"\t\"ONDC:LOG12\"\t\"gateway\"\t\"std:044\"\t\"reliancefresh.in\"
\"IND\"\t\"ONDC:MED01\"\t\"sellerApp\"\t\"std:080\"\t\"icicibank.com\"
\"IND\"\t\"ONDC:FOD01\"\t\"gateway\"\t\"std:044\"\t\"kotak.com\"
\"IND\"\t\"ONDC:MED01\"\t\"sellerApp\"\t\"std:011\"\t\"bharatpe.in\"
\"IND\"\t\"ONDC:MED01\"\t\"sellerApp\"\t\"std:080\"\t\"snapdeal.com\"
\"IND\"\t\"ONDC:LOG12\"\t\"sellerApp\"\t\"std:040\"\t\"bharatpe.in\"
\"IND\"\t\"ONDC:RET10\"\t\"buyerApp\"\t\"std:011\"\t\"indusind.com\"
\"IND\"\t\"ONDC:FOD01\"\t\"buyerApp\"\t\"std:022\"\t\"ajio.com\"
\"IND\"\t\"ONDC:FOD01\"\t\"gateway\"\t\"std:011\"\t\"reliancefresh.in\"
\"IND\"\t\"ONDC:MED01\"\t\"sellerApp\"\t\"std:080\"\t\"acer.in\"
\"IND\"\t\"ONDC:RET10\"\t\"sellerApp\"\t\"std:040\"\t\"idfcfirstbank.com\"
\"IND\"\t\"ONDC:MED01\"\t\"sellerApp\"\t\"std:011\"\t\"myntra.com\"
\"IND\"\t\"ONDC:RET10\"\t\"sellerApp\"\t\"std:044\"\t\"mobikwik.com\"
\"IND\"\t\"ONDC:TRV01\"\t\"gateway\"\t\"std:022\"\t\"vodafone.in\"
\"IND\"\t\"ONDC:FOD01\"\t\"sellerApp\"\t\"std:0755\"\t\"sbi.nic.in\"
\"IND\"\t\"ONDC:TRV01\"\t\"buyerApp\"\t\"std:044\"\t\"acer.in\"
\"IND\"\t\"ONDC:LOG12\"\t\"gateway\"\t\"std:011\"\t\"ajio.com\"
\"IND\"\t\"ONDC:FOD01\"\t\"gateway\"\t\"std:0755\"\t\"blinkit.com\"
\"IND\"\t\"ONDC:RET10\"\t\"buyerApp\"\t\"std:022\"\t\"reliancefresh.in\"
\"IND\"\t\"ONDC:LOG12\"\t\"buyerApp\"\t\"std:080\"\t\"blinkit.com\"
