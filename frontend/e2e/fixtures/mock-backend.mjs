import { createServer } from 'node:http';

const port = 3101;
const expiredSessions = new Set();
const firstCandidateID = '22222222-2222-4222-8222-222222222222';
const secondCandidateID = '33333333-3333-4333-8333-333333333333';
const confirmedEventDates = new Map([['confirmed-event', firstCandidateID]]);
const candidateSyncState = {
    enabled: false,
    calendarExists: false,
    failUpdate: false,
};
const calendarSettingsState = {
    primaryID: 'calendar-primary',
    primaryVisible: true,
    referenceVisible: true,
    failUpdate: false,
};

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

const candidateCalendar = {
    id: 'calendar-adjusta',
    calendar_id: 'calendar-entity-adjusta',
    google_calendar_id: 'adjusta-candidate@group.calendar.google.com',
    summary: 'Adjusta 調整用',
    timezone: 'Asia/Tokyo',
    role: 'adjusta_candidate',
    is_visible: true,
    sync_proposed_dates: candidateSyncState.enabled,
};

const candidateSyncSetting = () => ({
    enabled: candidateSyncState.enabled,
    calendar: candidateSyncState.calendarExists
        ? { ...candidateCalendar, sync_proposed_dates: candidateSyncState.enabled }
        : null,
});

const calendarSettings = () => [
    {
        id: 'calendar-primary',
        calendar_id: 'calendar-entity-primary',
        google_calendar_id: 'e2e@example.com',
        summary: 'メインカレンダー',
        timezone: 'Asia/Tokyo',
        role: calendarSettingsState.primaryID === 'calendar-primary' ? 'primary' : 'reference',
        is_visible: calendarSettingsState.primaryVisible,
        sync_proposed_dates: false,
    },
    {
        id: 'calendar-reference',
        calendar_id: 'calendar-entity-reference',
        google_calendar_id: 'team@example.com',
        summary: 'チームカレンダー',
        timezone: 'Asia/Tokyo',
        role: calendarSettingsState.primaryID === 'calendar-reference' ? 'primary' : 'reference',
        is_visible: calendarSettingsState.referenceVisible,
        sync_proposed_dates: false,
    },
    {
        id: 'calendar-holiday',
        calendar_id: 'calendar-entity-holiday',
        google_calendar_id: 'ja.japanese#holiday@group.v.calendar.google.com',
        summary: '日本の祝日',
        timezone: 'Asia/Tokyo',
        role: 'reference',
        is_visible: true,
        sync_proposed_dates: false,
    },
    ...(candidateSyncState.calendarExists ? [candidateSyncSetting().calendar] : []),
];

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

const confirmationEvent = (id, confirmedDateID = null) => {
    const confirmed = confirmedDateID !== null;

    return {
        ...event(id, confirmed ? '確定済みイベント' : '日程確定イベント'),
        status: confirmed ? 'confirmed' : 'active',
        confirmed_date_id: confirmedDateID,
        confirmed_google_event_id: confirmed ? 'confirmed-google-event' : undefined,
        proposed_dates: [
            {
                id: firstCandidateID,
                google_event_id: 'candidate-google-event-1',
                start: '2026-07-22T01:00:00.000Z',
                end: '2026-07-22T02:00:00.000Z',
                priority: 1,
                status: confirmedDateID === firstCandidateID ? 'confirmed' : confirmed ? 'not_selected' : 'active',
                sync_status: 'synced',
            },
            {
                id: secondCandidateID,
                google_event_id: 'candidate-google-event-2',
                start: '2026-07-23T03:00:00.000Z',
                end: '2026-07-23T04:00:00.000Z',
                priority: 2,
                status: confirmedDateID === secondCandidateID ? 'confirmed' : confirmed ? 'not_selected' : 'active',
                sync_status: 'synced',
            },
        ],
    };
};

const server = createServer((request, response) => {
    if (request.url === '/health') {
        response.writeHead(200).end('ok');
        return;
    }

    if (request.method === 'GET' && request.url === '/auth/google/login') {
        response.writeHead(307, {
            Location: `http://localhost:${port}/__e2e/google/authorize`,
            'Set-Cookie': 'session=oauth-state-session; Path=/; HttpOnly; SameSite=Lax',
        });
        response.end();
        return;
    }

    if (request.method === 'GET' && request.url === '/__e2e/google/authorize') {
        response.writeHead(307, {
            Location: 'http://localhost:3100/api/auth/google/callback?code=e2e-code&state=e2e-state',
        });
        response.end();
        return;
    }

    if (request.method === 'GET' && request.url?.startsWith('/auth/google/callback?')) {
        if (sessionToken(request) !== 'oauth-state-session') {
            json(response, 400, { code: 'bad_request', error: 'stateが不正です' });
            return;
        }

        response.writeHead(307, {
            Location: 'http://localhost:3100/dashboard',
            'Set-Cookie': 'session=authenticated-session; Path=/; HttpOnly; SameSite=Lax',
        });
        response.end();
        return;
    }

    const expireMatch = request.url?.match(/^\/__e2e\/sessions\/([^/]+)\/expire$/);
    if (request.method === 'POST' && expireMatch) {
        expiredSessions.add(decodeURIComponent(expireMatch[1]));
        response.writeHead(204).end();
        return;
    }

    const candidateSyncFixtureMatch = request.url?.match(
        /^\/__e2e\/candidate-sync\/(off|on|fail-update)$/,
    );
    if (request.method === 'POST' && candidateSyncFixtureMatch) {
        const mode = candidateSyncFixtureMatch[1];
        candidateSyncState.enabled = mode === 'on';
        candidateSyncState.calendarExists = mode === 'on';
        candidateSyncState.failUpdate = mode === 'fail-update';
        response.writeHead(204).end();
        return;
    }

    const calendarSettingsFixtureMatch = request.url?.match(
        /^\/__e2e\/calendar-settings\/(reset|fail-update)$/,
    );
    if (request.method === 'POST' && calendarSettingsFixtureMatch) {
        calendarSettingsState.primaryID = 'calendar-primary';
        calendarSettingsState.primaryVisible = true;
        calendarSettingsState.referenceVisible = true;
        calendarSettingsState.failUpdate = calendarSettingsFixtureMatch[1] === 'fail-update';
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
        json(response, 200, calendarSettings());
        return;
    }

    const calendarSettingUpdateMatch = request.url?.match(
        /^\/api\/user-calendars\/(calendar-primary|calendar-reference)$/,
    );
    if (request.method === 'PATCH' && calendarSettingUpdateMatch) {
        if (calendarSettingsState.failUpdate) {
            json(response, 500, { code: 'internal_error', error: 'カレンダー設定を更新できませんでした' });
            return;
        }

        const id = calendarSettingUpdateMatch[1];
        if (id === 'calendar-reference') {
            calendarSettingsState.primaryID = id;
        } else {
            calendarSettingsState.primaryVisible = !calendarSettingsState.primaryVisible;
        }
        json(response, 200, calendarSettings().find((setting) => setting?.id === id));
        return;
    }

    if (request.method === 'GET' && request.url === '/api/calendar-settings/candidate-sync') {
        json(response, 200, candidateSyncSetting());
        return;
    }

    if (request.method === 'PUT' && request.url === '/api/calendar-settings/candidate-sync') {
        if (candidateSyncState.failUpdate) {
            json(response, 500, { code: 'internal_error', error: '同期設定を更新できませんでした' });
            return;
        }

        candidateSyncState.enabled = !candidateSyncState.enabled;
        candidateSyncState.calendarExists = true;
        json(response, 200, candidateSyncSetting());
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
                    id: 'candidate-3',
                    start: '2026-07-22T01:00:00.000Z',
                    end: '2026-07-22T02:00:00.000Z',
                    priority: 3072,
                    status: 'active',
                    sync_status: 'not_synced',
                },
                {
                    id: 'candidate-2',
                    start: '2026-07-21T01:00:00.000Z',
                    end: '2026-07-21T02:00:00.000Z',
                    priority: 2048,
                    status: 'active',
                    sync_status: 'not_synced',
                },
                {
                    id: 'candidate-1',
                    start: '2026-07-20T01:00:00.000Z',
                    end: '2026-07-20T02:00:00.000Z',
                    priority: 1024,
                    status: 'active',
                    sync_status: 'not_synced',
                },
            ],
        });
        return;
    }

    if (request.method === 'GET' && request.url === '/api/calendar/event/draft/edit-event') {
        json(response, 200, {
            ...event('edit-event', '編集前イベント'),
            description: '編集前の説明',
            location: '会議室B',
            proposed_dates: [
                {
                    id: '11111111-1111-4111-8111-111111111111',
                    start: '2026-07-21T01:00:00.000Z',
                    end: '2026-07-21T02:00:00.000Z',
                    priority: 1,
                    status: 'active',
                    sync_status: 'not_synced',
                },
            ],
        });
        return;
    }

    if (request.method === 'PUT' && request.url === '/api/calendar/event/draft/edit-event') {
        response.writeHead(204).end();
        return;
    }

    const deletionDetailMatch = request.url?.match(
        /^\/api\/calendar\/event\/draft\/(delete-event|delete-error-event)$/,
    );
    if (request.method === 'GET' && deletionDetailMatch) {
        const id = deletionDetailMatch[1];
        json(response, 200, event(id, id === 'delete-event' ? '削除対象イベント' : '削除失敗イベント'));
        return;
    }

    if (request.method === 'DELETE' && request.url === '/api/calendar/event/draft/delete-event') {
        json(response, 200, { message: 'success' });
        return;
    }

    if (request.method === 'DELETE' && request.url === '/api/calendar/event/draft/delete-error-event') {
        json(response, 500, { code: 'internal_error', error: 'イベントの削除処理に失敗しました' });
        return;
    }

    const confirmationDetailMatch = request.url?.match(
        /^\/api\/calendar\/event\/draft\/(confirm-event|confirm-error-event|confirmed-event)$/,
    );
    if (request.method === 'GET' && confirmationDetailMatch) {
        const id = confirmationDetailMatch[1];
        json(response, 200, confirmationEvent(id, confirmedEventDates.get(id) ?? null));
        return;
    }

    if (request.method === 'PATCH' && request.url === '/api/calendar/event/confirm/confirm-event') {
        confirmedEventDates.set('confirm-event', firstCandidateID);
        json(response, 200, { message: 'success' });
        return;
    }

    if (request.method === 'PATCH' && request.url === '/api/calendar/event/confirm/confirmed-event') {
        confirmedEventDates.set('confirmed-event', secondCandidateID);
        json(response, 200, { message: 'success' });
        return;
    }

    if (request.method === 'PATCH' && request.url === '/api/calendar/event/confirm/confirm-error-event') {
        json(response, 500, { code: 'internal_error', error: '日程の確定処理に失敗しました' });
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
