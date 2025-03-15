import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
  stages: [
    { duration: '15s', target: 5 },
    { duration: '15s', target: 30 },
    { duration: '15s', target: 0 },
  ],
  thresholds: {
    http_req_failed: ['rate<0.01'], // Less than 5% errors
    http_req_duration: ['p(95)<500'], // 95% of requests should be below 500ms
  },
};

const BASE_URL = 'http://localhost:8080';
const SOURCE_TYPES = ['game', 'server', 'payment'];
const EXPECTED_ERRORS = [422]; // Add expected error codes

function getRandomSourceType() {
  return SOURCE_TYPES[Math.floor(Math.random() * SOURCE_TYPES.length)];
}

export default function () {
  const userId = Math.floor(Math.random() * 3) + 1;
  const transactionId = `tx_${new Date().getTime()}_${Math.random()}`;
  const sourceType = getRandomSourceType();

  // Test balance endpoint
  const balanceCheck = http.get(`${BASE_URL}/user/${userId}/balance`);
  check(balanceCheck, {
    'balance status is 200': (r) => r.status === 200,
  });

  // Test transaction endpoint
  const payload = JSON.stringify({
    state: Math.random() > 0.5 ? 'win' : 'lose',
    amount: (Math.random() * 10).toFixed(2),
    transactionId: transactionId,
  });

  const params = {
    headers: {
      'Content-Type': 'application/json',
      'Source-Type': sourceType,
    },
  };

  const transaction = http.post(
    `${BASE_URL}/user/${userId}/transaction`,
    payload,
    params
  );

  check(transaction, {
    'transaction status is valid': (r) => {
      if (r.status !== 200 && !EXPECTED_ERRORS.includes(r.status)) {
        console.log(`Transaction failed with status ${r.status}: ${r.body}`);
        return false;
      }
      return true;
    },
    'source type is valid': (r) => {
      // Get the first value from the header array
      const sourceTypeHeader = r.request.headers['Source-Type'] || r.request.headers['source-type'];
      const sentSourceType = Array.isArray(sourceTypeHeader) ? sourceTypeHeader[0] : sourceTypeHeader;
      
      console.log('Raw header value:', sourceTypeHeader);
      console.log('Processed source type:', sentSourceType);
      console.log('Valid types:', SOURCE_TYPES);
      
      const isValid = SOURCE_TYPES.includes(sentSourceType);
      console.log('Is valid source type:', isValid);
      
      return isValid;
    },
  });

  sleep(1);
}