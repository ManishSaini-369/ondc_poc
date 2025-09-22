import http from "k6/http";
import { check, sleep } from "k6";

// Load test configuration
export const options = {
  vus: 100,        // 100 concurrent users
  duration: "15s", // run for 15 seconds
};

// Dataset (you can expand this as needed)
const ukids = [
  "UKID001", "UKID002", "UKID003", "UKID004", "UKID005",
  "UKID006", "UKID007", "UKID008", "UKID009", "UKID010",
  "UKID011", "UKID012", "UKID013", "UKID014", "UKID015",
  "UKID016", "UKID017", "UKID018", "UKID019", "UKID020",
  "UKID021", "UKID022", "UKID023", "UKID024", "UKID025",
  "UKID026", "UKID027", "UKID028", "UKID029", "UKID030",
  "UKID031", "UKID032", "UKID033", "UKID034", "UKID035",
  "UKID036", "UKID037", "UKID038", "UKID039", "UKID040",
  "UKID041", "UKID042", "UKID043", "UKID044", "UKID045",
  "UKID046", "UKID047", "UKID048", "UKID049", "UKID050",
  // ... continue up to UKID100 or more
];

export default function () {
  const url = "http://localhost:8081/lookup";

  // Pick a random UKID for each request
  const randomUkid = ukids[Math.floor(Math.random() * ukids.length)];

  const payload = JSON.stringify({ ukid: randomUkid });

  const params = {
    headers: {
      "Content-Type": "application/json",
    },
  };

  let res = http.post(url, payload, params);

  check(res, {
    "status is 200": (r) => r.status === 200,
    "response contains subscriber_id": (r) =>
      Array.isArray(r.json()) && r.json()[0].hasOwnProperty("subscriber_id"),
  });

  sleep(1); // small wait between requests
}
