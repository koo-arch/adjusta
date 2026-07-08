import { NextRequest, NextResponse } from 'next/server';

const publicRoutes: string[] = [
    '/',
    '/login'
];

// cookie の「存在」だけを見る楽観的な UX 振り分け。セッションの検証はしない。
// 権威的な認証チェックは Go API(全 /api/* のセッションミドルウェア)であり、
// 期限切れ cookie がここをすり抜けてもデータ取得は 401 になり
// serverApi / QueryCache 側で cookie 失効後に /login へ回収される。
export function proxy(request: NextRequest) {
    const hasSessionCookie = request.cookies.has('session');
    const { pathname } = request.nextUrl;

    if (!hasSessionCookie && !publicRoutes.includes(pathname)) {
        return NextResponse.redirect(new URL('/login', request.url));
    }

    // RSC 経路の 401 でも /api/auth/session-expired が cookie を失効させるため、
    // 期限切れ cookie が /login と /dashboard を往復し続ける状態にはならない。
    if (hasSessionCookie && (pathname === '/' || pathname === '/login')) {
        return NextResponse.redirect(new URL('/dashboard', request.url));
    }

    return NextResponse.next();
}

export const config = {
    matcher: '/((?!api|_next/static|_next/image|favicon.ico|.*\\.(?:svg|png|jpg|jpeg|webp|gif)$).*)',
}
