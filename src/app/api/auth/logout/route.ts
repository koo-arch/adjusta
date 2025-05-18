import { NextResponse } from "next/server";

export async function GET(req: Request) {
    const backend = process.env.INTERNAL_BACKEND_URL!;
    const url = `${backend}/auth/logout`;

    // クッキーを引き継いで呼び出し
    const apiRes = await fetch(url, {
        method: "GET",
        headers: {
            cookie: req.headers.get("cookie") || ""
        },
        redirect: "manual",
    });

    // ブラウザにはトップ（/）へリダイレクト
    const res = NextResponse.redirect(new URL("/", req.url), 307);

    // バックエンドの Set-Cookie (Max-Age=-1 など) を転送
    apiRes.headers.forEach((value, key) => {
        if (key.toLowerCase() === "set-cookie") {
            res.headers.append("Set-Cookie", value);
        }
    });

    return res;
}