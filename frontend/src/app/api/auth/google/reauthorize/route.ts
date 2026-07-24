import { NextResponse } from 'next/server';
import { proxyOAuthRequest } from '@/lib/server/oauthProxy';

export async function GET(request: Request) {
    const { search } = new URL(request.url);
    const response = await proxyOAuthRequest(request, `/api/auth/google/reauthorize${search}`);
    if (response.status === 401) {
        return NextResponse.redirect(new URL('/api/auth/session-expired', request.url), 303);
    }
    return response;
}
