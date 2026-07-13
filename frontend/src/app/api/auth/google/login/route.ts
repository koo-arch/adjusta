import { proxyOAuthRequest } from '@/lib/server/oauthProxy';

export function GET(request: Request) {
    return proxyOAuthRequest(request, '/auth/google/login');
}
