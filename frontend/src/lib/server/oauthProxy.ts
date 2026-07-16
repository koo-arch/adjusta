import { NextResponse } from 'next/server';

const redirectStatuses = new Set([301, 302, 303, 307, 308]);

export const proxyOAuthRequest = async (request: Request, path: string) => {
    const backend = process.env.INTERNAL_BACKEND_URL?.replace(/\/$/, '');
    if (!backend) {
        return NextResponse.json(
            { error: '認証サービスの接続先が設定されていません' },
            { status: 500 },
        );
    }

    // Request 依存の読み取りは try の外で行い、Next.js のプリレンダリング中断を
    // backend 接続エラーとして握りつぶさない。
    const cookie = request.headers.get('cookie') ?? '';

    try {
        const backendResponse = await fetch(`${backend}${path}`, {
            method: 'GET',
            headers: {
                cookie,
            },
            redirect: 'manual',
            cache: 'no-store',
        });

        const location = backendResponse.headers.get('location');
        const response =
            location && redirectStatuses.has(backendResponse.status)
                ? NextResponse.redirect(location, backendResponse.status)
                : new NextResponse(backendResponse.body, {
                      status: backendResponse.status,
                      statusText: backendResponse.statusText,
                      headers: {
                          'content-type':
                              backendResponse.headers.get('content-type') ??
                              'application/json',
                      },
                  });

        // backendはOAuth stateまたはlogin sessionを1つのcookieとして返す。
        const setCookie = backendResponse.headers.get('set-cookie');
        if (setCookie) {
            response.headers.append('set-cookie', setCookie);
        }

        return response;
    } catch (error) {
        console.error('failed to proxy OAuth request', error);
        return NextResponse.json(
            { error: '認証サービスに接続できませんでした' },
            { status: 502 },
        );
    }
};
