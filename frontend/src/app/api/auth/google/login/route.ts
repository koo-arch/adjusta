import { NextResponse } from 'next/server';

export function GET() {
    return NextResponse.redirect(
        `${process.env.BACKEND_URL}/auth/google/login`
    );
}