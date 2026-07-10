# AGENTS.md

## Scope
- この指示はリポジトリ全体に適用する。
- より深い階層に `AGENTS.md` がある場合は、そちらを優先する。

## Project Overview
- Product: Adjusta
- Frontend: Next.js App Router / TypeScript
- Backend: Go
- Database: PostgreSQL
- ORM: ent
- Auth: Google OAuth
- External API: Google Calendar API

## Required Workflow
- 変更は最小スコープで行い、無関係ファイルは編集しない。
- 命名・import・コードスタイルは既存ファイルに合わせる。
- API 入出力、DB schema、frontend type の変更時は影響範囲を確認する。
- 仕様判断が必要な場合は、Documentation のドキュメントマップと参照ルールに従う。
- 推測で大きな設計変更をしない。

## Documentation

### ドキュメントマップ
- `docs/requirements.md`: 要件定義（機能・画面・API 要件、優先度）。
- `docs/db-design.md`: DB 設計（エンティティ・enum・制約・Soft Delete 方針）。
- `docs/rearchitecture-memo.md`: DDD 再設計の経緯とレイヤー配置の判断理由。
- `docs/screen-design.md`: 画面構成・各画面の表示内容・遷移・現実装との移行課題。
- `frontend/DESIGN.md`: デザイン仕様の正（色・タイポグラフィ・余白・elevation）。
- `docs/ui-guidelines.md`: コンポーネント利用規約・shadcn/ui 移行規定。
- `docs/ui-review.md`: UI 改善バックログ（P1〜P3。UI 実装タスクの出典）。

### 参照ルール
- 機能の追加・変更時は、まず `docs/requirements.md` で対象機能の要件を確認する。
- DB スキーマ・ent の schema 定義に触れる変更では、必ず `docs/db-design.md` を先に読む。
- レイヤー間の依存方向やパッケージ配置に迷ったら `docs/rearchitecture-memo.md` を参照する。
- 画面の構成・表示内容は `docs/screen-design.md` に、UI の見た目・コンポーネント選定は `frontend/DESIGN.md` と `docs/ui-guidelines.md` に従う。
- ドキュメントと実装が食い違っている場合は、勝手にどちらかに合わせず報告する。

## Commands
- Dev / build / lint / test コマンドは README または実際の構成を確認してから実行する。
- コマンドが未整備の場合は、推測で追記しない。

## Commit / Git
- コミットメッセージは既存履歴に合わせて Conventional Commits 形式にする。
- 形式は `type(scope): summary` を基本とする。
- 例: `refactor(auth): clarify session boundary`、`fix(api): fix query status validation typo`
- `type` は `feat` / `fix` / `refactor` / `docs` / `test` / `chore` などから変更内容に合うものを選ぶ。

## Pull Request
- PRタイトルは英語の sentence case とし、先頭を大文字にする。
- PRタイトル例: `Add GitHub CLI to the dev container`
- PR本文は `.github/pull_request_template.md` に従い、英語で書く。
- `Overview` には変更内容の要約を書く。
- `Background` には変更が必要になった理由・背景を書く。
- `Verification` には実行した検証コマンドと結果を書く。
- 未実行の検証がある場合は、その理由を `Notes` または `Verification` に明記する。

## Definition of Done
- 変更理由と影響範囲を説明できること。
- 変更ファイル一覧を示せること。
- 実行した検証コマンドと結果を報告すること。
- 未実行の検証がある場合は理由を明記すること。
