import { NextResponse } from "next/server";

export async function GET(req: Request) {
    const { search } = new URL(req.url);
    const backend = process.env.BACKEND_URL!;
    const url = `${backend}/auth/google/callback${search}`;

    // リダイレクトせずにバックエンドを呼び出し
    const apiRes = await fetch(url, {
        method: "GET",
        headers: {
            // ブラウザから受け取った Cookie をそのまま渡す
            cookie: req.headers.get("cookie") || ""
        },
        redirect: "manual",
    });

    // NextResponse でリダイレクト先とステータスを用意
    const res = NextResponse.redirect(new URL("/", req.url), 307);

    // バックエンドが返した Set-Cookie ヘッダーをすべて転送
    apiRes.headers.forEach((value, key) => {
        if (key.toLowerCase() === "set-cookie") {
            res.headers.append("Set-Cookie", value);
        }
    });

    return res;
}