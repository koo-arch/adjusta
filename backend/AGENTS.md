# apps/api/AGENTS.md

## Scope
- この指示は `backend` 配下に適用する。

## Architecture
- DDD を意識した層構成を維持する。
- `domain` は `infrastructure` に依存しない。
- `domain` は ent 生成コードに依存しない。
- repository interface は domain 側に置く。
- repository implementation は infrastructure 側に置く。
- handler に business logic を書かない。
- usecase は workflow、transaction、authorization、外部API連携の orchestration を担当する。

## ent / Persistence
- ent は infrastructure / persistence の詳細として扱う。
- domain model と ent model を混同しない。
- repository 実装以外から ent client を直接使わない。
- ent 生成ディレクトリ内に repository 実装を置かない。
- 複合PKは原則使用せず、unique index で重複を防ぐ。

## Validation
- API boundary で input validation を行う。
- validation 済みの値を usecase に渡す。
- handler で業務ルールを検証しすぎない。
- 業務ルールは domain または usecase に寄せる。

## Naming
- handler は HTTP 入出力の責務が分かる名前にする。
- usecase は操作単位が分かる名前にする。
- repository interface は domain の語彙で命名する。
- infrastructure 実装名には実装技術が分かる名前を使ってよい。
- Google Calendar 連携処理は calendar / googlecalendar など既存命名に合わせる。

## Error Handling
- 外部APIエラーやDBエラーをそのまま handler へ漏らさない。
- domain / usecase で扱える application error に変換する。
- 認証・認可・validation・外部API失敗を区別する。

## Testing
- usecase と repository を優先してテストする。
- 外部APIは必要に応じて mock / fake を使う。
- 認可、transaction、status 変更、同期失敗系を重点的に確認する。