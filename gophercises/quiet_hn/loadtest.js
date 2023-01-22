import http from 'k6/http';
import { sleep } from 'k6';
// k6 run loadtest.js
export const options = {
  vus: 500,
  duration: '30s',
};
export default function () {
  http.get('http://localhost:3000');
  sleep(1);
}
