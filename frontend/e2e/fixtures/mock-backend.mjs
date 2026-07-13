import { createServer } from 'node:http';

const port = 3101;

const server = createServer((request, response) => {
    if (request.url === '/health') {
        response.writeHead(200).end('ok');
        return;
    }

    if (request.url === '/api/users/me') {
        response.writeHead(401, { 'Content-Type': 'application/json' });
        response.end(JSON.stringify({ code: 'unauthorized', error: '認証情報がありません' }));
        return;
    }

    response.writeHead(404, { 'Content-Type': 'application/json' });
    response.end(JSON.stringify({ code: 'not_found', error: 'Not Found' }));
});

server.listen(port, 'localhost');

const closeServer = () => server.close();

process.on('SIGINT', closeServer);
process.on('SIGTERM', closeServer);
