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
| UI 基盤 | **shadcn/ui(Radix UI + Tailwind)**。2026-07-14 に旧自作プリミティブからの移行を完了した |
| バリアント管理 | class-variance-authority(cva)と shadcn 標準の `cn()`(clsx + tailwind-merge) |
| テーマ | **ライト固定**(2026-07-09 決定。ダークモードは将来拡張とし、`dark:` を新規に書かない) |
| ヘッドレス UI | shadcn/ui が使用する Radix UI |
| アイコン | **lucide-react を標準とする**(2026-07-14 に @heroicons/react からの移行を完了) |
| 関連ドキュメント | `frontend/DESIGN.md`(デザイン仕様の正)、`screen-design.md`(画面構成)、`db-design.md`(enum 定義)、`ui-review.md`(改善バックログ) |

---

## 2. デザイントークン(`frontend/DESIGN.md` へ移管)

色・ステータスカラー・タイポグラフィ・余白・elevation の定義は `frontend/DESIGN.md`(awesome-design-md-jp 形式のデザイン仕様書)に移管した。本章の旧内容は DESIGN.md の 2〜6 章を正とする。

本書側に残す規約は以下のみ:

- **variant と色の対応**: shadcn コンポーネントの `variant` は DESIGN.md「2. Color Palette & Roles」のセマンティックカラーに対応させる。Primary は 500(通常)→ 600(hover)→ 700(active)の階調を使う
- **ステータス表示**: enum(EventStatus / ProposedDateStatus / SyncStatus)→ 色の対応は DESIGN.md の Status 表に従う。マッピング関数はドメイン知識として feature 側に置く(3.3 参照)
- **ダークモード**: **ライト固定に決定(2026-07-09)**。`dark:` バリアントを新規に書かない。ThemeProvider / ThemeButton / `next-themes` は撤去済み

---

## 3. コンポーネント利用規約

### 3.1 UI パターン → 使用コンポーネント

新しい UI を作る前に、以下の既存コンポーネントで実現できないか確認する。プリミティブは `src/components/ui/`、複数機能で共有する組み立ては `src/components/common/` を使用する。

| UI パターン | 使用コンポーネント |
| --- | --- |
| ボタン(テキスト) | `ui/button` |
| ボタン(アイコンのみ) | `ui/button`(`variant="ghost"` + `size="icon"`) |
| セグメント切替 | `ui/tabs` または `toggle-group` |
| ON/OFF スイッチ | `ui/switch` |
| テキスト入力(1行 / 複数行) | `ui/input` / `ui/textarea` + `ui/label` |
| 選択(プルダウン) | `ui/select` |
| 日時選択 | `common/DateTimePicker`(shadcn calendar + popover + time入力) |
| モーダル | `ui/dialog` |
| 確認ダイアログ(破壊的操作) | `ui/alert-dialog` |
| メニュー | `ui/dropdown-menu` |
| ステータス表示 | `common/StatusBadge`(`ui/badge` ベース。色ドット+ラベル) |
| カード | `ui/card` |
| 並び替えリスト | `common/DraggableList`(dnd-kit) |
| 一時通知 | `ui/sonner` |
| スケルトン | `ui/skeleton` |

FullCalendar(カレンダー表示)も shadcn の対象外として存置する。**存置するサードパーティ UI はデフォルトテーマのまま使わず、`frontend/DESIGN.md`「4. Third-Party Components」の上書き規定に従ってアプリのトークンに統一する。**

### 3.2 使い分けの規約

- **重要な状態変更**: イベント削除・候補削除などの破壊的操作は `variant="destructive"` を使い、`AlertDialog` による確認を必須とする。日程確定などの主要操作はラベル付きの Primary ボタンを使う
- **処理結果の通知**: 画面遷移を伴わない成功・失敗は toast で通知する。フォームの入力エラーは toast ではなく該当フィールドの `error` + helperText で表示する
- **認証エラー**: 401 は cookie を失効させてログイン画面へフルリロード、`code = google_reauthorization_required`は Adjusta セッションを維持したまま `AuthErrorModal` から Google 再認可へ進む(`screen-design.md` 4.1 参照)。個別画面で独自ハンドリングしない
- **空状態**: 一覧・セクションが 0 件のときは空状態表示と次アクションへの導線を置く(dashboard の `EmptyStateCard` を参考にする)

### 3.3 配置規約(components と features の分担)

| 置き場所 | 役割 | 規約 |
| --- | --- | --- |
| `frontend/src/components/ui/` | shadcn/ui プリミティブ | CLI でコピーインして所有する。ドメイン・API 非依存。業務ロジックを載せない |
| `frontend/src/components/common/` | 複数 feature で使う共通の組み立て | ui/ プリミティブを組み合わせた汎用部品(EmptyState 等)。Storybook ストーリー必須 |
| `frontend/src/components/layout/` | Header などアプリ共通のレイアウト部品 | route group ごとのシェルを構成する部品を置く |
| `frontend/src/features/*/components/` | ドメイン固有の組み立て | 共通コンポーネントを import して構成する。enum → ラベル/色のマッピング、意味付けラッパー(DeleteButton 等)はこちらに置く |

- ページ(`src/app/**/page.tsx`)は薄く保ち、レイアウト(コンテナ幅・余白)は feature の container が持つ
- 同型のドメインスタイルロジック(例: ステータス色)が複数 feature で必要になった場合は、feature 内の共有モジュールへの切り出しを検討する

### 3.4 実装規約

- shadcn コンポーネントは CLI(**`npx shadcn@2.3.0 add <name>`**。Tailwind v3 のため 2.x 系を使う)で `src/components/ui/` に追加し、コードとして所有する(npm 依存にしない)
- **ui/ プリミティブは CLI 生成物のまま保つ**。改変はテーマ上の必然があるときだけ最小限に(例: tabs のアクティブ背景を `bg-card` に)。`add --overwrite` は button 等の**カスタマイズ済みファイルまで上書きする**ため使わず、add 実行後は必ず `git status` / `git diff` で意図しない上書きを確認する
- **制御ロジック・見た目の組み立ては primitive をいじらず別コンポーネントに分離**し、複数 feature で使うものは `src/components/common/` に置く(参照実装: `common/pagination/PaginationControls`、`common/DateTimePicker`)
- `frontend/DESIGN.md` のトークンを CSS 変数(`--primary` 等)として `globals.css` に定義し、shadcn のテーマ変数に接続する。色の直書きよりトークン経由を優先する
- shadcn のラッパー・カスタマイズに業務ロジックを載せない(`frontend/AGENTS.md` 準拠)
- Storybook: ui/ 配下の未改変プリミティブはストーリー任意、カスタマイズしたもの(badge ベースの StatusBadge 等)と common/ 配下は必須
- バリアントは cva で定義し、クラス結合は shadcn 標準の `cn()`(clsx + tailwind-merge)に統一する
- アイコンは shadcn 標準の lucide-react に統一する(`components.json` の `iconLibrary: "lucide"` により CLI 生成物も lucide になる)
- 色は `frontend/DESIGN.md` のセマンティクスに沿って選び、規約外の色をアドホックに追加しない

### 3.5 shadcn/ui 移行規定

移行時に採用した規定(2026-07-09 決定、2026-07-14 移行完了):

1. **再設計する画面から置換する**: `ui-review.md` のバックログで着手する画面・コンポーネントを shadcn ベースで作り直す。触らない画面の既存コンポーネントは無理に置換しない
2. **同一パターンの一括置換は可**: トースト(react-toastify → sonner)のような横断パターンは、画面単位でなく一括で置換してよい
3. **依存の削除タイミング**: @headlessui/react は Modal / PopupMenu / DropdownSelect / Header(Disclosure)の置換完了後に削除。react-toastify は sonner 置換後に削除
4. **旧コンポーネントの削除**: `src/components/` 直下の自作プリミティブは、参照がなくなった時点でストーリーごと削除する
5. **存置サードパーティのテーマ統一**: FullCalendar / Splide は `frontend/DESIGN.md`「4. Third-Party Components」に従いテーマ上書きでトークンに統一する(FullCalendar は適用済み 2026-07-11。react-datepicker は 2026-07-12 に shadcn へ置換済み)。上書き CSS は `globals.css` のサードパーティ節に集約し、コンポーネント内に散らさない

---

## 4. 既知の課題(実装との乖離)

本書の規約と現実装には以下の乖離がある。修正は本書のスコープ外とし、対応時は個別に判断すること。

1. ~~**未定義トークンの参照**~~ **解消済み(2026-07-14)**: 旧 `ToggleSwitch` を削除し、CSS変数へ接続した shadcn `switch` に統一
2. ~~**テーマ有効化経路の整理**~~ **解消済み(2026-07-14)**: ThemeButton / ThemeProvider / `next-themes` と旧プリミティブの `dark:` を撤去
3. ~~**テーマ系統の二重化**~~ **解消済み(2026-07-09)**: 未使用だった `--foreground-rgb` 等を削除し、`globals.css` は DESIGN.md トークンの CSS 変数(shadcn テーマ変数・ライト固定)に一本化した
4. ~~**色マッピングの分散**~~ **解消済み(2026-07-14)**: ボタンは shadcn variant、ステータスの enum 対応は feature、表示は `common/StatusBadge` に分離
5. ~~**Tailwind 非標準クラスの使用**~~ **解消済み(2026-07-14)**: 該当する旧入力プリミティブを削除
6. ~~**ファイル名タイポ**~~ **解消済み(2026-07-12)**: `common/DateTimePicker/DateTimePicker.tsx` へ統一
