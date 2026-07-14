# Adjusta 再設計メモ

## 1. 目的

本メモは、`docs/requirements.md` と `docs/db-design.md` を正本として実装を再設計するために、現実装との差分と着手順を整理することを目的とする。

本メモ自体は仕様の最終確定版ではなく、再設計の出発点として扱う。

---

## 2. 正本として扱うもの

再設計では、以下の優先順位で判断する。

1. `docs/requirements.md`
2. `docs/db-design.md`
3. `AGENTS.md`
4. 既存実装

既存実装は参照対象ではあるが、仕様の正しさを保証するものとしては扱わない。

また、`docs/requirements.md` と `docs/db-design.md` の記述が衝突する場合は、実装を進める前に docs 側で先に解消する。

---

## 3. 現実装との差分一覧

| 項目 | docs の方針 | 現実装の状況 | 影響 |
|---|---|---|---|
| 認証方式 | Google OAuth を用い、自前 JWT は原則持たない。認証状態はセッションまたは OAuth ベースで管理する | JWT 発行は撤去され、backend は `sessions` と session cookie を中心に認証状態を扱う形へ移行済み。frontend も `GET /api/users/me` と middleware の session 検証を起点に寄っている。HTTP 境界の cookie session 操作は `api/sessionctx.CookieSessionStore` に集約し、OAuth callback / logout の失敗系も session cookie と OAuth state の後始末を行う | Phase 2 の大枠は完了。Google 連携再認可の frontend 表示導線は後続で整える |
| User モデル | `users` は Adjusta 利用者の基本情報のみを保持し、Google アカウント識別情報と token は `accounts` で管理する | `users.refresh_token*` は撤去され、Google アカウント識別情報と token は `accounts` に集約済み | この差分は概ね解消済み。今後は Google 連携状態と Adjusta ログイン状態の扱いをさらに分離する |
| OAuth トークン管理 | Google Calendar API 用トークンを安全に保存する | `OAuthToken` は撤去され、Google token は `accounts.access_token` / `refresh_token` / `expires_at` / `scope` で管理している | 保存先の正規化は概ね解消済み。Google token refresh 失敗時の再認可導線は継続課題 |
| カレンダーの関係 | `users` と `calendars` は `user_calendars` を介した多対多 | `user_calendars` を導入し、`users` と `calendars` の関係は docs に沿った多対多へ移行済み | この差分は概ね解消済み。今後は `user_calendars.role` を前提に用途別ロジックを整理する |
| カレンダー属性の置き場所 | `calendars` が `google_calendar_id`、`summary`、`description`、`timezone` を持つ | `calendars` に Google Calendar 情報を集約し、`google_calendar_infos` は廃止済み | この差分は解消済み。以後 `googlecalendarinfo` を正本として扱わない |
| Event 所有者と登録先 | `events` は `user_id`、`primary_calendar_id`、`confirmed_date_id` を持つ | `events` は docs に沿って `user_id` / `primary_calendar_id` / `confirmed_date_id` を直接保持する形へ移行済み | この差分は概ね解消済み。認可、登録先、同期時の前提が明確になった |
| 確定予定の Google Event ID | 要件書では `Event.confirmed_google_event_id` の保持を想定している | `events.confirmed_google_event_id` を正本に統一し、backend 内部と ent schema から legacy `events.google_event_id` は除去済み。API の `google_event_id` は互換用の派生値として一時的に残る | DB 正本の曖昧さは解消済み。レスポンス契約の簡素化は後続で進める |
| ProposedDate の状態 | `proposed_dates` は `google_event_id`、`status`、`priority` を持ち、`status` は `active` / `confirmed` / `not_selected` / `cancelled` を取る | `proposed_dates` は docs に沿った状態語彙、`google_event_id`、priority を持つ形へ移行済み | この差分は概ね解消済み。状態遷移のさらなる domain 集約は継続課題 |
| 日程確定ロジック | 確定日程を明示的に状態遷移させ、非選択候補も識別する | 確定候補を `confirmed`、非選択候補を `not_selected` として扱う流れへ寄ってきている | docs にかなり近づいたが、状態遷移ルールの domain 集約は継続課題 |
| 同期状態管理 | Event / ProposedDate に `sync_status`、`last_synced_at`、`last_sync_error` を持たせる | `events` / `proposed_dates` の両方に同期状態系カラムを導入済みで、create / update / finalize / detail access sync で更新している | 同期失敗の保持と再同期前提は整ってきた。再試行ポリシーや運用ルールは継続課題 |
| バックエンドの層構成 | DDD を意識し、domain は ent に依存しない。repository interface は domain 側に置く | repository interface の domain 側移設、repository 実装の infrastructure 側集約、usecase ごとの port 分離は進んでいる。transaction は UoW / infrastructure 側に閉じ込め、domain repository から `WithTx(...)` を除去済み | shared model の残りなど、domain の純化は継続課題 |
| フロントエンド連携前提 | API 型の重複を避け、server data と draft state を分離する | event API 型は `status` / `sync_status` / `confirmed_google_event_id` 前提へ追従済み | server data と draft state の責務分離、互換項目の整理は継続課題 |

---

## 4. docs で確定した実装前提

現時点では、以下を実装前提として固定してよい。

### 4.1 認証方式

- 初期実装ではセッション主体で統一する
- `users:accounts = 1:1` とし、1ユーザーにつき1つの Google アカウント連携情報を持つ
- `accounts` は Google 前提の論理設計とし、provider は持たない

#### 4.1.1 認証セッションの正本

- Adjusta のログイン状態は `sessions` を正本とする
- ブラウザが保持する認証情報は、アプリ用のセッション cookie のみとする
- session cookie には Google の access token や refresh token を入れない
- 認証済み判定は「有効な session cookie があり、対応する `sessions` レコードが有効期限内であること」で判断する
- backend の保護 API は、session から user を解決して認可コンテキストを構築する

#### 4.1.2 Google OAuth token の責務

- Google Calendar API 用の token は `accounts` に保存する
- `accounts.access_token` / `refresh_token` / `expires_at` / `scope` は、Google API 呼び出し時にのみ参照する
- Google token の期限切れ時は backend で refresh し、更新結果を `accounts` に書き戻す
- Google token の refresh に失敗した場合は、「Adjusta のログイン状態が無効」ではなく「Google 連携の再認可が必要」として扱う
- したがって、アプリ認証と Google 外部連携状態は分けて扱う

#### 4.1.3 cookie 方針

- cookie は session cookie のみを採用する
- `access_token` cookie や `refresh_token` cookie は採用しない
- session cookie は `HttpOnly` を前提とし、`Secure` は環境ごとに設定する。`SameSite` はfrontendのsame-origin proxyを前提に`Lax`へ固定する
- frontend は cookie の中身を認証の「判定」に使わない。proxy(middleware)は cookie の存在だけを見る楽観的ルーティング(UX)に徹し、権威的な認証判定は Go API のセッション検証に集約する(詳細は 4.1.8)

#### 4.1.4 ログインフローの最終形

1. ユーザーが `/auth/google/login` にアクセスする
2. backend は OAuth state を発行し、必要なら一時的に session に保持する
3. Google callback で code を exchange する
4. Google の userinfo を取得し、`users` と `accounts` を upsert する
5. backend は `sessions` にログインセッションを作成する
6. ブラウザには session cookie のみを返却する
7. frontend は認証済み画面へ遷移する

補足:

- このフローではアプリ独自 JWT を発行しない
- Google OAuth callback は、Google token 保存と Adjusta session 作成を同一 transaction で扱える形を優先する

#### 4.1.5 logout / middleware / 認可の責務

- logout は `sessions` の無効化または削除と session cookie の破棄を行う
- auth middleware は session を検証し、`user_id` と必要最小限の user 情報を context に積む
- auth middleware は Google access token の検証や refresh を担当しない
- Google token の取得と refresh は、Google Calendar を呼ぶ usecase または service 側で行う
- これにより「ログイン確認」と「Google 連携確認」の責務を分離する

#### 4.1.6 Phase 2 完了時に廃止対象とするもの

- `JWTManager`
- `KeyManager`
- `JWTKey` schema / table
- `access_token` cookie / `refresh_token` cookie
- `users.refresh_token`
- `users.refresh_token_expiry`
- `OAuthToken` schema / table / repository
- app JWT の refresh を行う auth middleware
- frontend の「cookie があるか」を認証の判定として扱う実装(cookie 存在を proxy の楽観的ルーティングのヒントに使うことは 4.1.8 のとおり設計内)

#### 4.1.7 現時点の変更の扱い

- `accounts` に Google token を集約する変更は、最終アーキテクチャでも継続利用できる
- JWT ベース実装は撤去済みであり、今後 auth 実装を進める際は session と Google 連携状態の責務分離を前提に変更を入れる
- 既存の session 主体実装は最終形の土台として扱い、HTTP 境界の cookie session 操作は `api/sessionctx.CookieSessionStore` に閉じる
- OAuth callback の失敗系では OAuth state を破棄する
- logout では DB session 削除に失敗しても browser session cookie の破棄を試みる
- Google token refresh / Google API 認可失敗は、Adjusta ログイン失効ではなく Google 連携再認可要求として扱う

#### 4.1.8 frontend の認証境界(2026-07-07 ADR)

frontend の認可は次の責務分離で扱う(CVE-2025-29927 の教訓として、middleware/proxy を認証の唯一の砦にしない)。

- `proxy.ts` は session cookie の**存在チェックのみ**で UX 振り分けを行い、backend への検証 fetch はしない
  - cookie なし + 保護ルート → `/login` へ redirect
  - cookie あり + `/` または `/login` → `/dashboard` へ redirect(要件 5.1.1「認証済みユーザーがログイン画面へ → ダッシュボードへ」に対応)
- 期限切れ cookie のループは `/api/auth/session-expired`(Route Handler)が断ち切る。session cookie を失効(Max-Age=0、backend と同じ domain/path 規則)させて `/login` へ 303 することで、「401 → `/login` → `/dashboard` → 401 → …」の往復が起きる前に cookie そのものが消える。RSC レンダリング中は Set-Cookie できないため、401 の着地はこの handler に集約する
- セキュリティ境界は Go API のみ(全 `/api/*` のセッションミドルウェア)。frontend のチェックはすべて UX 装置であり防御ではない
- Server Component からの認証必須データ取得は DAL(`frontend/src/lib/server/api.ts` の `serverApi` / `requireUser`)に統一し、401 は `redirect("/api/auth/session-expired")` に集約する
- ブラウザ側の 401 / 409 は TanStack Query の QueryCache / MutationCache の onError で一元回収する
  - 401(Adjusta ログイン失効)→ `window.location.assign("/api/auth/session-expired")`(cookie 失効 + フルリロードで認証済みキャッシュを破棄。redirect は一回化)
  - 409(`google_reauthorization_required`)→ 再認可導線のモーダル表示。ログイン失効とは扱いを分ける

route group は layout shell の違いで分割し、認証について次の不変条件を守る。

- `(app)` 以外の route group(`(marketing)`・`(auth)` など)は**認証状態を参照しない**
- ログイン済み/未ログインの出し分けは proxy redirect と `(app)` 側の認証境界でのみ扱う
- `(marketing)`・`(auth)` は `next/headers`・TanStack Query・認証 API 呼び出しを推移的にも import せず、静的レンダリングを維持する

### 4.2 カレンダー用途の語彙

初期実装で扱うロールは以下で統一する。

| role | 用途 |
|---|---|
| `primary` | 確定予定の登録先 |
| `adjusta_candidate` | 候補予定の登録先 |
| `reference` | 空き時間判定の参照先 |

#### 4.2.1 Google Calendar 情報の正本

- Google Calendar の識別子とメタデータは `calendars` に保持する
- ユーザーごとの用途や表示設定は `user_calendars` に保持する
- `googlecalendarinfo` / `google_calendar_infos` は旧設計の名残として扱い、再設計後の正本には含めない
- `is_primary` のようなユーザー依存の意味は、`user_calendars.role = primary` に統一する

### 4.3 ProposedDate の状態語彙

候補日程の状態は以下で統一する。

| status | 意味 |
|---|---|
| `active` | 現在有効な候補 |
| `confirmed` | 確定済み候補 |
| `not_selected` | 確定されなかった候補 |
| `cancelled` | ユーザー操作で取り消した候補 |

### 4.4 確定予定の Google Event ID の置き場

初期実装では `events.confirmed_google_event_id` を持つ。

### 4.5 同期状態の初期スコープ

初期実装では、`events` と `proposed_dates` に以下を持たせる。

- `sync_status`
- `last_synced_at`
- `last_sync_error`

MVP で扱う `sync_status` は、`not_synced` / `pending_sync` / `synced` / `sync_failed` を基本とする。

### 4.6 メール機能の初期スコープ

初期スコープでは「候補日程一覧をコピーできる」を主導線とする。

- 候補日程一覧のコピーは初期導線として扱う
- メール文面作成やテンプレート保存は後続フェーズで拡張する

---

## 5. 目標アーキテクチャ

### 5.1 フロントエンド

- Next.js App Router を前提とする
- 画面単位の container と再利用 UI component を分離する
- server data、form/draft state、UI state を分離する
- backend の API 契約と vocabulary に合わせる

### 5.2 バックエンド

- handler は HTTP 入出力と validation 済み入力の受け渡しに集中する
- usecase は認可、transaction、Google Calendar 連携、状態遷移を orchestration する
- domain model は ent 生成コードに依存しない
- repository interface は domain/usecase 側の語彙で定義する
- ent は infrastructure/persistence 実装に閉じ込める

#### 5.2.1 エラー処理の境界

- usecase / infrastructure は `gin` や `net/http` を知らず、application error を返す
- handler / middleware は、下位層から受け取った error を HTTP に変換するとき `respond.Error` を使う
- query / path / body / session 保存失敗など、HTTP 境界でその場で種類が決まる失敗は handler / middleware で `respond.BadRequest` などを使う
- これにより「アプリケーション上の失敗」と「HTTP 入出力上の失敗」を分離する

### 5.3 永続化

初期案として、以下のテーブル構成を目標とする。

- `users`
- `accounts`
- `sessions`
- `calendars`
- `user_calendars`
- `events`
- `proposed_dates`

必要に応じて後続で以下を拡張対象とする。

- `email_templates`
- `operation_logs`
- `sync_jobs` または同等の補助テーブル

### 5.4 現在の到達度（2026-07-07 時点）

現時点では、バックエンドの層構成は docs の目標形に概ね近づいており、`events` / `proposed_dates` の schema と同期語彙もかなり docs に寄ってきている。auth は session 主体の実装へ移行し、composition root は `cmd/server` と `internal/app` に分かれた。migration はAtlasによるversioned migrationへ移行し、開発composeとデプロイパイプラインから同じ履歴を明示的に適用する。

Phase 3 以降についても、backend は未着手というより大枠の再編成が進んだ状態である。今後の中心は、残る shared model 依存の削減、domain の純化、API DTO と usecase DTO の境界仕上げ、単なる repository 操作の言い換えになっている port / adapter / Record の整理である。frontend は 2026-07-07 の認証アーキテクチャ再編(4.1.8)で、proxy の cookie 存在チェックへの一本化、401 と `google_reauthorization_required` の扱い分け、再認可表示導線の一元化まで進んだ。残る改善余地は server data と draft state の責務分離が中心である。

#### 5.4.1 バックエンド層ごとの到達度

| 層 | 現在の主な配置 | 到達度の目安 | できていること | 主な残課題 |
|---|---|---:|---|---|
| interface | `backend/api/handlers` `backend/api/dto` `backend/api/middlewares` `backend/api/queryparser` `backend/api/requestctx` `backend/api/respond` `backend/api/sessionctx` `backend/api/validation` | 88% 前後 | HTTP 入出力、validation、request context、HTTP error 変換が概ねこの層に集約されている。events の request DTO は `api/dto` へ寄せ始め、`api.Server` と root `api` port package は廃止した。cookie session 操作は `api/sessionctx.CookieSessionStore` に集約した | API response DTO と usecase output の分離は未完 |
| application | `backend/internal/usecase` `backend/internal/google` `backend/internal/errors` | 82% 前後 | usecase ごとの port 分離、transaction orchestration、Google Calendar 連携の orchestration、イベント詳細アクセス時の候補予定再同期が進んでいる。events port では候補日程作成時の selected date DTO を usecase 側へ切り出し始めた | Google 連携の共通語彙を `internal/google` に寄せ始めたが、usecase 専用 DTO / domain model への分離は継続課題 |
| domain | `backend/internal/domain` | 78% 前後 | repository interface の移動、priority / confirm ルール抽出、`events.confirmed_google_event_id` 正本化、transaction 技術要素の domain interface からの除去、ProposedDate repository の create/update option 分離、shared value の `domain/value` 配置が進んでいる | 状態遷移や同期方針の一部はまだ usecase 側にある |
| infrastructure | `backend/internal/infrastructure` | 90% 前後 | repository 実装、UoW、Google Calendar adapter、auth/calendar/events adapter、ent schema の docs 寄せが進んでいる。tx 付き repository の組み立ても infrastructure に集約され、DB接続は`internal/infrastructure/database`、migration履歴は`backend/migrations`へ分離した | versioned migration運用では、破壊的変更をexpand/contractで段階適用する規約の継続が必要 |

#### 5.4.2 現在の補助的な位置づけ

- `backend/cmd/server` は HTTP server 起動のentrypoint、DB migrationはAtlas CLIと`backend/migrations`を正本として扱う
- `backend/internal/app` は composition root の組み立て、router 登録、HTTP server lifecycle を扱う
- handler / middleware 用 port は、それぞれ `backend/api/handlers` / `backend/api/middlewares` に閉じる
- HTTP 境界の cookie session 操作は `backend/api/sessionctx.CookieSessionStore` に閉じる
- `backend/internal/google` は Google 連携の共通語彙として扱い、大きくなりすぎる場合は usecase ごとの contract へ分ける
- Google 連携の実装は `internal/infrastructure/googleoauth` / `internal/infrastructure/googlecalendar` へ寄っている
- `cookie` は HTTP 境界の `api/cookie`、cache は infrastructure、config は `internal/config` として整理済み。環境変数供給経路は DB 系を root `.env` + compose、アプリ固有設定を `backend/.env` + godotenv で扱っている
- `EventTxStore` や usecase 専用 Record は、repository interface 移設前の保護層として有効だったが、domain repository interface が整ってきたため薄くする対象として扱う

#### 5.4.3 この時点で解消できた差分

- repository interface を `internal/domain/*` 側へ移し、repository 実装を `internal/infrastructure/repository/*` へ集約した
- `ent` 依存は repository 実装と composition root へかなり閉じ込められた
- usecase ごとの port 分離を進め、events / auth / calendar の orchestration は usecase に寄せた
- repository interface から transaction 技術要素を外し、tx 付き repository の組み立ては UoW / infrastructure に集約した
- auth callback / middleware / logout は session token を中心に扱う流れへ寄ってきている
- `respond.Error` と application error の境界を整理し、validation error も `APIError` に統一した
- `events` schema は `user_id` / `primary_calendar_id` / `confirmed_google_event_id` / sync 系カラムを持つ docs 寄りの形へ移行した
- backend 内部と ent schema から legacy `events.google_event_id` を除去し、確定予定の Google Event ID は `confirmed_google_event_id` を正本にした
- ProposedDate repository の create/update option を分離し、domain repository から selected date 由来の appmodel 依存を除去した
- events port の候補日程作成入力を usecase DTO に寄せ、infrastructure events adapter から selected date の appmodel 依存を除去した
- calendar sync port の Google Calendar list を usecase DTO に寄せ、Google Calendar infrastructure から API DTO 依存を一部除去した
- events の確定処理では `ConfirmationRequest` を導入し、更新処理が API DTO を組み立て直す流れを除去した
- events の作成・更新処理では `DraftCreationRequest` / `DraftUpdateRequest` を導入し、handler で API DTO から usecase input へ変換する境界を作り始めた
- events の API request DTO を `backend/api/dto` へ移し、validation / handler が interface 層の DTO を扱う形に寄せた
- events の draft / upcoming / needs-action response は usecase output と API response DTO に分け、`appmodel` からイベント API 入出力の責務をさらに剥がした
- Google Calendar event list response も usecase output と API response DTO に分け、`appmodel/event.go` を削除した
- events usecase の port 周辺型を `ports.go` / `records.go` / `requests.go` / `outputs.go` / `mutations.go` / `query.go` に分け、interface と DTO / Record / mutation の見通しを改善した
- events adapter では tx store が reader を再利用する形にし、transaction scope 用 store の read 系重複を削減した
- events usecase の mutation 型は repository update option の alias に寄せ、adapter 側の薄い詰め替えを削減した
- events usecase の `EventRecord` / `ProposedDateRecord` は domain model alias に寄せ、`EventReader` / `EventTxStore` / `internal/infrastructure/events` package を削除した。usecase は domain repository interface の bundle を直接扱い、transaction は tx scope の repository bundle を渡すだけにした
- events usecase の test fake repository を専用ファイルへ切り出し、`CalendarRecord` / `EventTxRepositories` も役割別ファイルに分けた
- calendar sync usecase も `SyncStore` を廃止し、tx scope の domain repository bundle を直接扱う形へ寄せた
- calendar sync usecase の user reader adapter を廃止し、domain user repository を直接受け取る形へ寄せた
- auth usecase も `SignInReader` / `SignInStore` / `SessionStore` を廃止し、domain repository bundle と transaction adapter を直接扱う形へ寄せた
- auth usecase の `service_ports.go` を廃止し、repository bundle、transaction port、OAuth port、output、sign-in plan / mutation builder を役割別ファイルへ分けた
- OAuth handler / auth middleware / session middleware の session cookie 操作を `api/sessionctx` に集約し、HTTP 境界での session 操作の重複を削減した
- `internal/infrastructure/cookie` から Gin 依存を外し、cookie の HTTP response 書き込みは `api/sessionctx` 側で扱う形にした
- OAuth callback 成功時の state 削除と session token 保存を `api/sessionctx.CompleteOAuthSignIn` にまとめた
- `api/sessionctx` は `CookieSessionStore` に集約し、OAuth state 発行・取得、session token 取得、OAuth sign-in 完了、session renewal / clear を handler / middleware から直接扱える境界にした
- 未使用だった gin request context への session token 書き込みを削除し、cookie session store と request context の責務混在を解消した
- OAuth callback で Google から `error` が返った場合、`code` が欠落した場合、state が不正な場合は、usecase を呼ばずに Bad Request とし、保存済み OAuth state を破棄する形にした
- logout では DB session 削除後に cookie を破棄するだけでなく、DB session 削除に失敗した場合も browser session cookie の破棄を試みる形にした
- Google token refresh / Google API 401・403 は `google_reauthorization_required` として扱い、Adjusta のログイン失効を表す 401 と分離した
- auth middleware は `AuthenticateSession` の application error を `respond.Error` へ通し、内部エラーを 401 に潰さない形へ寄せた
- auth middleware の user context 書き込みを `api/requestctx` に寄せ、request context key の直書きを閉じた
- `internal/appmodel` に残っていた Google token / user profile 型を `internal/google` へ移し、Google 連携の共通語彙として扱う方針へ寄せた
- shared domain value を `internal/domain/value` へ移し、event / proposed date / sync / user calendar の役割別ファイルへ分けた
- account / user / OAuth / event handler は `api.Server` 経由ではなく必要な usecase port を直接受け取る形へ寄せ始め、旧 `Handler` wrapper を削除した
- `api.Server` と root `api` port package を廃止し、handler / middleware はそれぞれの package 内 port と専用 dependencies を使う形へ分けた
- middleware 共通の依存束ねを廃止し、auth / calendar / session middleware が必要な依存だけを直接受け取る形へ分けた
- handler / middleware から呼び出す application 入口の port 命名は `XxxService` ではなく `XxxUsecase` に寄せた
- event handler が呼び出す port は `GoogleEventUsecase` / `EventQueryUsecase` / `EventCommandUsecase` に分け、外部 API 取得・読み取り系・状態変更系の入口を分離した
- events usecase は操作単位の `*_usecase.go` / 細かい query ファイルから、`draft` / `detail` / `schedule` / `confirmation` / `google_calendar` の機能単位ファイルへ寄せ始めた
- イベント詳細アクセス時に、`sync_proposed_dates` と `adjusta_candidate` カレンダーを見て候補予定を再同期する流れを実装した
- frontend 側の event API 型は、`status` / `sync_status` / `confirmed_google_event_id` を含めて backend 契約に近づけた
- frontend の認証判定は、`authAtom` / `api/auth/cookie` ではなく `GET /api/users/me` と middleware 上の session 検証結果を起点にする形へ寄せた
- frontend の認証境界を 4.1.8 の形へ再編した。proxy は cookie 存在チェックのみの楽観的ルーティングに簡素化し、middleware からの backend 検証 fetch を廃止した
- `app/` を `(marketing)` / `(auth)` / `(app)` の route group に分割し、`(app)` 以外は認証状態を参照しない不変条件を導入した
- Server Component 用 DAL(`lib/server/api.ts` の `serverApi` / `requireUser`)を導入し、401 → `/login` redirect を DAL に集約した。ヘッダーのユーザーメニューは Suspense + async server component で描画する形にした
- ブラウザ側の 401 / 409 は QueryCache / MutationCache の onError に一元化した。401 はフルリロードで `/login` へ、409(`google_reauthorization_required`)は再認可モーダル表示に分離し、apiClient 組み込みの 401 副作用は削除した
- `backend/main.go` を廃止し、HTTP server entrypoint を `backend/cmd/server`、composition root を `backend/internal/app` に分けた
- `internal/app.Run`はgraceful shutdownを扱い、DB接続は`internal/infrastructure/database`に寄せた
- server起動時のauto migrationと`cmd/migrate`を廃止し、Atlasのversioned migrationを明示的に適用する形へ移行した
- `backend/migrations`のSQLをスキーマ変更履歴の正本とし、`atlas.sum`でディレクトリの完全性を検証し、環境ごとの適用状況はAtlasリビジョンテーブルで管理する形にした
- アプリケーションテーブルを`adjusta` schemaへ分離し、entの`sql/schemaconfig`でruntime queryもschema修飾する形にした
- `docker-compose.yml` では `migrate` service を profile 付きで追加し、`db` healthcheck と `depends_on.condition: service_healthy` により migration 実行前に Postgres 起動完了を待つ形にした
- DB 接続文字列は root `.env` の `DB_USER` / `DB_PASSWORD` / `DB_NAME` から compose の `x-db-env` で組み立て、backend / migrate が同じ `DATABASE_URL` を参照する形にした
- `user_calendars` の部分 unique index に明示的な storage key を付け、ent auto migration 時に `usercalendar_user_id` index 名が衝突する問題を解消した

#### 5.4.4 主な残課題

- **解消済み(2026-07-13)**: ent auto migrationを廃止し、Atlasによるversioned migrationと初期migrationへ移行した
- auth の Phase 2 は大枠完了。session 主体の基盤、HTTP 境界の cookie session store、OAuth callback / logout 失敗系の処理、Google 連携再認可を Adjusta ログイン失効と分ける backend error kind は整った。frontend 側の再認可表示導線と 401 / `google_reauthorization_required` の扱い分けも 4.1.8 の形で完了した
- **解消済み(2026-07-13)**: backend のエラーボディに機械可読な `code`を追加し、frontend の`google_reauthorization_required`判定をHTTP 409依存からcode判定へ移行した
- Google 連携の共通語彙を置いた `backend/internal/google` が大きくなりすぎないか確認し、必要なら usecase ごとの contract へさらに分ける
- 残る usecase 専用 store / adapter を見直し、単なる repository 操作の言い換えになっているものは domain repository interface を直接扱える形へ寄せる
- domain model とほぼ同じ usecase Record を見直し、必要なものだけ残す
- proposed date / event の状態遷移ルールや同期方針を、usecase から domain へさらに引き上げる
- `backend/internal/usecase/events` は draft / confirmation / sync / google などの関心が増えているため、まず同一 package 内でファイル prefix による責務整理を進める。依存関係が十分薄くなった段階で、サブドメイン単位の package 分割も検討する
- 環境変数供給経路は DB 系とアプリ固有設定で分かれているため、必要になれば root `.env` + compose にさらに寄せる
- frontend では API server data と draft state の責務分離をさらに進める(再認可表示導線と middleware の整理は 4.1.8 で完了)

#### 5.4.5 次の作業候補

1. interface 層で API DTO と usecase input / output の境界を仕上げ、互換レスポンス項目や HTTP DTO が usecase 側へ漏れない形へ寄せる
2. frontend の API server data と draft state の責務分離を進める
3. 残る過剰 adapter を薄くし、transaction callback で tx scope の domain repository bundle を扱う形へ寄せる
4. domain rule と usecase orchestration の境界を再確認する
5. `events` usecase は draft / confirmation / sync / google のファイル責務を明確にし、将来的なサブドメイン package 分割に備える
6. versioned migrationはアプリの新旧revisionが並行しても成立するexpand/contractを基本とし、破壊的変更を単一デプロイで行わない
7. Google 連携状態 API の責務を整理する。現実装の `GET /api/account/list` は accounts の一覧や `expires_at` / `scope` ではなく Google プロフィールを返し、`GET /api/users/me` と責務が重複している。1 ユーザー 1 Google 連携の現状では `/account/list` という命名も実態と合わないため、廃止・改名または連携状態専用 API(`provider` / `email` / `status` など必要最小限の DTO)への再定義を判断し、frontend の `fetchAccount` / `useAccounts` と API テストを合わせて更新する。access token の `expires_at` は refresh 可能なため、単独で再認可要否を判定しない

---

## 6. 推奨する実装着手順

### Phase 0: docs の整合性確定

目的:

- `requirements.md`、`db-design.md`、ER 図の整合を確認する
- 実装前提として使う語彙とスコープを固定する

成果物:

- 更新済み `docs/requirements.md`
- 更新済み `docs/db-design.md`
- 更新済み ER 図

補足:

- 現時点では、Phase 0 は概ね完了している前提でよい

### Phase 1: ドメインモデルと DB スキーマの再定義

目的:

- `users`、`accounts`、`sessions`、`calendars`、`user_calendars`、`events`、`proposed_dates` を docs に沿って定義し直す

主な作業:

- ent schema の再設計
- enum と sync_status 関連カラムの反映
- unique index / partial unique index の方針確定
- 既存データ移行方針の検討

成果物:

- 更新済み ent schema
- migration 方針メモ

補足:

- `events` / `proposed_dates` を中心に ent schema の docs 寄せはかなり進んでいる
- server起動時のauto migrationは廃止済みであり、Atlas CLIからversioned migrationを明示的に適用する
- ent schema変更時はmigrationを生成してSQLをレビューし、`atlas.sum`は手動編集せずAtlasのコマンドで更新する

### Phase 2: 認証基盤の是正

目的:

- docs に合わせて認証方式を統一する

主な作業:

- session 作成 / 検証 / logout の usecase と repository interface を定義する
- OAuth callback で `users` / `accounts` / `sessions` を扱うフローへ寄せる
- auth middleware を session 検証専用に差し替える
- Google token refresh を Google Calendar 利用側の service / usecase に寄せる
- `JWTManager` / `KeyManager` / `JWTKey` / `OAuthToken` / `users.refresh_token*` の依存を順に除去する
- frontend の認証判定を session ベースに寄せる

成果物:

- 更新済み認証フロー
- 認証関連の API / middleware 実装
- JWT 非依存化された認証基盤

補足:

- backend は session 主体の認証基盤、HTTP 境界の cookie session store、OAuth callback / logout 失敗系の後始末、Google 連携再認可 error kind まで整っており、Phase 2 は大枠完了として扱う
- frontend は session 検証に寄っているが、Google 連携再認可要求の表示導線と、middleware の cookie presence fallback の整理が残る

### Phase 3: backend の層構成整理

目的:

- handler / usecase / repository / infrastructure の責務を分離する

主な作業:

- repository interface から `ent` 依存を外す
- usecase 層で transaction と Google Calendar 連携を扱う
- 認証・認可・validation・外部 API エラーの扱いを整理する

成果物:

- 再編成された backend ディレクトリ構成
- usecase 単位の API 入出力定義

補足:

- repository 実装の infrastructure 側集約、port 分離、adapter 整理はかなり進んでいる
- 今後は shared model 依存の削減、API DTO と usecase DTO の境界仕上げ、domain 純化が中心課題になる

### Phase 4: イベント・候補日程ユースケースの再構築

目的:

- イベント作成、候補追加、候補編集、候補削除、日程確定を新モデルに合わせて実装する

主な作業:

- Event 作成・更新・削除 usecase
- ProposedDate の状態遷移実装
- Google Calendar 候補予定と確定予定の同期処理
- 整合性エラー時のロールバック方針整理

成果物:

- イベント系 API の再実装
- usecase / repository テスト

補足:

- create / update / finalize / detail access sync はかなり docs に寄ってきている
- 候補予定削除時の Google 側方針、状態遷移ルールの domain 集約、残る edge case の整理は継続課題である

### Phase 5: frontend の追従

目的:

- 新しい API 契約と状態モデルに合わせて画面を整理する

主な作業:

- API 型の見直し
- draft state と server data の責務分離
- 一覧、詳細、編集、確定 UI の調整
- 候補一覧コピー機能の導線を整理する

成果物:

- 更新済み frontend 画面
- 主要導線の動作確認

補足:

- event API 型の追従と、`confirmed_google_event_id` 前提の調整は進んでいる
- 認証判定は `GET /api/users/me` と middleware の session 検証へ寄っているが、cookie presence fallback と Google 連携再認可要求の表示導線は改善余地がある
- server data と draft state の責務分離はまだ途中である

### Phase 6: 同期失敗と移行対応

目的:

- Google Calendar 失敗時の扱いと既存データの移行を仕上げる

主な作業:

- 同期失敗表示
- 再試行ポリシー検討
- 旧スキーマからのデータ移行手順整理
- 運用手順書の追加

成果物:

- 移行手順
- 障害時運用メモ

---

## 7. 想定変更範囲

再設計を進める場合、主に以下の領域に変更が入る想定である。

### docs

- `docs/requirements.md`
- `docs/db-design.md`

### backend

- `backend/internal/infrastructure/ent/schema/*`
- `backend/internal/infrastructure/ent/*`
- `backend/internal/domain/*`
- `backend/internal/usecase/*`
- `backend/internal/infrastructure/*`
- `backend/api/handlers/*`
- `backend/api/middlewares/*`
- `backend/api/respond/*`
- `backend/cmd/*`
- `backend/internal/app/*`

### frontend

- `frontend/src/app/*`
- `frontend/src/features/events/*`
- `frontend/src/features/calendar/*`
- `frontend/src/hooks/*`
- `frontend/src/middleware.ts`

---

## 8. 実装時の注意

- docs に未確定事項が残っている状態で backend schema 変更に進まない
- API 契約を固める前に frontend の詳細実装を進めすぎない
- Google Calendar 連携は usecase から service 経由で扱い、画面や handler に漏らさない
- `ent` 型を domain model として使わない
- 既存データがある前提で、削除より移行を優先する

---

## 9. 次の具体アクション

次に着手する候補は以下の順とする。

1. frontend で Google 連携再認可要求の表示導線を整え、401 と `google_reauthorization_required` の扱いを分ける
2. interface 層で API 入出力 DTO と usecase input / output の境界を仕上げ、usecase に HTTP DTO を直接渡さない形へ寄せる
3. frontend の middleware と認証系 data flow を見直し、cookie presence fallback と session 検証の責務を整理する
4. frontend の server data / draft state 整理へ進む
5. domain rule と usecase orchestration の境界を再確認し、状態遷移や同期方針を domain へ寄せる
6. versioned migrationのexpand/contract運用を維持し、DB変更をCloud Runのrevision切り替えと後方互換にする

以上を起点に、Google 連携再認可の frontend 導線、interface DTO と usecase DTO の分離、frontend の server data / draft state 整理、domain 純化を並行して進める。
