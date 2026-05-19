# Adjusta DB設計書

## 1. 概要

### 1.1 目的

本設計書は、Adjusta におけるデータベース設計を明確化し、実装・保守・機能追加時の基準とすることを目的とする。

Adjusta は、Google Calendar と連携し、日程調整における候補日程の作成、管理、確定、候補日程共有を支援する Web アプリケーションである。

本設計書では、要件定義で定義されたデータ要件、Google Calendar 同期方針、Soft Delete 方針、将来的な拡張余地を踏まえ、主要エンティティ・リレーション・制約・enum を定義する。

### 1.2 対象範囲

本設計書では、以下のデータ構造を対象とする。

* ユーザー
* OAuth アカウント / トークン
* Google Calendar 情報
* ユーザーとカレンダーの関連
* 日程調整イベント
* 候補日程
* Google Calendar 同期状態
* enum 定義
* 制約・インデックス
* Soft Delete 方針

### 1.3 前提

| 項目 | 内容 |
| --- | --- |
| データベース | PostgreSQL |
| ORM | ent |
| 認証 | Google OAuth |
| 外部連携 | Google Calendar API |
| 削除方針 | 主要エンティティは Soft Delete |
| 同期方針 | Adjusta 側の DB を正とし、Google Calendar は同期先として扱う |


---

## 2. 設計方針

### 2.1 Google OAuth を前提としたユーザー管理

Adjusta では、独自のメールアドレス・パスワード認証は持たず、Google OAuth を用いて認証する。

そのため、users テーブルではアプリ利用者としての基本情報を保持し、`google_user_id`、access token、refresh token などは accounts テーブルで管理する。

### 2.2 Google Calendar との同期方針

Adjusta 側で作成した events および proposed_dates は、Adjusta の DB を正とする。

Google Calendar は同期先として扱い、Google Calendar 側で Adjusta 管理下の予定が削除・変更された場合でも、原則として Adjusta 側の状態を優先して再作成または上書きする。

### 2.3 候補予定と確定予定の配置

候補予定と確定予定は Google Calendar 上で表示先を分ける。

* 候補予定は原則として Adjusta 専用カレンダーに作成する
* 確定予定はユーザーのメインカレンダーに本予定として作成する
* 確定後、候補予定は削除ではなくステータス変更して残す方針を基本とする
* 確定対象以外の候補予定も not_selected などの状態として管理する

### 2.4 カレンダー共有を考慮した多対多設計

Google Calendar には、祝日カレンダーや共有カレンダーなど、複数ユーザーが参照する可能性のあるカレンダーが存在する。

そのため、users と calendars は直接 1 対多にせず、user_calendars を介した多対多構造とする。

calendars.google_calendar_id は Google Calendar API 上のカレンダーを一意に識別する値であるため、単独 UNIQUE 制約を設定する。

### 2.5 UserCalendar によるユーザー別カレンダー用途の管理

user_calendars は、ユーザーにとってそのカレンダーがどのような役割を持つかを管理する。

主な用途は以下である。

* 確定予定を登録するメインカレンダー
* 候補予定を表示する Adjusta 専用カレンダー
* 空き時間確認に利用する参照用カレンダー

また、表示対象をユーザー単位で制御できるようにする。

同期対象かどうかは is_sync_target のような boolean では持たず、role によって判断する。

### 2.6 明示的な中間テーブルの採用

user_calendars は、単なる中間テーブルではなく、role、is_visible などの属性を持つ。

そのため、ent の自動生成される中間テーブルではなく、明示的なテーブルとして定義する。

### 2.7 Soft Delete

サービス全体として Soft Delete を採用する。

物理削除ではなく deleted_at による論理削除を基本とし、復元・調査・誤削除対策を可能にする。

ただし、accounts や sessions などの認証関連データについては、利用する認証方式・セキュリティ要件に応じて別途削除方針を定める。

### 2.8 操作ログ・同期ログ

初期実装では、AuditLog / OperationLog / SyncLog のような専用ログテーブルは持たない。

代わりに、events と proposed_dates に以下を持たせる。

* sync_status
* last_synced_at
* last_sync_error

詳細な履歴が必要になった場合は、将来的に SyncLog / AuditLog を追加する。

---

## 3. ER図

![](ER.drawio.svg)

> ※ 図版は更新途中のため、カラム定義・制約の正本は 5章以降の本文とする。

主要な関係は以下の通りである。

| 関係 | 内容 |
| --- | --- |
| users - accounts | 1人のユーザーは1つの Google アカウント連携情報を持つ |
| users - user_calendars - calendars | ユーザーとカレンダーは多対多 |
| users - events | 1人のユーザーは複数のイベントを作成できる |
| calendars - events | カレンダーは複数のイベントに紐づく |
| events - proposed_dates | 1つのイベントは複数の候補日程を持つ |
| events.confirmed_date_id - proposed_dates.id | イベントは確定済み候補日程を1つ参照する |
| events.primary_calendar_id - calendars.id | イベントは確定予定の登録先カレンダーを参照する |


---

## 4. テーブル一覧

| テーブル名 | 概要 |
| --- | --- |
| users | Adjusta の利用ユーザーを管理する |
| accounts | Google アカウント連携情報と token 情報を管理する |
| sessions | ログインセッションを管理する |
| calendars | Google Calendar のカレンダー情報を管理する |
| user_calendars | ユーザーとカレンダーの関連および用途を管理する |
| events | 日程調整イベントを管理する |
| proposed_dates | イベントに対する候補日程を管理する |


---

## 5. テーブル定義

### 5.1 users

Adjusta のユーザー情報を管理する。

| カラム名 | 型 | NULL | 制約 | 説明 |
| --- | --- | --- | --- | --- |
| id | uuid | NO | PK | ユーザーID |
| name | varchar | YES |  | ユーザー名 |
| email | varchar | NO | UNIQUE | メールアドレス |
| avatar_url | text | YES |  | プロフィール画像URL |
| created_at | timestamp | NO |  | 作成日時 |
| updated_at | timestamp | NO |  | 更新日時 |
| deleted_at | timestamp | YES |  | 論理削除日時 |


#### 補足

* パスワードは保持しない。
* Google アカウントとの紐づけは accounts で管理する。
* email はユーザー識別や表示に利用するため UNIQUE とする。

---

### 5.2 accounts

Google アカウント連携情報と Google Calendar API 利用に必要な token 情報を管理する。

| カラム名 | 型 | NULL | 制約 | 説明 |
| --- | --- | --- | --- | --- |
| id | uuid | NO | PK | アカウントID |
| user_id | uuid | NO | FK | users.id |
| google_user_id | varchar | NO | UNIQUE | Google 側のユーザー識別子 |
| access_token | text | YES |  | Google Calendar API 呼び出し用アクセストークン |
| refresh_token | text | YES |  | アクセストークン更新用リフレッシュトークン |
| expires_at | timestamp | YES |  | アクセストークン有効期限 |
| scope | text | YES |  | 許可された OAuth scope |
| created_at | timestamp | NO |  | 作成日時 |
| updated_at | timestamp | NO |  | 更新日時 |


#### 制約

| 制約 | 内容 |
| --- | --- |
| UNIQUE(user_id) | 1ユーザーにつき1つの Google アカウント連携情報を保証する |
| UNIQUE(google_user_id) | 同一 Google アカウントの重複登録を防ぐ |


#### 補足

* 初期実装では 1ユーザーにつき1つの Google アカウント連携情報を持つ。
* token 情報は安全に保存する。
* 論理設計では provider を持たない。認証ライブラリの都合で必要な場合のみ、物理設計で追加を検討する。

---

### 5.3 sessions

ログインセッションを管理する。

| カラム名 | 型 | NULL | 制約 | 説明 |
| --- | --- | --- | --- | --- |
| id | uuid | NO | PK | セッションID |
| user_id | uuid | NO | FK | users.id |
| session_token | varchar | NO | UNIQUE | セッショントークン |
| expires_at | timestamp | NO |  | セッション有効期限 |
| created_at | timestamp | NO |  | 作成日時 |
| updated_at | timestamp | NO |  | 更新日時 |


#### 補足

* 認証状態はセッションまたは OAuth に基づいて管理する。
* 認証ライブラリを利用する場合、実際のテーブル構造はライブラリ仕様に合わせる。

---

### 5.4 calendars

Google Calendar 上のカレンダー情報を管理する。

| カラム名 | 型 | NULL | 制約 | 説明 |
| --- | --- | --- | --- | --- |
| id | uuid | NO | PK | カレンダーID |
| google_calendar_id | varchar | NO | UNIQUE | Google Calendar API 上のカレンダーID |
| summary | varchar | YES |  | カレンダー名 |
| description | text | YES |  | 説明 |
| timezone | varchar | YES |  | タイムゾーン |
| created_at | timestamp | NO |  | 作成日時 |
| updated_at | timestamp | NO |  | 更新日時 |
| deleted_at | timestamp | YES |  | 論理削除日時 |


#### 補足

google_calendar_id は単独 UNIQUE とする。

同一 Google Calendar を Adjusta 内で重複して保持しないためである。

祝日カレンダーや共有カレンダーのように複数ユーザーが同じカレンダーを参照する場合も、calendars には1件のみ保持し、ユーザーとの関連は user_calendars で表現する。

---

### 5.5 user_calendars

ユーザーとカレンダーの関連を管理する中間テーブルである。

| カラム名 | 型 | NULL | 制約 | 説明 |
| --- | --- | --- | --- | --- |
| id | uuid | NO | PK | 関連ID |
| user_id | uuid | NO | FK | users.id |
| calendar_id | uuid | NO | FK | calendars.id |
| role | enum | NO |  | ユーザーにとってのカレンダー用途 |
| is_visible | boolean | NO |  | アプリ上の予定確認で表示対象にするか |
| created_at | timestamp | NO |  | 作成日時 |
| updated_at | timestamp | NO |  | 更新日時 |
| deleted_at | timestamp | YES |  | 論理削除日時 |


#### 制約

| 制約 | 内容 |
| --- | --- |
| UNIQUE(user_id, calendar_id) | 同一ユーザーと同一カレンダーの重複関連を防止する |
| UNIQUE(user_id) WHERE role = ‘adjusta_candidate’ | 1ユーザーにつき Adjusta 専用カレンダーを原則1つにする |
| UNIQUE(user_id) WHERE role = ‘primary’ | 1ユーザーにつきメインカレンダーを原則1つにする |


#### 補足

* role = primary は、確定予定の登録先として扱う。
* role = adjusta_candidate は、Adjusta が自動作成する候補予定用カレンダーを表す。
* role = reference は、空き時間確認・予定確認に利用する参照用カレンダーを表す。
* user_calendars は FK 以外に業務属性を持つため、明示的なテーブルとして定義する。

---

### 5.6 events

日程調整イベントを管理する。

| カラム名 | 型 | NULL | 制約 | 説明 |
| --- | --- | --- | --- | --- |
| id | uuid | NO | PK | イベントID |
| user_id | uuid | NO | FK | 作成ユーザー。users.id を参照する |
| primary_calendar_id | uuid | NO | FK | 確定予定を登録するメインカレンダー。calendars.id を参照する |
| confirmed_date_id | uuid | YES | FK | 確定した候補日程。proposed_dates.id を参照する |
| confirmed_google_event_id | varchar | YES |  | メインカレンダー上に作成した確定予定の Google Calendar Event ID |
| title | varchar | NO |  | イベントタイトル |
| description | text | YES |  | 説明 |
| location | varchar | YES |  | 場所 |
| status | enum | NO |  | イベント状態 |
| sync_status | enum | NO |  | Google Calendar との同期状態 |
| last_synced_at | timestamp | YES |  | 最終同期日時 |
| last_sync_error | text | YES |  | 最後の同期エラー内容 |
| created_at | timestamp | NO |  | 作成日時 |
| updated_at | timestamp | NO |  | 更新日時 |
| deleted_at | timestamp | YES |  | 論理削除日時 |


#### 補足

* primary_calendar_id は、確定予定を登録するカレンダーを表す。
* 候補予定の登録先は、user_calendars.role = adjusta_candidate の Calendar から取得する。
* confirmed_date_id は proposed_dates.id への外部キーとする。
* confirmed_date_id に指定できる proposed_dates は、同じ events.id に属するものに限定する。
* confirmed_google_event_id は、メインカレンダー上に作成した確定予定の Google Calendar Event ID を保持する。

---

### 5.7 proposed_dates

イベントに紐づく候補日程を管理する。

| カラム名 | 型 | NULL | 制約 | 説明 |
| --- | --- | --- | --- | --- |
| id | uuid | NO | PK | 候補日程ID |
| event_id | uuid | NO | FK | events.id |
| google_event_id | varchar | YES |  | Adjusta 専用カレンダー上の Google Calendar Event ID |
| start_time | timestamp | NO |  | 開始日時 |
| end_time | timestamp | NO |  | 終了日時 |
| priority | int | NO |  | 優先順位 |
| status | enum | NO |  | 候補日程の状態 |
| sync_status | enum | NO |  | Google Calendar との同期状態 |
| last_synced_at | timestamp | YES |  | 最終同期日時 |
| last_sync_error | text | YES |  | 最後の同期エラー内容 |
| created_at | timestamp | NO |  | 作成日時 |
| updated_at | timestamp | NO |  | 更新日時 |
| deleted_at | timestamp | YES |  | 論理削除日時 |


#### 補足

* 候補日程は原則として Adjusta 専用カレンダーに登録する。
* google_event_id には、Google Calendar API で作成された候補予定のイベントIDを保持する。
* Google Calendar に候補予定を表示しない設定の場合、google_event_id は NULL になり得る。
* is_finalized は持たず、確定状態は status = confirmed で表現する。
* 日程確定時、確定されなかった候補日程は status = not_selected とする。

---

## 6. enum 定義

### 6.1 UserCalendarRole

| 値 | 説明 |
| --- | --- |
| primary | 確定予定を登録するメインカレンダー |
| adjusta_candidate | 候補予定を表示する Adjusta 専用カレンダー |
| reference | 予定確認・空き時間確認に利用する参照用カレンダー |


### 6.2 EventStatus

| 値 | 説明 |
| --- | --- |
| draft | 作成途中。まだ候補日程を正式に提示していない状態 |
| active | 候補日程を提示中・調整中の状態 |
| confirmed | 日程が確定済みの状態 |
| cancelled | 日程調整自体を中止した状態 |


#### 状態遷移例

```text
draft -> active -> confirmed
active -> cancelled
draft -> cancelled
confirmed -> active
```

confirmed -> active は、確定後に再調整する場合に利用する。

### 6.3 ProposedDateStatus

| 値 | 説明 |
| --- | --- |
| active | 現在有効な候補日程 |
| confirmed | 確定された候補日程 |
| not_selected | 確定時に選択されなかった候補日程 |
| cancelled | ユーザーが取り下げた候補日程 |


#### 補足

* proposed_dates.status = confirmed は、events.confirmed_date_id が指す候補日程と一致する。
* 日程確定時、確定された候補日程は confirmed とする。
* 日程確定時、確定されなかった候補日程は not_selected とする。
* is_finalized は持たず、status = confirmed で確定済み候補を表す。

### 6.4 SyncStatus

Google Calendar との同期状態を表す。

events と proposed_dates で共通の enum として扱う。

| 値 | 説明 |
| --- | --- |
| not_synced | まだ Google Calendar に同期されていない状態 |
| pending_sync | 同期が必要な状態 |
| synced | Google Calendar と同期済みの状態 |
| sync_failed | 同期処理に失敗した状態 |


#### 将来的な追加候補

| 値 | 説明 |
| --- | --- |
| external_missing | Google Calendar 側で削除されている状態 |
| external_modified | Google Calendar 側で変更されている状態 |


## 7. リレーション定義

| 親 | 子 | 関係 | 説明 |
| --- | --- | --- | --- |
| users | accounts | 1:1 | 1人のユーザーは1つの Google アカウント連携情報を持つ |
| users | sessions | 1:N | 1人のユーザーは複数セッションを持ち得る |
| users | events | 1:N | 1人のユーザーは複数イベントを作成できる |
| users | user_calendars | 1:N | 1人のユーザーは複数カレンダーと関連を持つ |
| calendars | user_calendars | 1:N | 1つのカレンダーは複数ユーザーと関連できる |
| calendars | events | 1:N | 1つのカレンダーは複数イベントの確定予定登録先になり得る |
| events | proposed_dates | 1:N | 1つのイベントは複数候補日程を持つ |
| events | proposed_dates | 1:1 | confirmed_date_id により確定日程を参照する |


---

## 8. 制約・インデックス

### 8.1 Unique 制約

| 対象 | 制約 | 目的 |
| --- | --- | --- |
| users.email | UNIQUE | メールアドレスの重複防止 |
| accounts.user_id | UNIQUE | 1ユーザーにつき1つの Google アカウント連携情報を保証する |
| accounts.google_user_id | UNIQUE | 同一 Google アカウントの重複登録防止 |
| sessions.session_token | UNIQUE | セッショントークンの重複防止 |
| calendars.google_calendar_id | UNIQUE | 同一 Google Calendar の重複登録防止 |
| user_calendars(user_id, calendar_id) | UNIQUE | 同一ユーザーと同一カレンダーの重複関連を防止 |
| user_calendars(user_id) WHERE role = ‘adjusta_candidate’ | Partial UNIQUE | 1ユーザーにつき Adjusta 専用カレンダーを原則1つにする |
| user_calendars(user_id) WHERE role = ‘primary’ | Partial UNIQUE | 1ユーザーにつきメインカレンダーを原則1つにする |


#### PostgreSQL 制約例

```sql
CREATE UNIQUE INDEX uq_user_calendar
ON user_calendars (user_id, calendar_id);
CREATE UNIQUE INDEX uq_user_adjusta_candidate_calendar
ON user_calendars (user_id)
WHERE role = 'adjusta_candidate';
CREATE UNIQUE INDEX uq_user_primary_calendar
ON user_calendars (user_id)
WHERE role = 'primary';
```

### 8.2 Index

| 対象 | 目的 |
| --- | --- |
| accounts.user_id | ユーザーに紐づく Google アカウント連携情報取得 |
| sessions.user_id | ユーザーに紐づくセッション取得 |
| user_calendars.user_id | ユーザーごとのカレンダー一覧取得 |
| user_calendars.calendar_id | カレンダーに紐づくユーザー関連取得 |
| user_calendars.role | カレンダー用途による絞り込み |
| events.user_id | ユーザーごとのイベント一覧取得 |
| events.primary_calendar_id | カレンダー単位のイベント検索 |
| events.confirmed_date_id | 確定日程参照 |
| events.status | イベント状態による絞り込み |
| events.sync_status | 同期状態による再同期対象検索 |
| proposed_dates.event_id | イベント詳細表示時の候補日程取得 |
| proposed_dates.start_time | 日程検索 |
| proposed_dates.status | 候補日程状態による絞り込み |
| proposed_dates.sync_status | 同期状態による再同期対象検索 |


---

## 9. 日程確定時の更新ルール

日程確定時は、以下の DB 更新を行う。

| 対象 | 更新内容 |
| --- | --- |
| events.status | confirmed に更新する |
| events.confirmed_date_id | 確定した proposed_dates.id を設定する |
| events.confirmed_google_event_id | メインカレンダー上に作成した確定予定の Google Calendar Event ID を保存する |
| events.sync_status | 同期結果に応じて synced または sync_failed に更新する |
| 確定対象の proposed_dates.status | confirmed に更新する |
| 非確定の proposed_dates.status | not_selected に更新する |
| proposed_dates.sync_status | 候補予定側の同期結果に応じて更新する |


#### 補足

* 確定予定は、候補予定として作成済みの Google Calendar イベントをそのまま本予定にするのではなく、メインカレンダーに確定予定として登録する。
* Adjusta 専用カレンダー上の候補予定は、削除ではなくステータス変更して残す方針を基本とする。
* confirmed_date_id に指定できる proposed_dates は、同じ events.id に属するものに限定する。

---

## 10. Google Calendar 同期方針

### 10.1 基本方針

Adjusta 側で作成したイベントおよび候補日程については、Adjusta の DB を正とし、Google Calendar は同期先として扱う。

同期対象は、Adjusta が作成した Google Calendar 予定に限定する。

Google Calendar 側で作成された通常予定は、Adjusta 側では管理しない。

### 10.2 同期状態の管理

events と proposed_dates には、Google Calendar との同期状態を管理するために以下を持たせる。

| カラム | 内容 |
| --- | --- |
| sync_status | 同期状態 |
| last_synced_at | 最終同期日時 |
| last_sync_error | 最後の同期エラー内容 |


### 10.3 候補予定

候補予定は、原則として Adjusta 専用カレンダーに作成する。

候補予定を Google Calendar に表示しない設定の場合、Google Calendar には候補予定を作成せず、Adjusta 側の DB のみで管理する。

### 10.4 確定予定

確定日程は、ユーザーのメインカレンダーに本予定として作成する。

その Google Calendar Event ID は events.confirmed_google_event_id に保存する。

### 10.5 外部変更への対応

Google Calendar 側で Adjusta 管理下の予定が削除または変更されていた場合でも、原則として Adjusta 側の状態を優先する。

* Google Calendar 側で削除された候補予定は、必要に応じて再作成する
* Google Calendar 側で変更された候補予定は、Adjusta 側の内容で上書きする
* Google Calendar イベント ID が無効になった場合は、再作成して ID を更新する

---

## 11. Soft Delete 方針

各主要テーブルには deleted_at を持たせる。

通常の検索処理では、deleted_at IS NULL のデータのみを対象とする。

物理削除は原則行わない。

#### 対象候補

* events
* proposed_dates
* calendars
* user_calendars

#### 補足

Google Calendar 側の予定削除またはステータス変更と、Adjusta 側の Soft Delete の意味は区別する。

DB 上の proposed_dates を論理削除した場合でも、対応する Google Calendar イベントを削除するか、ステータス変更して残すかはアプリケーション層で制御する。

---

## 12. ent 実装方針

### 12.1 スキーマ配置

ent のスキーマ定義は infrastructure 層に配置する。

```text
backend/
  internal/
    infrastructure/
      persistence/
        ent/
          schema/
            user.go
            account.go
            session.go
            calendar.go
            usercalendar.go
            event.go
            proposeddate.go
```

### 12.2 中間テーブル

user_calendars は FK 以外に role、is_visible などの属性を持つため、ent の自動生成される中間テーブルではなく、明示的なスキーマとして定義する。

### 12.3 複合PK

ent では複合PKを基本的に使用せず、各テーブルに id を持たせる。

重複防止には unique index を用いる。

#### 例：

user_calendars:
- id: uuid
- unique(user_id, calendar_id)
- unique(user_id) where role = 'adjusta_candidate'
- unique(user_id) where role = 'primary'

### 12.4 confirmed_date_id の注意点

events.confirmed_date_id は proposed_dates.id を参照する。

ただし、events と proposed_dates の間で循環参照に近い構造になるため、ent 実装時には以下を検討する。

* FK 制約を DB で張るか、アプリケーション層で整合性を保証するか
* confirmed_date_id が参照する proposed_dates が同じ event_id に属することをどこで検証するか
* 日程確定処理を usecase 層のトランザクション内で実行すること

---

## 13. 今後の拡張余地

以下は初期実装では対象外または優先度を下げるが、将来的な拡張候補とする。

* SyncLog テーブル
* AuditLog / OperationLog テーブル
* メール文面自動生成
* メールテンプレート管理
* Gmail API 連携
* 調整相手用回答フォーム
* 通知・リマインダー
* 複数ユーザー共同編集
* 組織・チーム利用
* 外部カレンダー連携

---

## 14. 未確定事項

以下は今後の実装・要件整理の中で確定する。

| 項目 | 検討内容 |
| --- | --- |
| 認証ライブラリの実テーブル | accounts / sessions の実カラムは採用ライブラリ仕様に合わせる |
| confirmed_date_id の FK 制約 | DB 制約でどこまで保証するか、アプリケーション層で保証するか |
| priority の仕様 | 数値が小さいほど優先度が高いか、大きいほど高いかを統一する |
| 候補予定を Google Calendar に表示しない設定 | DB 上でユーザー設定として持つか、初期実装では固定値にするか |
| external_missing / external_modified | MVP の sync_status に含めるか、将来的な追加にするか |

