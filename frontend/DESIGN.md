# DESIGN.md — Adjusta

> このファイルは AI エージェントが正確な日本語 UI を生成するためのデザイン仕様書です。
> セクションヘッダーは英語、値の説明は日本語で記述しています。
>
> 本書は Adjusta が**目指すデザイン仕様**を規定します(現状実装の転記ではない)。現実装との差分は末尾「Migration Notes」を参照。
> 実装規約(コンポーネント利用・状態設計)は `../docs/ui-guidelines.md`、画面構成は `../docs/screen-design.md` を参照。

---

## 1. Visual Theme & Atmosphere

- **デザイン方針**: クリーン、効率的、信頼感。日程調整という事務的なタスクを、迷いなく素早く完了できる道具としての UI
- **密度**: 中密度。カード単位で情報をまとめ、余白にゆとりを持たせる。一覧は情報量を絞り、詳細で全情報を見せる
- **キーワード**: 整然 / 落ち着き / 実直 / 軽快 / 邪魔をしない
- **テーマ**: ライトテーマのみ(ダークモードは将来拡張。中途半端な対応はしない)

---

## 2. Color Palette & Roles

<!-- Tailwind CSS v3 デフォルトパレットを採用し、hex に展開して記述 -->

### Primary(ブランドカラー)

- **Primary** (`#6366f1` / indigo-500): メインのブランドカラー。CTA ボタン、フォーカスリング、リンク的操作
- **Primary Dark** (`#4f46e5` / indigo-600): ホバー時
- **Primary Darker** (`#4338ca` / indigo-700): プレス(active)時

### Semantic(意味的な色)

- **Danger** (`#ef4444` / red-500): エラー、削除、破壊的操作
- **Warning** (`#eab308` / yellow-500): 警告、注意喚起
- **Success** (`#22c55e` / green-500): 成功、完了、確定
- **Secondary** (`#ec4899` / pink-500): 補助アクション(使用は最小限に)

ホバー/プレスは Primary と同じ規則で各色の 600 / 700 階調を使う。

### Status(ドメインステータス色)

イベント・候補日程・同期状態のバッジに使う。色+ラベルを必ず併記する(色のみで意味を伝えない)。

| ステータス | 色 | hex |
| --- | --- | --- |
| draft(下書き) | blue | `#3b82f6` |
| active(調整中)/ pending_sync(同期待ち) | yellow | `#eab308` |
| confirmed(確定)/ synced(同期済み) | green | `#22c55e` |
| cancelled(キャンセル)/ sync_failed(同期失敗) | red | `#ef4444` |
| not_selected(非選択)/ not_synced(未同期)/ 不明 | gray | `#6b7280` |

### Neutral(ニュートラル)

- **Text Primary** (`#374151` / gray-700): 本文テキスト
- **Text Strong** (`#111827` / gray-900): ページタイトル等の強調見出し
- **Text Secondary** (`#6b7280` / gray-500): 補足テキスト、ラベル、説明
- **Text Disabled** (`#9ca3af` / gray-400): 無効状態、最弱の補足
- **Border** (`#e5e7eb` / gray-200): 区切り線、カード枠
- **Border Input** (`#d1d5db` / gray-300): 入力欄の枠
- **Background** (`#f9fafb` / gray-50): ページ背景
- **Surface** (`#ffffff`): カード、モーダル、入力欄の面

---

## 3. Typography Rules

### 3.1 和文フォント

- **ゴシック体**: Noto Sans JP(next/font/google で読み込む。ウェイト: 400 / 500 / 700)
- **明朝体**: 使用しない

### 3.2 欧文フォント

- **サンセリフ**: Inter(next/font/google で読み込み済み)
- **等幅**: ui-monospace, SFMono-Regular, Menlo, Consolas(コード・ID 表示が必要な場合のみ)

### 3.3 font-family 指定

```css
/* 本文・UI 全体 */
font-family: Inter, "Noto Sans JP", sans-serif;

/* 等幅(必要な場合のみ) */
font-family: ui-monospace, SFMono-Regular, Menlo, Consolas, monospace;
```

**フォールバックの考え方**:
- 欧文グリフ(日時の数字・英字)は Inter を優先し、和文グリフは Noto Sans JP が受ける混植構成
- 両フォントとも next/font で読み込み、CSS 変数経由で Tailwind の `fontFamily` に登録する
- 和文フォント 1 つだけの指定や、和文未指定(現状)にしない

### 3.4 文字サイズ・ウェイト階層

| Role | Size | Weight | Line Height | Letter Spacing | 備考 |
| --- | --- | --- | --- | --- | --- |
| Display | 24px | 700 | 1.4 | 0 | ページタイトル(h1) |
| Heading 1 | 20px | 700 | 1.4 | 0 | セクション見出し |
| Heading 2 | 18px | 700 | 1.4 | 0 | カード見出し |
| Heading 3 | 16px | 500 | 1.5 | 0 | 小見出し、フォーム label |
| Body | 16px | 400 | 1.7 | 0.02em | 本文 |
| Caption | 14px | 400 | 1.6 | 0.02em | 補足、helperText、日時表示 |
| Small | 12px | 400 | 1.5 | 0.02em | バッジ内ラベル等の最小テキスト |

Tailwind スケール対応: 24px=`text-2xl` / 20px=`text-xl` / 18px=`text-lg` / 16px=`text-base` / 14px=`text-sm` / 12px=`text-xs`。`text-md` は Tailwind に存在しないため使用しない。

### 3.5 行間・字間

- **本文の行間**: 1.7(日本語は欧文より広めが標準)
- **見出しの行間**: 1.4
- **本文の字間**: 0.02em(和文の可読性向上)
- **見出しの字間**: 0
- 日時・数値中心の行(候補日程リスト等)は行間 1.5 まで詰めてよい

### 3.6 禁則処理・改行ルール

```css
overflow-wrap: break-word;  /* 長い URL・英単語の折り返し */
line-break: strict;         /* 厳格な禁則処理 */
```

- イベントタイトル等のユーザー入力テキストは `overflow-wrap: break-word` を必須とする(長文でレイアウトを壊さない)
- 一覧では `truncate`(1行省略)を許容するが、詳細画面では全文表示する

### 3.7 OpenType 機能

```css
/* 見出し・バッジ等の短いテキストのみ許容 */
font-feature-settings: "palt" 1;
```

- **palt**(プロポーショナル字詰め)は見出し・ナビゲーションに限り許容。本文・日時表示には適用しない(数字の桁ズレ防止)

### 3.8 縦書き

該当なし。

---

## 4. Component Stylings

### Buttons

**Primary(solid)**
- Background: `#6366f1` → hover `#4f46e5` → active `#4338ca`
- Text: `#ffffff`
- Padding: 12px 16px(md。sm は 8px 12px)
- Border Radius: 6px
- Font Size: 16px / Font Weight: 500
- 最小高さ 40px(タッチターゲット確保。アイコンのみボタンは 44×44px)

**Secondary(outline)**
- Background: `transparent`
- Text / Border: 対応する intent 色の 500 階調(例: `#6366f1`)
- その他は Primary と同じ

**規則**:
- 画面の主目的となる操作(確定・登録・保存)は必ずラベル付きボタンで置く。アイコンのみにしない
- 破壊的操作(削除)は Danger 色 + 確認モーダル必須
- disabled 時は opacity 50% + `cursor-not-allowed`

### Inputs

- Background: `#ffffff`
- Border: 1px solid `#d1d5db`
- Border (focus): `#6366f1`(ring 2px)
- Border (error): `#ef4444`(ring も red)
- Border Radius: 6px
- Padding: 8px 12px
- Font Size: 16px(iOS のズーム回避のため 16px 未満にしない)
- label は入力と `htmlFor`/`id` で関連付け、エラー時は `aria-invalid` + helperText(`#ef4444`, 14px)

### Cards

- Background: `#ffffff`
- Border: 1px solid `#e5e7eb`
- Border Radius: 8px
- Padding: 16px
- Shadow: Level 1(Depth & Elevation 参照)
- クリック可能なカードはネイティブの button / a 要素でラップし、hover で Shadow を Level 2 に上げる

### Status Badge

- 構成: 色ドット(8px 円)+ ラベルテキスト(14px、ステータス色)
- 色は「2. Color Palette」の Status 表に従う

### Modal

- Surface: `#ffffff`、Border Radius: 12px、Padding: 24px、Max Width: 448px
- Backdrop: `rgba(0, 0, 0, 0.3)`
- Shadow: Level 3

### Third-Party Components

shadcn/ui の置換対象外として存置するサードパーティ UI も、**デフォルトテーマのまま使わず本書のトークン(色・フォント・radius)に合わせて上書きする**(アプリから浮かせない)。フォントはすべて `Inter, "Noto Sans JP", sans-serif` を継承させる。上書き CSS は `globals.css` のサードパーティ節(または専用 CSS ファイル)に集約し、コンポーネント内に散らさない。

**FullCalendar**(v6。`--fc-*` CSS 変数で上書き)

```css
--fc-border-color: #e5e7eb;
--fc-today-bg-color: #eef2ff;          /* indigo-50 */
--fc-button-bg-color: #6366f1;         /* Primary。hover は #4f46e5、active は #4338ca */
--fc-button-border-color: #6366f1;
--fc-highlight-color: rgba(99, 102, 241, 0.15);  /* 日程ドラッグ選択中の範囲 */
--fc-event-bg-color: #6366f1;
--fc-page-bg-color: #ffffff;
```

- ツールバーボタンは radius 6px・階調規則(500→600→700)をアプリの Button 仕様に揃える
- イベントの色: Adjusta 管理の候補・確定予定は Status 表に従う。Google カレンダー由来の予定は Neutral 系で区別する

**react-datepicker** — **置換済み(2026-07-12)**

- shadcn `calendar`(react-day-picker)+ `popover` + time 入力を組んだ共通 `DateTimePicker`(`src/components/common/DateTimePicker/`)に統一し、react-datepicker は依存ごと削除した(ui-guidelines 3.1 更新済み)

**Splide(スライダー)**

- 矢印は Text Secondary `#6b7280`(hover で Primary)、ページネーションドットは Border `#e5e7eb` / 現在位置 Primary に合わせる(既存の globals.css 上書きを拡張)

**sonner(トースト)**

- Surface `#ffffff` / Border `#e5e7eb` / radius 8px / Shadow Level 2 とし、success・error のアクセントに Semantic 色を接続する

---

## 5. Layout Principles

### Spacing Scale

4px 基準(Tailwind スケール)。

| Token | Value | Tailwind |
| --- | --- | --- |
| XS | 4px | `1` |
| S | 8px | `2` |
| M | 16px | `4` |
| L | 24px | `6` |
| XL | 32px | `8` |
| XXL | 48px | `12` |

- コンポーネント内は XS〜M、セクション間は L〜XL、ページブロック間は XXL を目安とする

### Container

ページの最大幅は用途で 3 段階に分ける。画面ごとの適用は `../docs/screen-design.md` の画面定義に合わせる。

| Max Width | 用途 |
| --- | --- |
| 640px(`max-w-screen-sm`) | 単一フォーム画面(ログイン) |
| 768px(`max-w-screen-md`) | リスト中心の画面(一覧、アカウント) |
| 1024px(`max-w-screen-lg`) | 情報量の多い画面(ダッシュボード、作成・詳細・編集) |

- Padding (horizontal): 16px

### Grid

- カード一覧: 1 列(Mobile)/ 2 列(≥640px)/ 3 列(≥1024px)、Gutter 16px
- フォーム+カレンダーの 2 ペイン: ≥768px で左 4 : 右 6、未満は縦積み

---

## 6. Depth & Elevation

### サーフェス戦略(2026-07-10 決定)

**囲い(カード)は、複数の同種要素を並べて区切る必要がある場面でのみ使う。単一の対象を表示する画面では、キャンバス(Background `#f9fafb`)に直置きし、余白とタイポグラフィで構造を作る。**(枠は「ここまでが一塊」という区切りの信号であり、区切る相手がいない画面では機能せず窮屈さだけが残るため)

- **一覧・グリッド系**(区切る相手が多数): 項目ごとの白カード(Surface + Border + Shadow Level 1)
- **詳細など単一対象の画面**: フラット(囲いなし)。セクションが複数あるときは区切り線(Border `#e5e7eb`)を使い、囲いは使わない
- 設定画面のような中間ケースも同じ原則で判断する(セクション複数=区切り線、囲いなし)
- 強調は薄い塗りのパネル(例: 確定日時の green-50)を限定的に使い、面の入れ子はしない

**フラット画面の規律**(囲いがない分、構造は余白だけで示すため以下を固定する):

- 余白の階層: **セクション間は 48px**(区切り線を挟む場合は線の上下に 24px ずつ)/ **セクション内のブロック間は 16px** / 密接な要素(見出しと補足、ラベルと値)は 8〜12px
- コンテンツ幅: 読む画面(詳細等)は **720〜840px** に制限する(`max-w-screen-md` = 768px を標準)。囲いがないと幅の間延びを止めるものがないため必須

| Level | Shadow | 用途 |
| --- | --- | --- |
| 0 | none | フラットな要素、ページ背景上のテキスト、大きな読み物サーフェス(Border のみ) |
| 1 | `0 1px 2px 0 rgba(0,0,0,0.05)` | カード |
| 2 | `0 4px 6px -1px rgba(0,0,0,0.1), 0 2px 4px -2px rgba(0,0,0,0.1)` | ドロップダウン、ポップオーバー、hover 中のカード |
| 3 | `0 20px 25px -5px rgba(0,0,0,0.1), 0 8px 10px -6px rgba(0,0,0,0.1)` | モーダル、ダイアログ |

- 階層は最大 3 まで。装飾目的で影を重ねない

---

## 7. Do's and Don'ts

### Do(推奨)

- font-family は必ず `Inter, "Noto Sans JP", sans-serif` のフォールバックチェーンで指定する
- 日本語本文の line-height は 1.7 を基本とする(最低 1.5)
- 色のコントラスト比は WCAG AA(4.5:1)以上を確保する
- 余白は Spacing Scale に従う
- ステータスは色+ラベルを併記する
- 画面の主要アクションはラベル付きボタンで目立たせる
- データ表示領域は loading(スケルトン)/ error(再試行導線)/ empty(次アクション導線)の3状態を持つ

### Don't(禁止)

- `font-family` に和文フォント 1 つだけ、または欧文フォントのみを指定しない
- 日本語本文に line-height 1.2 以下を使わない
- テキスト色に純粋な `#000000` を使わない(Text Strong は `#111827`)
- 本パレット外の色をアドホックに追加しない(必要なら本書を先に更新する)
- `dark:` バリアントを新規に書かない(ライト固定。ダーク対応は将来まとめて行う)
- 主要操作をアイコンのみのボタンにしない
- ロゴ・アイコン等の画像資産を外部 URL から直接ロードしない(ローカル資産にする)
- レイアウト分岐に JS のメディアクエリ(react-responsive 等)を使わない(CSS ブレークポイントのみ)
- `text-md` 等 Tailwind に存在しないクラスを使わない
- サードパーティ UI(FullCalendar / react-datepicker / Splide 等)をデフォルトテーマのまま使わない(「4. Third-Party Components」に従って上書きする)

---

## 8. Responsive Behavior

### Breakpoints

Tailwind デフォルトに従う。

| Name | Width | 説明 |
| --- | --- | --- |
| Mobile | < 640px | 1 カラム縦積み、ハンバーガーメニュー |
| Tablet | ≥ 640px (`sm:`) / ≥ 768px (`md:`) | カード 2 列、フォーム 2 ペイン化(md) |
| Desktop | ≥ 1024px (`lg:`) | カード 3 列、フル幅レイアウト |

### タッチターゲット

- 最小サイズ: 44px × 44px(アイコンボタン・リスト行の操作要素を含む)
- ドラッグ&ドロップ操作(優先順位並び替え)には、タッチ環境向けの代替手段(上下移動ボタン)を必ず併設する

### フォントサイズの調整

- 本文 16px はモバイルでも維持する(入力欄は iOS ズーム回避のため 16px 必須)
- Display のみモバイルで 20px に縮小してよい

---

## 9. Agent Prompt Guide

### クイックリファレンス

```
Primary Color: #6366f1
Text Color: #374151
Text Secondary: #6b7280
Background: #f9fafb
Surface: #ffffff
Border: #e5e7eb
Font: Inter, "Noto Sans JP", sans-serif
Body Size: 16px
Line Height: 1.7(本文)/ 1.4(見出し)
Radius: 6px(ボタン・入力)/ 8px(カード)/ 12px(モーダル)
```

### プロンプト例

```
Adjusta のデザインシステム(frontend/DESIGN.md)に従って、候補日程リストを作成してください。
- プライマリカラー: #6366f1、フォント: Inter, "Noto Sans JP", sans-serif
- 行間: 本文 1.7、日時行は 1.5 まで許容
- カード: 背景 #ffffff、枠 #e5e7eb、radius 8px、padding 16px
- ステータスは色ドット + ラベル併記(confirmed=#22c55e、not_selected=#6b7280)
- loading / error / empty の3状態を用意する
```

---

## Migration Notes(現実装との差分)

本書は目指す姿の仕様であり、現実装(2026-07 時点)とは以下の差分がある。実装時に本書へ寄せる。

1. ~~**和文フォント未導入**~~ **解消済み(2026-07-09)**: Inter + Noto Sans JP を next/font で読み込み、CSS 変数経由で Tailwind の `fontFamily.sans` に登録済み
2. **`dark:` 残骸の整理**: ライト固定の決定に伴い、部分的に存在する `dark:` バリアント(4ファイル)と未配置の ThemeButton / ThemeProvider の扱いを整理する
3. **ボタン radius の変更**: 現実装は `rounded`(4px)。本書は 6px を規定
4. ~~**ページ背景**~~ **解消済み(2026-07-09)**: body に Background `#f9fafb`(`bg-background`)を適用し、Surface はカード側(`bg-card`)で分離
5. **未定義トークンの解消**: `ToggleSwitch` の `bg-primary-600` 等(Tailwind config に未定義)を本パレットの実クラスに置換する
6. **`text-md` の排除**: `text-base` へ統一(`TextField` 等)
7. ~~**行間・字間**~~ **解消済み(2026-07-09)**: body の base スタイルに 1.7 / 0.02em と禁則処理(`line-break: strict` / `overflow-wrap: break-word`)を導入。見出しは各コンポーネント側で行間 1.4 相当・字間 0 に詰める
8. **状態設計・a11y**: loading / error / empty の3状態、label 関連付け、タッチターゲット 44px は `../docs/ui-review.md` のバックログとして実装する
9. **shadcn/ui への段階的移行**(2026-07-09 決定・**セットアップ済み**): 本書のトークンは CSS 変数(HSL)として `globals.css` に定義し shadcn のテーマ変数に接続済み。`src/components/ui/` に button / card / badge / switch / skeleton / alert-dialog / radio-group を導入し、アカウント画面から利用開始。既存コンポーネントの置換は `../docs/ui-guidelines.md` 3.5 の移行規定に従い順次進める
10. **サードパーティ UI のテーマ未統一**: 現状 FullCalendar・react-datepicker はほぼ素のテーマのまま(react-datepicker の dist CSS を `(app)/layout.tsx` で素の状態で import、Splide は矢印の非表示 CSS のみ)。「4. Third-Party Components」の上書きが未実施
