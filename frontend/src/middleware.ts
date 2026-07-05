import { NextRequest, NextResponse } from 'next/server';

const publicRoutes: string[] = [
    '/',
    '/login'
];

const hasValidSession = async (request: NextRequest) => {
    const backendURL = process.env.INTERNAL_BACKEND_URL || process.env.NEXT_PUBLIC_API_BASE_URL;
    if (!backendURL) {
        return request.cookies.has('session');
    }

    try {
        const response = await fetch(`${backendURL}/api/users/me`, {
            method: 'GET',
            headers: {
                cookie: request.headers.get('cookie') || '',
            },
            cache: 'no-store',
            redirect: 'manual',
        });

        return response.status === 200;
    } catch {
        return request.cookies.has('session');
    }
};

export async function middleware(request: NextRequest) {
    const hasSessionCookie = request.cookies.has('session');
    const { pathname } = new URL(request.url);

    if (!hasSessionCookie) {
        if (!publicRoutes.includes(pathname)) {
            return NextResponse.redirect(new URL('/login', request.url));
        }

        return NextResponse.next();
    }

    const isAuthenticated = await hasValidSession(request);

    if (!isAuthenticated) {
        if (!publicRoutes.includes(pathname)) {
            return NextResponse.redirect(new URL('/login', request.url));
        }

        return NextResponse.next();
    }

    if (publicRoutes.includes(pathname)) {
        return NextResponse.redirect(new URL('/dashboard', request.url));
    }

    return NextResponse.next();
}

export const config = {
    matcher: '/((?!api|_next/static|_next/image|favicon.ico|.*\\.svg$).*)',
}
