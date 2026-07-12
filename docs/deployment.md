# デプロイ方針メモ

2026-07-12 作成。初回デプロイの構成判断と、実施時に必要な設定の一覧。未決事項は末尾に明示する。

## 結論(推奨構成)

**Vercel(frontend)+ Cloud Run(backend)+ マネージド PostgreSQL** — 「Cloud Run + Vercel で良いか」への回答は **Yes**。

| レイヤー | サービス | 根拠 |
| --- | --- | --- |
| frontend(Next.js) | Vercel | App Router・middleware(認証 proxy)・next/font をそのまま運用でき、Git 連携で CI/CD が不要。セルフホスト(Cloud Run に載せる)も可能だが、得られるものが少ない |
| backend(Go + ent) | Cloud Run | Go は Vercel に載らないためコンテナ運用が前提。`backend/Dockerfile` が既にあり、scale-to-zero で個人運用のコストが小さい |
| DB(PostgreSQL) | Neon(個人運用)/ Cloud SQL(本格運用) | 下記「DB の選択」 |

## この構成が本アプリと相性が良い理由(重要)

ブラウザからの API 呼び出しは **Next.js の `/api/[...path]` route handler がサーバー側で backend へ転送する proxy 構成**になっている(`frontend/src/lib/api/client.ts` の baseURL は `NEXT_PUBLIC_API_BASE_URL` 未設定時に空 = same-origin)。つまり:

- **session cookie は Vercel ドメインの first-party のまま** — クロスドメイン cookie / SameSite=None / CORS の問題が発生しない
- backend(Cloud Run)の URL はブラウザに露出せず、`INTERNAL_BACKEND_URL` としてサーバー側にだけ設定する
- 本番では **`NEXT_PUBLIC_API_BASE_URL` は未設定(空)のままにする**こと。値を入れるとブラウザ直アクセスになり、上記の利点が崩れる
- OAuth 開始・コールバックも `/api/auth/*` の route handler 経由で成立している(`frontend/src/app/api/`)

## DB の選択

| 選択肢 | 向き | 備考 |
| --- | --- | --- |
| **Neon**(serverless Postgres) | 個人運用・初回リリース | 無料枠で開始でき、scale-to-zero。Cloud Run からは接続文字列だけで繋がる。コールドスタート時の接続レイテンシは許容範囲 |
| Cloud SQL for PostgreSQL | 本格運用 | 同一リージョンでレイテンシ最小・安定。最小構成でも固定費(月 $10〜)が掛かる。Cloud Run からは Cloud SQL コネクタ or プライベート IP |
| Supabase | Neon 同等 | Auth 等の付随機能は使わない(自前 OAuth があるため)なら Neon とほぼ等価 |

推奨: **まず Neon で開始し、負荷・運用要件が固まったら Cloud SQL を再検討**。

## 実施時に必要な設定

### 環境変数

| 変数 | 置き場所 | 値 |
| --- | --- | --- |
| `INTERNAL_BACKEND_URL` | Vercel(server) | Cloud Run の URL(例: `https://adjusta-api-xxxx.a.run.app`) |
| `NEXT_PUBLIC_API_BASE_URL` | Vercel | **設定しない**(same-origin proxy を維持) |
| `SESSION_SECRET` | Cloud Run | 十分に長いランダム値 |
| DB DSN(`DATABASE_URL` 等、`backend/internal/config` 参照) | Cloud Run | Neon / Cloud SQL の接続文字列 |
| `CORS_ALLOW_ORIGINS` | Cloud Run | Vercel のドメイン(proxy 経由ならサーバー間通信のため実質不要だが、設定しておく) |
| Google OAuth client ID / secret | Cloud Run | GCP コンソールで発行 |

### Google OAuth

- 承認済みリダイレクト URI に **backend の公開 URL**(`https://<cloud-run>/auth/google/callback`)を登録
- OAuth 同意画面の公開設定(テストユーザー→本番公開)を確認

### Cloud Run

- `backend/Dockerfile` をそのままビルド(GitHub Actions から `gcloud run deploy` が簡単)
- **min-instances**: 0 だと初回アクセス・OAuth コールバックにコールドスタート遅延が乗る。体感を優先するなら 1(常時課金)。まず 0 で始めて気になったら上げる
- cookie の `Secure` / ドメイン設定が本番 URL 前提になっているか `backend/api/cookie` を確認

### CI/CD

- frontend: Vercel の Git 連携(main への push で本番、PR ごとに Preview)
- backend: GitHub Actions → Artifact Registry へ push → Cloud Run デプロイ(`.github/workflows` は未整備、要作成)

## 未決事項(実施時に決める)

1. **独自ドメイン**を取るか(取るなら Vercel に割当。backend はブラウザ非公開のため Cloud Run ドメインのままでよい)
2. **DB の最終選択**(推奨は Neon スタート)
3. **staging 環境**の要否(Vercel Preview + Cloud Run のリビジョンタグで軽く済ませる案が有力)
4. Cloud Run の **min-instances**(0 か 1 か)
5. DB マイグレーションの実行方法(ent のマイグレーションをデプロイパイプラインに組み込むか、手動か)
