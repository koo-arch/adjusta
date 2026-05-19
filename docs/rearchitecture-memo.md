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
| 認証方式 | Google OAuth を用い、自前 JWT は原則持たない。認証状態はセッションまたは OAuth ベースで管理する | バックエンドで独自 JWT を発行し、フロントエンドは `access_token` cookie を見て画面遷移を制御している | 認証基盤、middleware、cookie 設計、API 認可方式を見直す必要がある |
| User モデル | `users` は Adjusta 利用者の基本情報のみを保持し、Google アカウント識別情報と token は `accounts` で管理する | `users` に `refresh_token` と `refresh_token_expiry` を保持している | User と Token の責務が混在しており、モデル再定義が必要 |
| OAuth トークン管理 | Google Calendar API 用トークンを安全に保存する | `OAuthToken` テーブルと `users.refresh_token` が併存している | `accounts` に集約する前提でトークン保存先の正規化が必要 |
| カレンダーの関係 | `users` と `calendars` は `user_calendars` を介した多対多 | 現行スキーマは `User -> Calendar` の 1 対多に近い構造で、`user_calendars` が存在しない | 共有カレンダー、用途別カレンダー、メインカレンダー管理の前提を満たせない |
| カレンダー属性の置き場所 | `calendars` が `google_calendar_id`、`summary`、`description`、`timezone` を持つ | `google_calendar_id` と `is_primary` が `google_calendar_infos` にあり、`calendars` は実質 ID のみ | カレンダー集約の設計を見直す必要がある |
| Event 所有者と登録先 | `events` は `user_id`、`primary_calendar_id`、`confirmed_date_id` を持つ | 現行は `calendar` エッジ中心で、`user_id` と `primary_calendar_id` を直接持たない | イベントの所有権、登録先カレンダー、認可判定が曖昧 |
| 確定予定の Google Event ID | 要件書では `Event.confirmed_google_event_id` の保持を想定している | 現行 `events.google_event_id` の意味が候補予定か確定予定か曖昧 | Google Calendar との同期対象が曖昧になりやすい |
| ProposedDate の状態 | `proposed_dates` は `google_event_id`、`status`、`priority` を持ち、`status` は `active` / `confirmed` / `not_selected` / `cancelled` を取る | 現行は `start_time`、`end_time`、`priority` のみで、`status` と `google_event_id` がない | 候補、確定、非選択、取り下げなどの状態管理ができない |
| 日程確定ロジック | 確定日程を明示的に状態遷移させ、非選択候補も識別する | 現行実装は主に `priority` の振り直しで確定を表現している | 状態遷移が DB で表現されず、UI や同期で不整合が起こりやすい |
| 同期状態管理 | Event / ProposedDate に `sync_status`、`last_synced_at`、`last_sync_error` を持たせる | 現行スキーマに同期状態系のカラムがない | Google Calendar 連携失敗時の再試行や表示が難しい |
| バックエンドの層構成 | DDD を意識し、domain は ent に依存しない。repository interface は domain 側に置く | repository interface や usecase 相当が `ent` 型や `*ent.Client` に直接依存している | 再設計時に handler/usecase/repository/infrastructure の分離が必要 |
| フロントエンド連携前提 | API 型の重複を避け、server data と draft state を分離する | 現状は API 契約が現行 backend モデルに強く引っ張られている | backend 再設計後に API 契約と状態責務の整理が必要 |

---

## 4. docs で確定した実装前提

現時点では、以下を実装前提として固定してよい。

### 4.1 認証方式

- 初期実装ではセッション主体で統一する
- `users:accounts = 1:1` とし、1ユーザーにつき1つの Google アカウント連携情報を持つ
- `accounts` は Google 前提の論理設計とし、provider は持たない

### 4.2 カレンダー用途の語彙

初期実装で扱うロールは以下で統一する。

| role | 用途 |
|---|---|
| `primary` | 確定予定の登録先 |
| `adjusta_candidate` | 候補予定の登録先 |
| `reference` | 空き時間判定の参照先 |

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

### Phase 2: 認証基盤の是正

目的:

- docs に合わせて認証方式を統一する

主な作業:

- セッション主体の認証フローへ寄せる
- `accounts` / `sessions` を前提にした永続化と middleware を整理する
- cookie 保存内容と責務を明確にする

成果物:

- 更新済み認証フロー
- 認証関連の API / middleware 実装

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

- `backend/ent/schema/*`
- `backend/internal/repo/*`
- `backend/internal/apps/*` または後継 usecase 層
- `backend/api/handlers/*`
- `backend/api/middlewares/*`
- `backend/internal/auth/*`

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

1. `backend/ent/schema` の現状と docs の差分を洗い出す
2. `accounts` / `sessions` を前提に認証フローと middleware の再設計方針を固める
3. Event / ProposedDate / UserCalendar を中心に repository interface と usecase 入出力の語彙を決める
4. その前提で ent schema の再設計と API の再実装に入る

以上を起点に、Phase 1 から順に実装へ着手する。
