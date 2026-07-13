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

const event = (id, title) => ({
    id,
    title,
    description: '',
    location: '',
    status: 'active',
    sync_status: 'not_synced',
    confirmed_date_id: null,
    proposed_dates: [],
});

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

    if (
        request.url === '/api/event/draft/needs-action' ||
        request.url === '/api/event/confirmed/upcoming'
    ) {
        json(response, 200, []);
        return;
    }

    if (request.url === '/api/calendar/list') {
        json(response, 200, { events: [], warning: { failed_calendars: [] } });
        return;
    }

    if (request.url?.startsWith('/api/event/draft/search')) {
        const url = new URL(request.url, `http://localhost:${port}`);
        const page = Number(url.searchParams.get('page') ?? '1');
        const title = url.searchParams.get('title');

        if (title === '表示確認') {
            json(response, 200, {
                items: [event('visible-event', '表示確認イベント')],
                pagination: {
                    page,
                    per_page: 20,
                    total_items: 1,
                    total_pages: 1,
                },
            });
            return;
        }

        if (title === 'ページング') {
            json(response, 200, {
                items: [event(`page-${page}-event`, `${page}ページ目のイベント`)],
                pagination: {
                    page,
                    per_page: 20,
                    total_items: 21,
                    total_pages: 2,
                },
            });
            return;
        }

        json(response, 200, {
            items: [],
            pagination: {
                page,
                per_page: 20,
                total_items: 0,
                total_pages: 0,
            },
        });
        return;
    }

    if (request.method === 'POST' && request.url === '/api/calendar/event/draft') {
        json(response, 201, { id: 'created-event' });
        return;
    }

    if (request.method === 'GET' && request.url === '/api/calendar/event/draft/created-event') {
        json(response, 200, event('created-event', 'E2E作成イベント'));
        return;
    }

    if (request.method === 'GET' && request.url === '/api/calendar/event/draft/detail-event') {
        json(response, 200, {
            ...event('detail-event', '詳細確認イベント'),
            description: 'イベント詳細の説明',
            location: '会議室A',
            proposed_dates: [
                {
                    id: 'candidate-1',
                    start: '2026-07-20T01:00:00.000Z',
                    end: '2026-07-20T02:00:00.000Z',
                    priority: 1,
                    status: 'active',
                    sync_status: 'not_synced',
                },
            ],
        });
        return;
    }

    if (request.method === 'GET' && request.url === '/api/calendar/event/draft/missing-event') {
        json(response, 404, { code: 'not_found', error: 'イベントが見つかりません' });
        return;
    }

    json(response, 404, { code: 'not_found', error: 'Not Found' });
});

server.listen(port, 'localhost');

const closeServer = () => server.close();

process.on('SIGINT', closeServer);
process.on('SIGTERM', closeServer);
