import http from 'k6/http';
import {sleep, check} from 'k6';

export let options = {
    stages: [
        {duration: '2m', target: 100},
        {duration: '3m', target: 100},
        {duration: '1m', target: 0},
    ],
    thresholds: {
        http_req_duration: ['p(95)<50'],
        http_req_failed: ['rate<0.01'],
    }
};

export default function () {
    let response = http.get('http://localhost:8080/user_banner?tag_id=1&feature_id=4', {
        headers: {
            'Authorization': 'admin_token',
            'Content-Type': 'application/json'
        }
    });

    check(response, {
        'is status 200': r => r.status === 200,
    });

    sleep(1);
}
