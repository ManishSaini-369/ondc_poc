import http from 'k6/http';
import { check, sleep } from 'k6';
import { uuidv4 } from 'https://jslib.k6.io/k6-utils/1.4.0/index.js';

export const options = {
  vus: 100,
  duration: '10s',
};

const url = 'http://localhost:8081/search';

export default function () {
  const transactionId = uuidv4();
  const messageId = uuidv4();

  const payload = JSON.stringify({
  "context": {
    "domain": "retail",
    "city": "DEL",
    "action": "search",
    "core_version": "1.0.0",
    "bap_id": "buyer-app.ondc.org",
    "bap_uri": "https://buyer-app.ondc.org",
    "transaction_id": "12345",
    "message_id": "abc-123",
    "timestamp": "2025-09-29T18:30:00Z"
  },
  "message": {
    "intent": {
      "item": {
        "descriptor": {
          "name": "Laptop"
        }
      }
    }
  }
});

  // WARNING: The Authorization and Digest headers are static values from the example.
  // For a valid test, these would need to be dynamically generated based on the payload and a private key.
  const params = {
    headers: {
      'Content-Type': 'application/json',
      'Authorization': 'Signature keyId="UK1001",algorithm="ed25519",signature="QR3mInFDxRVaCPiWrKgkNBWuhATbOPn5blUsxmzQfdXRgy4N0fjWcp8PDGbg7QeDFMe+3ek0035T+R/CxnhBBQ=="',
      'Digest': 'Rfv1pdi+f9ZFwMJW4egU26kNX3WrdbsKZt8lAZCprKI=',
    },
  };

  const res = http.post(url, payload, params);

  check(res, {
    'search status is 200': (r) => r.status === 200,
  });

  sleep(1);
}
