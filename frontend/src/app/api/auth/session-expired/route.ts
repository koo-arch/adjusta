import { NextResponse } from "next/server";

const sessionCookieName = "session";

const sessionCookieDomain = process.env.DOMAIN || undefined;
const appEnv = process.env.GO_ENV ?? process.env.NODE_ENV;
const isSecureCookie = appEnv !== "development";

export async function GET(request: Request) {
    const response = NextResponse.redirect(new URL("/login", request.url), 303);

    response.cookies.set(sessionCookieName, "", {
        domain: sessionCookieDomain,
        path: "/",
        maxAge: 0,
        httpOnly: true,
        secure: isSecureCookie,
        sameSite: isSecureCookie ? "none" : "lax",
    });

    return response;
}
