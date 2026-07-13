import { proxyOAuthRequest } from '@/lib/server/oauthProxy';

export function GET(request: Request) {
    const { search } = new URL(request.url);
    return proxyOAuthRequest(request, `/auth/google/callback${search}`);
}
