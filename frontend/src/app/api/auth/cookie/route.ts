import { cookies } from 'next/headers';
import { NextResponse } from 'next/server';

export function GET() {
  const token = cookies().get('access_token');
  if (token) {
    return NextResponse.json({ exist : true }, { status: 200 });
  } else {
    return NextResponse.json({ exist : false }, { status: 401 });
  }
}