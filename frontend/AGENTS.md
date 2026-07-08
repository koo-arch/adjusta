# AGENTS.md

## Scope
- この指示は `frontend` 配下に適用する。

## Architecture
- Next.js App Router を前提とする。
- feature-based directory を基本とする。
- Server Component と Client Component の責務を分ける。
- server data、form/draft state、UI state を混同しない。

## Route Groups
- `src/app/*` は layout shell の違いで route group に分割する（public かどうかでは分けない）。
  - `(marketing)`: 公開LP・訴求ページ。MarketingHeader + MarketingFooter。
  - `(auth)`: 認証導線ページ（`/login` など）。認証状態を参照しない最小シェル。
  - `(app)`: 認証後アプリ。Providers・Header・AuthErrorModal・`requireUser` を使える。
- 不変条件:
  - `(app)` 以外の route group は認証状態を参照しない。
  - ログイン済み/未ログインの出し分けは proxy redirect と `(app)` 側の認証境界でのみ扱う。
  - `(marketing)` と `(auth)` は `next/headers`・TanStack Query・`useAuth`・`requireUser`・認証 API 呼び出しを推移的に import しない（静的レンダリング維持）。
  - cookie 保持者は `/` と `/login` の両方から `/dashboard` へ redirect する。期限切れ cookie は `/api/auth/session-expired`(Route Handler)が失効させて `/login` へ 303 するためループしない。RSC レンダリング中は Set-Cookie できないため、401 の着地はこの handler に集約する。

## Directory Rules
- `src/app/*` にはルーティングとページエントリのみを置く。
- `src/features/<domain>/<feature>/*` を機能単位の実装場所とする。
- `src/features/<...>/api/*` には API 呼び出し関数を置く。
- `src/features/<...>/hooks/*` には `useQuery` / `useMutation` フックを置く。
- `src/features/<...>/containers/*` にはページ接続層を置く。Server Component 主体で扱う。
- `src/features/<...>/components/*` には表示責務中心の UI を置く。
- `src/features/<...>/store/*` には Jotai atoms を置く。
- `src/features/<...>/queryKeys.ts` には TanStack Query の query key 定義を置く。
- `src/components/ui/*` には shadcn/ui ベースの共通 UI を置く。
- `src/components/common/*` には複数 feature で使う共通部品を置く。
- `src/lib/server/*` には Server Component 専用の DAL（`serverApi` / `requireUser`）を置く。cookie 転送と 401 → `/login` redirect はここに集約する。

## Components
- 画面単位の container と再利用可能な UI component を分ける。
- UI component はできるだけ props で制御する。
- business logic を presentation component に入れすぎない。
- 既存の component / style / naming に合わせる。

## Container Rules
- `containers/*` はページ/機能の接続層として扱う。
- `containers` の責務は、ルート引数の受け取りと feature への受け渡しに限定する。
- `containers` では、秘匿不要データの prefetch を行ってよい。
- `containers` では、`HydrationBoundary` による dehydrated state の受け渡しを行ってよい。
- `containers` では、初期 UI 状態を Jotai に橋渡しする薄いラッパーを置いてよい。
- `containers` では mutation を行わない。
- 認証必須データ・ユーザー固有データのサーバー側取得は DAL（`src/lib/server/api.ts` の `serverApi` / `requireUser`）経由に限る。DAL を通さないアドホックなサーバー fetch / prefetch は行わない。
- prefetch する query key は `queryKeys.ts` の builder を使い、`useQuery` 側と一致させる。

## Component Rules
- 1コンポーネント1責務を基本にする。
- データ取得や状態接続を含むロジックは `hooks` または `containers` に置く。
- 表示中心の UI は `components` に置く。
- feature 内再利用は feature 配下の `components` に留める。
- feature 横断で再利用が必要な場合のみ `src/components/common` へ昇格する。
- shadcn/ui のラッパーに業務ロジックを載せすぎない。

## State Management
- 不要な `useEffect` は避ける。
- server data は fetch / server action / query など既存方針に合わせる。
- form 送信方式は React Hook Form 固定にしない。
- `useActionState`, server actions, mutation, event handler などから画面に合うものを選ぶ。
- URL に保持すべき state と component 内 state を区別する。

## useEffect Policy
- `useEffect` は原則使用しない。escape hatch 扱いとする。
- データ取得は TanStack Query または Server Component の prefetch で行う。
- 派生値は render 時の計算で解決し、必要な場合のみ `useMemo` を使う。
- ユーザー操作起点の処理は event handler / action に寄せる。
- `useEffect` を使ってよいのは、DOM API 連携・購読/解除・タイマーなど副作用が不可避な場合のみとする。
- `useEffect` を使う場合は、なぜ不可避かをコードコメントで1行残す。

## Jotai Policy
- Jotai は、複数コンポーネントにまたがるクライアント状態をシンプルに共有・更新するために利用する。
- Jotai は、`useEffect` 依存の状態同期を減らすためにも利用する。
- 状態遷移は atom、特に write-only atom / derived atom で表現し、`useEffect` での後追い同期を避ける。
- フォーム途中状態・UI状態・ステップ遷移などのクライアント状態は Jotai に置く。
- Query 結果から必要な値を UI 都合で保持する場合は、最小限の派生状態のみ atom に持つ。

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
