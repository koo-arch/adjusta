# AGENTS.md

## Scope
- この指示は `frontend` 配下に適用する。

## Architecture
- Next.js App Router を前提とする。
- feature-based directory を基本とする。
- Server Component と Client Component の責務を分ける。
- server data、form/draft state、UI state を混同しない。

## Components
- 画面単位の container と再利用可能な UI component を分ける。
- UI component はできるだけ props で制御する。
- business logic を presentation component に入れすぎない。
- 既存の component / style / naming に合わせる。

## State Management
- 不要な `useEffect` は避ける。
- server data は fetch / server action / query など既存方針に合わせる。
- form 送信方式は React Hook Form 固定にしない。
- `useActionState`, server actions, mutation, event handler などから画面に合うものを選ぶ。
- URL に保持すべき state と component 内 state を区別する。

## Validation
- frontend validation は UX 用とする。
- backend validation を authoritative とする。
- API 型や validation schema を重複定義しない。
- API 入出力型は共有定義または既存の型定義に合わせる。

## Naming
- component は PascalCase。
- hooks は `use*`。
- query key builder は `build*QueryKey` など既存命名に合わせる。
- action / mutation / api 関数は操作内容が分かる名前にする。
- domain vocabulary は backend / docs と揃える。

## Styling
- Tailwind を前提とする。
- 既存の spacing、色、コンポーネント構成に合わせる。
- 無関係なデザイン変更を含めない。

## Testing
- 画面の主要導線、form submit、validation 表示、権限による表示制御を重点的に確認する。
- 実装変更に関係ない snapshot / style 差分を増やさない。