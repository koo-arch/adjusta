# Adjusta UI ガイドライン

## 1. 概要

### 1.1 目的

本書は、Adjusta のフロントエンド UI を実装・変更する際の判断基準(デザイントークンとコンポーネント利用規約)を定めることを目的とする。

画面単位の構成・表示内容・導線は `screen-design.md` に定義する。本書は画面横断の「見た目と使い分けの規約」のみを扱う。

### 1.2 位置づけ

- **ビジュアル仕様(色・タイポグラフィ・余白・elevation)の正は `frontend/DESIGN.md`** とする。DESIGN.md は目指す姿を規定するデザイン仕様書であり、本書には転記しない
- コンポーネント個別の見た目・バリエーションの正は **Storybook**(`frontend/.storybook/`、各コンポーネントの `*.stories.tsx`)とする
- 本書が定めるのは「どの UI パターンにどのコンポーネントを使うか」の対応と、実装上の規約のみ
- 実装と規約が乖離している箇所は「4. 既知の課題」に列挙する(`frontend/DESIGN.md` の Migration Notes と併読すること)

### 1.3 前提

| 項目 | 内容 |
| --- | --- |
| CSS | Tailwind CSS v3(`frontend/tailwind.config.ts`、色トークンの独自拡張なし) |
| UI 基盤 | **shadcn/ui(Radix UI + Tailwind)へ段階的移行中**(2026-07-09 決定)。`src/components/ui/` にコピーインし、既存の自作 cva コンポーネントと共存(3.5 参照) |
| バリアント管理 | class-variance-authority(cva)中心。`tailwind-merge` は依存に存在するが現実装では未使用(shadcn 導入時に `cn()` ユーティリティとして標準化される) |
| テーマ | **ライト固定**(2026-07-09 決定。ダークモードは将来拡張とし、`dark:` を新規に書かない) |
| ヘッドレス UI | @headlessui/react(Dialog / Menu / Listbox / Disclosure) |
| アイコン | **lucide-react を標準とする**(2026-07-10 決定。shadcn 標準に合わせる。既存の @heroicons/react 使用箇所は触るファイルから段階的に移行) |
| 関連ドキュメント | `frontend/DESIGN.md`(デザイン仕様の正)、`screen-design.md`(画面構成)、`db-design.md`(enum 定義)、`ui-review.md`(改善バックログ) |

---

## 2. デザイントークン(`frontend/DESIGN.md` へ移管)

色・ステータスカラー・タイポグラフィ・余白・elevation の定義は `frontend/DESIGN.md`(awesome-design-md-jp 形式のデザイン仕様書)に移管した。本章の旧内容は DESIGN.md の 2〜6 章を正とする。

本書側に残す規約は以下のみ:

- **intent と色の対応**: ボタン等の `intent` バリアント(primary / secondary / danger / warning / success / clear)は DESIGN.md「2. Color Palette & Roles」のセマンティックカラーに対応させる。階調は 500(通常)→ 600(hover)→ 700(active)
- **ステータス表示**: enum(EventStatus / ProposedDateStatus / SyncStatus)→ 色の対応は DESIGN.md の Status 表に従う。マッピング関数はドメイン知識として feature 側に置く(3.3 参照)
- **ダークモード**: **ライト固定に決定(2026-07-09)**。`dark:` バリアントを新規に書かない。既存の部分的な `dark:` と ThemeProvider / ThemeButton の扱いは整理対象(4章)

---

## 3. コンポーネント利用規約

### 3.1 UI パターン → 使用コンポーネント

新しい UI を作る前に、以下の既存コンポーネントで実現できないか確認する。shadcn/ui への段階的移行中のため、「現行」と「移行先」を併記する。**新規・再設計するコンポーネントは移行先(shadcn)を使う**。

| UI パターン | 現行(`src/components/`) | 移行先(shadcn / `src/components/ui/`) |
| --- | --- | --- |
| ボタン(テキスト) | `Button`(variant/intent/size) | `button`(intent は variant にマップ) |
| ボタン(アイコンのみ) | `IconButton` | `button`(`variant="ghost"` + `size="icon"`) |
| セグメント切替 | `ToggleButton` | `tabs` または `toggle-group` |
| ON/OFF スイッチ | `ToggleSwitch` | `switch` |
| テキスト入力(1行 / 複数行) | `TextField` / `TextArea` | `input` / `textarea` + `label`(htmlFor 関連付けが標準で解決) |
| 選択(プルダウン) | `DropdownSelect` | `select` |
| 日時選択 | `common/DateTimePicker`(shadcn calendar + popover + time 入力) | **置換済み(2026-07-12)**: react-datepicker は依存ごと削除 |
| モーダル | `Modal` | `dialog` |
| 確認ダイアログ(破壊的操作) | `Modal` | `alert-dialog` |
| コンテキストメニュー | `PopupMenu` | `dropdown-menu` |
| ステータス表示 | `StatusBadge`(色は `frontend/DESIGN.md` の Status 表に従う) | `badge` ベースのカスタム(色ドット+ラベル構成は維持) |
| カード | `Card` | `card` |
| 並び替えリスト | `DraggableList` + `SortableItem` | 存置(dnd-kit。描画は自前のためトークン準拠で実装) |
| 一時通知 | react-toastify | `sonner` |
| スケルトン | (なし) | `skeleton`(新規導入) |

FullCalendar(カレンダー表示)も shadcn の対象外として存置する。**存置するサードパーティ UI はデフォルトテーマのまま使わず、`frontend/DESIGN.md`「4. Third-Party Components」の上書き規定に従ってアプリのトークンに統一する。**

### 3.2 使い分けの規約

- **重要な状態変更**: イベント削除・候補削除などの破壊的操作は `intent="danger"`、日程確定は `intent="success"` のボタンを使い、いずれも `Modal` による確認を必須とする
- **処理結果の通知**: 画面遷移を伴わない成功・失敗は toast で通知する。フォームの入力エラーは toast ではなく該当フィールドの `error` + helperText で表示する
- **認証エラー**: 401 はフルリロード、409(再認可)は `AuthErrorModal`(`screen-design.md` 4.1 参照)。個別画面で独自ハンドリングしない
- **空状態**: 一覧・セクションが 0 件のときは空状態表示と次アクションへの導線を置く(dashboard の `EmptyStateCard` を参考にする)

### 3.3 配置規約(components と features の分担)

| 置き場所 | 役割 | 規約 |
| --- | --- | --- |
| `frontend/src/components/ui/` | shadcn/ui プリミティブ | CLI でコピーインして所有する。ドメイン・API 非依存。業務ロジックを載せない |
| `frontend/src/components/common/` | 複数 feature で使う共通の組み立て | ui/ プリミティブを組み合わせた汎用部品(EmptyState 等)。Storybook ストーリー必須 |
| `frontend/src/components/`(直下・既存) | 旧・自作 cva プリミティブ | **凍結: 新規追加禁止**。shadcn への置換が完了したものから削除する |
| `frontend/src/features/*/components/` | ドメイン固有の組み立て | 共通コンポーネントを import して構成する。enum → ラベル/色のマッピング、意味付けラッパー(DeleteButton 等)はこちらに置く |

- ページ(`src/app/**/page.tsx`)は薄く保ち、レイアウト(コンテナ幅・余白)は feature の container が持つ
- 同型のドメインスタイルロジック(例: ステータス色)が複数 feature で必要になった場合は、feature 内の共有モジュールへの切り出しを検討する

### 3.4 実装規約

- shadcn コンポーネントは CLI(`npx shadcn@latest add <name>`)で `src/components/ui/` に追加し、コードとして所有・改変する(npm 依存にしない)
- `frontend/DESIGN.md` のトークンを CSS 変数(`--primary` 等)として `globals.css` に定義し、shadcn のテーマ変数に接続する。色の直書きよりトークン経由を優先する
- shadcn のラッパー・カスタマイズに業務ロジックを載せない(`frontend/AGENTS.md` 準拠)
- Storybook: ui/ 配下の未改変プリミティブはストーリー任意、カスタマイズしたもの(badge ベースの StatusBadge 等)と common/ 配下は必須
- バリアントは cva で定義し、クラス結合は shadcn 標準の `cn()`(clsx + tailwind-merge)に統一する
- アイコンは shadcn 標準の lucide-react に統一する(2026-07-10 決定。`components.json` の `iconLibrary: "lucide"` により CLI 生成物も lucide になる)。既存の @heroicons/react は新規使用を避け、画面を触る際に段階的に置き換える
- 色は `frontend/DESIGN.md` のセマンティクスに沿って選び、規約外の色をアドホックに追加しない

### 3.5 shadcn/ui 移行規定

段階的移行の進め方(2026-07-09 決定):

1. **再設計する画面から置換する**: `ui-review.md` のバックログで着手する画面・コンポーネントを shadcn ベースで作り直す。触らない画面の既存コンポーネントは無理に置換しない
2. **同一パターンの一括置換は可**: トースト(react-toastify → sonner)のような横断パターンは、画面単位でなく一括で置換してよい
3. **依存の削除タイミング**: @headlessui/react は Modal / PopupMenu / DropdownSelect / Header(Disclosure)の置換完了後に削除。react-toastify は sonner 置換後に削除
4. **旧コンポーネントの削除**: `src/components/` 直下の自作プリミティブは、参照がなくなった時点でストーリーごと削除する
5. **存置サードパーティのテーマ統一**: FullCalendar / Splide は `frontend/DESIGN.md`「4. Third-Party Components」に従いテーマ上書きでトークンに統一する(FullCalendar は適用済み 2026-07-11。react-datepicker は 2026-07-12 に shadcn へ置換済み)。上書き CSS は `globals.css` のサードパーティ節に集約し、コンポーネント内に散らさない

---

## 4. 既知の課題(実装との乖離)

本書の規約と現実装には以下の乖離がある。修正は本書のスコープ外とし、対応時は個別に判断すること。shadcn/ui への置換(3.5)で自然解消するものはその旨を注記する。

1. **未定義トークンの参照**: `ToggleSwitch` の `color: primary/secondary` が `bg-primary-600` 等の未定義クラスを参照しており実質無効(`tailwind.config.ts` に primary/secondary 色は未定義)。→ `switch` への置換 + CSS 変数トークン導入で解消予定
2. **`dark:` 残骸の整理**(ライト固定決定に伴う): 部分的に存在する `dark:` バリアント、未配置の `ThemeButton`、`ThemeProvider` の扱いを整理する(撤去 or 将来のダーク対応まで凍結の判断)
3. ~~**テーマ系統の二重化**~~ **解消済み(2026-07-09)**: 未使用だった `--foreground-rgb` 等を削除し、`globals.css` は DESIGN.md トークンの CSS 変数(shadcn テーマ変数・ライト固定)に一本化した
4. **色マッピングの分散**: intent → 色の定義が `BaseButton` / `IconButton` / `ToggleSwitch` に、ステータス → 色が `EventCard` に、それぞれ個別定義されており一元化されていない
5. **Tailwind 非標準クラスの使用**: 一部コンポーネントが `text-md` を使っているが、Tailwind デフォルトには `text-md` がない。`text-base` へ統一する必要がある。→ `input` / `textarea` への置換時に解消予定
6. **ファイル名タイポ**: `frontend/src/components/DateTimePicker/DateTImePicker.tsx`(大文字位置の誤り)
