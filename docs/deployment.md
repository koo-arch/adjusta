# デプロイ方針メモ

2026-07-12 作成、2026-07-26更新。初回デプロイの構成判断と、実施時に必要な設定の一覧。未決事項は末尾に明示する。

実際の初回構築、通常デプロイ、動作確認、ロールバックの操作手順は[`deployment-runbook.md`](./deployment-runbook.md)を参照する。

## 結論(推奨構成)

**Vercel(frontend)+ Cloud Run(backend)+ Cloud Tasks / Cloud Scheduler + マネージド PostgreSQL**を採用する。

| レイヤー | サービス | 根拠 |
| --- | --- | --- |
| frontend(Next.js) | Vercel | App Router・middleware(認証 proxy)・next/font をそのまま運用でき、Git 連携で CI/CD が不要。セルフホスト(Cloud Run に載せる)も可能だが、得られるものが少ない |
| backend(Go + ent) | Cloud Run | Go は Vercel に載らないためコンテナ運用が前提。`backend/Dockerfile` が既にあり、scale-to-zero で個人運用のコストが小さい |
| 非同期タスク | Cloud Tasks | Google Calendar API 呼び出しをユーザー向けリクエストから分離し、再試行と流量制御を行う |
| Outbox再送 | Cloud Scheduler | `dispatched_at`が未設定のOutbox Messageを定期的に検出し、Cloud Tasksへ再送する |
| DB(PostgreSQL) | Neon(個人運用)/ Cloud SQL(本格運用) | 下記「DB の選択」 |

## この構成が本アプリと相性が良い理由(重要)

通常のAPI呼び出しは **Next.js の `/api/[...path]` route handlerがサーバー側でbackendへ転送するproxy構成**になっている(`frontend/src/lib/api/client.ts`のbaseURLは`NEXT_PUBLIC_API_BASE_URL`未設定時に空 = same-origin)。OAuth開始・callbackも専用Route Handlerがbackendのredirectとcookieを中継する。

- **session cookie は Vercel ドメインの first-party のまま** — クロスドメイン cookie / SameSite=None / CORS の問題が発生しない
- session cookieは`HttpOnly`・`Secure`・host-only・`SameSite=Lax`とし、OAuth callbackを含む認証導線もsame-origin proxy経由で扱う
- backend(Cloud Run)の URL はブラウザに露出せず、`INTERNAL_BACKEND_URL` としてサーバー側にだけ設定する
- 本番では **`NEXT_PUBLIC_API_BASE_URL` は未設定(空)のままにする**こと。値を入れるとブラウザ直アクセスになり、上記の利点が崩れる
- OAuth開始・callbackも`/api/auth/*`のroute handler経由へ統一する(`frontend/src/app/api/`)

## DB の選択

| 選択肢 | 向き | 備考 |
| --- | --- | --- |
| **Neon**(serverless Postgres) | 個人運用・初回リリース | 無料枠で開始でき、scale-to-zero。Cloud Run からは接続文字列だけで繋がる。コールドスタート時の接続レイテンシは許容範囲 |
| Cloud SQL for PostgreSQL | 本格運用 | 同一リージョンでレイテンシ最小・安定。最小構成でも固定費(月 $10〜)が掛かる。Cloud Run からは Cloud SQL コネクタ or プライベート IP |
| Supabase | Neon 同等 | Auth 等の付随機能は使わない(自前 OAuth があるため)なら Neon とほぼ等価 |

初回productionは **Neonで開始**し、負荷・運用要件が固まったらCloud SQLを再検討する。

## Google Calendar 非同期同期

イベント作成・編集・確定・削除のユーザー向け API は、PostgreSQL の更新完了時点で応答する。Google Calendar API 呼び出しは Cloud Tasks から Cloud Run の task handler を呼び出して実行する。

```text
Cloud Run API
  ├─ 業務データ + adjusta.outbox_messages を同一 transaction で保存
  ├─ transaction commit 後に Cloud Tasks へ Outbox Message ID を enqueue
  └─ enqueue 失敗時は dispatched_at を設定せず配送エラーを保持

Cloud Tasks
  └─ OIDC token付きで Cloud Run task handler を呼び出す
       ├─ DBの最新状態から同期内容を構築
       ├─ Google Calendar APIを実行
       └─ sync_status と outbox_messages を更新

Cloud Scheduler
  └─ Outbox dispatcherを定期起動
       └─ dispatched_atが未設定かつavailable_atを過ぎたOutbox MessageをCloud Tasksへ再送
```

- queue名は`google-calendar-sync`とする
- task payloadにはOutbox Message IDだけを含める
- Cloud Tasks adapterはtask名をOutbox Message IDから安定して生成する。task名はoutbox_messagesへ保存しない
- Cloud Tasksとtask handlerは重複実行を前提とし、同期処理は冪等にする
- task handlerは通常のsession認証を使用せず、Cloud Tasks専用service accountのOIDC tokenを検証する
- `adjusta-api`はVercel proxyからの呼び出しに対応するため公開serviceのままとするが、task handlerのpathはアプリケーション側でOIDC tokenを必須とする
- Cloud SchedulerはGoogle Calendar同期そのものを実行せず、`dispatched_at`が未設定のOutbox MessageをCloud Tasksへ再送する回復経路として使用する

構成と認証方法は[Cloud TasksからCloud Runを呼び出す公式手順](https://docs.cloud.google.com/run/docs/triggering/using-tasks)に従う。

## 実施時に必要な設定

### 環境変数

| 変数 | 置き場所 | 値 |
| --- | --- | --- |
| `INTERNAL_BACKEND_URL` | Vercel(server) | `https://adjusta-api-1001878278191.asia-northeast1.run.app` |
| `NEXT_PUBLIC_API_BASE_URL` | Vercel | **設定しない**(same-origin proxy を維持) |
| `SESSION_SECRET` | Cloud Run | 十分に長いランダム値 |
| DB DSN(`DATABASE_URL` 等、`backend/internal/config` 参照) | Cloud Run | Neon / Cloud SQL の接続文字列 |
| `CORS_ALLOW_ORIGINS` | Cloud Run | `https://adjusta.vercel.app` |
| Google OAuth client ID / secret | Cloud Run | GCP コンソールで発行 |
| `CLOUD_TASKS_PROJECT_ID` | Cloud Run | Cloud Tasks queueを配置するGCP project ID |
| `CLOUD_TASKS_LOCATION` | Cloud Run | Cloud Runと同じ`asia-northeast1` |
| `CLOUD_TASKS_QUEUE` | Cloud Run | `google-calendar-sync` |
| `CLOUD_TASKS_HANDLER_URL` | Cloud Run | `${Cloud Run URL}/internal/tasks/google-calendar-sync` |
| `CLOUD_TASKS_OIDC_AUDIENCE` | Cloud Run | Cloud Run serviceのorigin URL |
| `CLOUD_TASKS_INVOKER_SERVICE_ACCOUNT` | Cloud Run | Cloud TasksがOIDC tokenに使用する専用service account |

### GitHub production environment

`.github/workflows/backend-deploy.yml`はGitHubの`production` environmentを使用する。初回実行前に以下を登録する。

Repository / environment variables:

| 名前 | 内容 |
| --- | --- |
| `GCP_PROJECT_ID` | GCP project ID |
| `GCP_REGION` | Artifact Registry / Cloud Runのregion |
| `GCP_ARTIFACT_REPOSITORY` | Docker repository名 |
| `CLOUD_RUN_SERVICE` | Cloud Run service名 |
| `CLOUD_RUN_SERVICE_ACCOUNT` | Cloud Run実行用service accountのメールアドレス |
| `CLOUD_TASKS_QUEUE` | Cloud Tasks queue名 |
| `CLOUD_TASKS_INVOKER_SERVICE_ACCOUNT` | task handler呼び出し用service accountのメールアドレス |
| `CLOUD_TASKS_HANDLER_URL` | Cloud Run上のtask handlerの完全URL |
| `CLOUD_TASKS_OIDC_AUDIENCE` | task handlerが検証するCloud Run serviceのorigin URL |
| `DATABASE_URL_SECRET` | Secret Manager上のDB接続文字列secret名 |
| `SESSION_SECRET_SECRET` | Secret Manager上のsession secret名 |
| `GOOGLE_CLIENT_SECRET_SECRET` | Secret Manager上のOAuth client secret名 |
| `GOOGLE_CLIENT_ID` | OAuth client ID |
| `GOOGLE_REDIRECT_URI` | Googleへ登録したcallback URL |
| `REDIRECT_URL_AFTER_LOGIN` | ログイン完了後のfrontend URL |
| `CORS_ALLOW_ORIGINS` | 許可するfrontend origin |
| `COOKIE_DOMAIN` | 空のままにし、Vercel originのhost-only cookieとして扱う |

Environment secrets:

| 名前 | 内容 |
| --- | --- |
| `GCP_WORKLOAD_IDENTITY_PROVIDER` | Workload Identity Providerの完全なresource名 |
| `GCP_DEPLOY_SERVICE_ACCOUNT` | GitHub Actionsがimpersonateするservice account |

長期service account keyは置かず、GitHub OIDCとWorkload Identity Federationを使う。deploy用identityにはArtifact Registryへのpush、Cloud Run更新、DB migration用secretの参照、およびCloud Run実行service accountを利用するための最小権限を付与する。Cloud Runの実行service accountには、runtimeで参照するSecret Manager secretへのaccessor権限とCloud Tasksへのenqueue権限が必要になる。Cloud TasksのOIDC tokenにはenqueueを行うruntime identityとは別の専用service accountを使用する。

`production` environmentのdeployment branch policyとWorkload Identity Providerのattribute conditionは、どちらも`main`からの実行だけを許可する。初回本番確認までは`koo-arch`をrequired reviewerとし、個人運用で承認不能にならないようself-reviewは許可する。

### Google OAuth

- 承認済みリダイレクトURIに **frontendのcallback URL**(`https://adjusta.vercel.app/api/auth/google/callback`)を登録
- OAuth 同意画面の公開設定(テストユーザー→本番公開)を確認

### Cloud Run

Productionで使用するGCPリソースは以下とする。

| 項目 | 値 |
| --- | --- |
| Project | `adjusta` |
| Region | `asia-northeast1` |
| Cloud Run service | `adjusta-api` |
| Cloud Run URL | `https://adjusta-api-1001878278191.asia-northeast1.run.app` |
| Artifact Registry | `adjusta` |
| Runtime service account | `adjusta-runtime@adjusta.iam.gserviceaccount.com` |
| Cloud Tasks queue | `google-calendar-sync` |
| Task invoker service account | `adjusta-tasks-invoker@adjusta.iam.gserviceaccount.com` |
| Deploy service account | `adjusta-deploy@adjusta.iam.gserviceaccount.com` |

- `backend/Dockerfile` をそのままビルド(GitHub Actions から `gcloud run deploy` が簡単)
- **min-instances**: 0 だと初回アクセス・OAuth コールバックにコールドスタート遅延が乗る。体感を優先するなら 1(常時課金)。まず 0 で始めて気になったら上げる
- Cloud Runの実行identityには専用service accountを使用し、workflowの`CLOUD_RUN_SERVICE_ACCOUNT`で明示する
- runtime service accountには`roles/cloudtasks.enqueuer`を付与する
- task invoker service accountのOIDC tokenをtask handlerで検証し、想定するissuer、audience、service account email以外を拒否する
- cookie の `Secure` / ドメイン設定が本番 URL 前提になっているか `backend/api/cookie` を確認

初回作成時はVercelのserver-side proxyから到達できるようCloud Run invocationのIAMを設定する。workflowはrevisionのデプロイのみを担当し、公開・非公開のIAM設定は変更しない。

### DB migration

- Ent Schemaをdesired state、`backend/migrations`配下のSQLをスキーマ変更履歴の正本とする
- `atlas.sum`はマイグレーションディレクトリの完全性検証に使用し、手動編集せずAtlasのコマンドで更新する
- 各環境への適用状況はDB上のAtlasリビジョンテーブルで管理する
- アプリケーションテーブルは`adjusta` schemaへ配置し、`public`にはAtlasの履歴テーブル`atlas_schema_revisions`以外のテーブルを置かない
- Neon productionではdatabase scopeの履歴schema自動判定が不整合にならないよう、migration適用時は`--revisions-schema public`を明示する。ローカルPostgreSQLはdatabase scopeのデフォルトである専用schemaを使用する
- ent queryはschema修飾されるため、`DATABASE_URL`へ`search_path`を追加しない
- schema変更時は`cd backend && atlas migrate diff <name> --env local`でmigrationを生成し、SQLをレビューする
- ローカル適用は`docker compose --profile tools run --rm migrate`を使う
- PRではmigration履歴を一時PostgreSQLへ再適用し、ent schemaとの差分が残っていないことを検査する
- productionではbackend imageのbuild / push後、Cloud Run revisionの切り替え前に独立したstepとしてpending migrationを適用する
- 上記の実行順序だけでは互換性は保証されない。Cloud Runのトラフィック分割やロールバックで旧revisionが動作する可能性を前提に、migration後も新旧両方が動作できる状態を維持する
- 削除・rename・型変更・NOT NULL化・外部キー追加など後方互換性を損なう変更は、expand/contractで複数回のreleaseに分ける

初期migrationは空DB向けである。auto migrationで`public`に作成済みの既存開発DBは、不要なデータであればvolumeを作り直してから適用する。保持が必要な場合は`public`から`adjusta`への明示的なデータ移行が必要であり、そのままbaseline扱いにはしない。productionは初回デプロイ前のため、空DBへ通常適用する。

### CI/CD

- frontend: Vercel の Git 連携(main への push で本番、PR ごとに Preview)
- backend CI: PRでGoテストとAtlas migration整合性検査を実行
- backend deploy: `Backend Deploy`を手動実行し、Artifact Registryへpush → DB migration適用 → Cloud Run deployの順で実行
- 初回本番確認が完了するまでは自動deployにせず、GitHub `production` environmentのapprovalを利用する。安定後に`main` push triggerを追加する

## 未決事項(実施時に決める)

1. **独自ドメイン**を取るか(取るなら Vercel に割当。backend はブラウザ非公開のため Cloud Run ドメインのままでよい)
2. **staging 環境**の要否(Vercel Preview + Cloud Run のリビジョンタグで軽く済ませる案が有力)
3. 初回本番確認後、backend deployを`main` pushで自動実行するか
