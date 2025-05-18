import { NextResponse } from "next/server";

export async function GET(request: Request) { return proxy(request) }
export async function POST(request: Request) { return proxy(request) }
export async function PUT(request: Request) { return proxy(request) }
export async function DELETE(request: Request) { return proxy(request) }
export async function PATCH(request: Request) { return proxy(request) }

async function proxy(request: Request) {
    const { pathname, search } = new URL(request.url);
    const backend = process.env.INTERNAL_BACKEND_URL;
    const url = `${backend}${pathname}${search}`;

    // リクエストボディとヘッダをそのまま転送
    const apiRes = await fetch(url, {
        method: request.method,
        headers: request.headers,
        body: ["GET", "HEAD"].includes(request.method) ? null : await request.text(),
        redirect: "manual",
    });

    const res = new NextResponse(apiRes.body, {
        status: apiRes.status,
        statusText: apiRes.statusText,
    });

    // ヘッダーをコピーしつつ、Set-Cookie は append で渡す
    apiRes.headers.forEach((value, key) => {
        if (key.toLowerCase() === "set-cookie") {
            // 複数ある場合も順に append すれば OK
            res.headers.append("Set-Cookie", value);
        } else {
            res.headers.set(key, value);
        }
    });

    return res;
}