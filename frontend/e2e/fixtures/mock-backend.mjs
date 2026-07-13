import { createServer } from 'node:http';

const port = 3101;
const expiredSessions = new Set();

const json = (response, status, body) => {
    response.writeHead(status, { 'Content-Type': 'application/json' });
    response.end(JSON.stringify(body));
};

const sessionToken = (request) => {
    const cookie = request.headers.cookie ?? '';
    return cookie.match(/(?:^|;\s*)session=([^;]+)/)?.[1];
};

const user = {
    sub: 'e2e-user',
    name: 'E2E User',
    email: 'e2e@example.com',
    picture: '',
};

const server = createServer((request, response) => {
    if (request.url === '/health') {
        response.writeHead(200).end('ok');
        return;
    }

    const expireMatch = request.url?.match(/^\/__e2e\/sessions\/([^/]+)\/expire$/);
    if (request.method === 'POST' && expireMatch) {
        expiredSessions.add(decodeURIComponent(expireMatch[1]));
        response.writeHead(204).end();
        return;
    }

    const token = sessionToken(request);
    if (!token || token === 'expired-session' || expiredSessions.has(token)) {
        json(response, 401, { code: 'unauthorized', error: '認証情報がありません' });
        return;
    }

    if (request.url === '/api/users/me' || request.url === '/api/account/list') {
        json(response, 200, user);
        return;
    }

    if (request.url === '/api/user-calendars') {
        json(response, 200, []);
        return;
    }

    if (request.url === '/api/calendar-settings/candidate-sync') {
        json(response, 200, { enabled: false, calendar: null });
        return;
    }

    json(response, 404, { code: 'not_found', error: 'Not Found' });
});

server.listen(port, 'localhost');

const closeServer = () => server.close();

process.on('SIGINT', closeServer);
process.on('SIGTERM', closeServer);
