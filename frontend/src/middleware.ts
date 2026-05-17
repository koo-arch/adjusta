import { NextRequest, NextResponse } from 'next/server';
import { cookies } from 'next/headers';

const publicRoutes: string[] = [
    '/',
    '/login'
];

export function middleware(request: NextRequest) {
    const token = cookies().get('access_token');
    const { pathname } = new URL(request.url);


    // 認証トークンがない場合、ログインページにリダイレクト
    if (!token) {
        if (!publicRoutes.includes(pathname)) {
            return NextResponse.redirect(new URL('/login', request.url));
        }
    } else {
        if (publicRoutes.includes(pathname)) {
            return NextResponse.redirect(new URL('/dashboard', request.url));
        }
    }

    return NextResponse.next();
}

export const config = {
    matcher: '/((?!api|_next/static|_next/image|favicon.ico|.*\\.svg$).*)',
}