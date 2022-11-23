import http from 'k6/http';
import { check, group, sleep } from 'k6';

export let options = {
    stages: [
        { duration: '1m', target: 200 }, 
        { duration: '1m', target: 400 },
        { duration: '1m', target: 800 },
        { duration: '1m', target: 1600 },
        { duration: '1m', target: 2400 },
    ],
    thresholds: {
        http_req_duration: ['p(99)<2000'], // 99% of requests must complete below 2s
    },
};

export default () => {
    let res = http.get('http://10.10.10.10:8080/api/v1/summary');
    check(res, {
        'is status 200': (r) => r.status === 200,
    });
    sleep(1);
}
