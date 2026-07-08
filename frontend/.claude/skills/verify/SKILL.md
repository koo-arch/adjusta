---
name: verify
description: Adjusta frontend の変更をエンドツーエンドで動作確認する手順(build → start → curl でルーティング/認証振り分けを観察)
---

# Adjusta frontend の動作確認

## ビルドと起動

```bash
cd /workspace/frontend
npm run build          # ルートテーブルで / と /login が ○(Static)、(app) 配下が ƒ(Dynamic)、Proxy (Middleware) の行が出ることを確認
nohup npm run start -- -p 3000 -H 0.0.0.0 > /tmp/next-start.log 2>&1 &
```

- Go バックエンドは devcontainer から `http://backend:8080` で到達可能(`curl http://backend:8080/api/users/me` → 401 `{"error":"認証情報がありません"}` なら生きている)
- サーバー側 fetch は `INTERNAL_BACKEND_URL`、ブラウザ側は `NEXT_PUBLIC_API_BASE_URL`(`.env.local`)

## 認証ルーティングの確認(curl)

セッション cookie 名は `session`。proxy は存在チェックのみなので偽値で「期限切れ cookie」を再現できる。

```bash
B=http://localhost:3000
curl -s -o /dev/null -w "%{http_code} -> %{redirect_url}\n" $B/dashboard                        # 307 -> /login
curl -s -o /dev/null -w "%{http_code} -> %{redirect_url}\n" -H "Cookie: session=fake" $B/       # 307 -> /dashboard
curl -s -o /dev/null -w "%{http_code} -> %{redirect_url}\n" -H "Cookie: session=fake" $B/login  # 307 -> /dashboard(ログイン済みは /login を見ない)
curl -s -H "Cookie: session=fake" $B/dashboard | grep -o "NEXT_REDIRECT[^\"]*"                  # NEXT_REDIRECT;replace;/api/auth/session-expired;303
curl -si -H "Cookie: session=fake" $B/api/auth/session-expired | grep -i "location\|set-cookie" # 303 /login + Set-Cookie: session=; Max-Age=0(cookie 失効がループブレーカー)
curl -s -o /dev/null -w "%{http_code}\n" $B/api/users/me                                        # 401(proxy は api を素通し、backend 401 がそのまま返る)
curl -s -o /dev/null -w "%{http_code}\n" $B/images/schedule_manage.jpg                          # 200(matcher が画像を除外)
```

期限切れ cookie のループ収束は cookie jar で全周回を実測できる:

```bash
J=$(mktemp); curl -s -c $J -b "session=fake" -o /dev/null $B/api/auth/session-expired
grep -c session $J   # 0(jar から消えている)→ 以後 /login は 200、/dashboard は 307 -> /login
```

シェル構成の確認: LP は `ログイン`(MarketingHeader)を含み `イベント一覧`(app Header)を含まない。`/login` は両方含まない。stale cookie の `/dashboard` は `ホーム` / `イベント一覧` / `animate-pulse`(UserMenuSkeleton)を含む。

## serverApi の送信ヘッダを実測する(実ログインなしで DAL を検証)

`INTERNAL_BACKEND_URL` を Node 製フェイク backend に向けると、`requireUser` が送る生の Cookie ヘッダを観察でき、有効セッション時の描画(UserButton など)も再現できる。

```bash
node scratchpad/echo_backend.mjs &   # :9999 で Cookie をログし、既知の値なら AuthUser JSON / それ以外 401 を返す
INTERNAL_BACKEND_URL=http://localhost:9999 npm run start -- -p 3000 &
curl -s -H "Cookie: session=MTc0OTk+dGVzdC9zaWc=" http://localhost:3000/dashboard | grep "Test User"
```

**注意(実バグの教訓)**: gin のセッション cookie 値は base64 パディング `=` を含む。Next の `cookies().toString()` は値を `encodeURIComponent` して再構築するため `%3D` に化けて backend 検証が壊れる。サーバー側で cookie を転送するときは必ず `(await headers()).get('cookie')` の生ヘッダを使うこと(`lib/server/api.ts` 参照)。cookie 転送を触ったら、`=` `+` `/` を含む値がフェイク backend に無傷で届くことを必ず確認する。

## ここでは検証できないもの(手動確認が必要)

- QueryCache / MutationCache の 401 → `window.location.assign('/api/auth/session-expired')` と 409 → AuthErrorModal(ブラウザJS + 実セッションが必要。headless ブラウザ未導入)
- 実ログイン後の UserMenu アバター表示(Google OAuth が必要)
- Header 回帰(スクロールシャドウ、モバイルメニュー、テーマ切替)
