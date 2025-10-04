import http from 'k6/http';
import { check, sleep } from 'k6';
import { uuidv4 } from 'https://jslib.k6.io/k6-utils/1.4.0/index.js';

export const options = {
  vus: 100,
  duration: '10s',
};

const url = 'http://localhost:8081/vlookup';
const params = {
  headers: {
    'Content-Type': 'application/json',
  },
};

export default function () {
  const payload = JSON.stringify({
    "sender_subscriber_id": "Amazon.com",
    "request_id": uuidv4(), // Generate a unique request_id for each iteration
    "timestamp": new Date().toISOString(), // Use the current timestamp
    "signature": "BSCm2wTw3lpCOK7aHf27CNSOfDhrkInwheIf5qja9y7CdsD23YKQb5pUJvKJW+XJIr/vjSW8AtUHI7zyn5RwAQ==", // Note: Using a static signature
    "search_parameters": {
        "country": "IN",
        "domain": "ONDC:ENT01",
        "type": "Seller",
        "city": "Hyderabad",
        "subscriber_id": "Paytm.com"
    }
  });

  const res = http.post(url, payload, params);

  check(res, {
    'vlookup status is 200': (r) => r.status === 200,
  });

  sleep(1); // Wait for 1 second between requests
}
