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
- 仕様判断が必要な場合は、要件定義書・DB設計書を参照する。
- 推測で大きな設計変更をしない。

## Commands
- Dev / build / lint / test コマンドは README または実際の構成を確認してから実行する。
- コマンドが未整備の場合は、推測で追記しない。

## Commit / Git
- コミットメッセージは既存履歴に合わせて Conventional Commits 形式にする。
- 形式は `type(scope): summary` を基本とする。
- 例: `refactor(auth): clarify session boundary`、`fix(api): fix query status validation typo`
- `type` は `feat` / `fix` / `refactor` / `docs` / `test` / `chore` などから変更内容に合うものを選ぶ。

## Definition of Done
- 変更理由と影響範囲を説明できること。
- 変更ファイル一覧を示せること。
- 実行した検証コマンドと結果を報告すること。
- 未実行の検証がある場合は理由を明記すること。
