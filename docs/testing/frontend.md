# フロントエンドテスト方針

## 目的

主要なユーザー導線と認証境界の回帰を、自動テストで早期に検出する。見た目の細かな差分ではなく、画面遷移、フォーム操作、validation、権限による表示制御を優先する。

## テストの役割

| 種別 | 対象 | 方針 |
| --- | --- | --- |
| E2E | 主要導線、画面遷移、認証境界 | Playwright を使用し、利用者の操作として確認する |
| Component | 複雑な入力、状態遷移、表示分岐 | 導入時にテストランナーを選定する。単純な表示だけのテストは増やさない |
| Storybook | UI部品の状態と目視確認 | デザイン確認に使用し、業務導線の保証はE2Eへ寄せる |
| Backend | validation、認可、業務ルール、API応答 | Go側のテストを正とし、フロントエンドで同じ条件を網羅し直さない |

snapshot は原則使用しない。文言やstyleだけの変更で、本質的でない差分が増えることを避ける。

## E2Eの実行

初回のみブラウザをインストールする。

```bash
cd frontend
npx playwright install chromium
```

テストを実行する。

```bash
npm run test:e2e
```

実装済みのテストケース一覧を表示する。

```bash
npm run test:e2e:list
```

直近のHTML reportをブラウザで表示する。

```bash
npm run test:e2e:report
```

`test:e2e` は通常の開発サーバーと競合しないよう、Next.js 開発サーバーを `http://localhost:3100` で自動起動し、生成物を `.next-e2e` に分離する。同じURLですでにサーバーが起動している場合は、そのサーバーを再利用する。

E2E起動時はローカルの `.env.local` に左右されないよう `NEXT_PUBLIC_API_BASE_URL` を空に固定し、ブラウザからのAPI呼び出しをNext.jsのsame-origin proxy経由でmock backendへ送る。

## 実行結果

`npm run test:e2e` は次の結果を生成する。

| 出力 | 用途 |
| --- | --- |
| `playwright-report/index.html` | ケースごとの成否、エラー、添付データを人が確認する |
| `test-results/e2e-results.xml` | CIやテスト結果集計で利用するJUnit XML |
| `test-results/*/trace.zip` | 失敗時の操作、通信、DOM状態を追跡する |
| `test-results/*/*.png` | 失敗時の画面状態を確認する |

これらは実行ごとに更新される生成物のためGitへコミットしない。CI導入時は、テストが失敗した場合も確認できるよう `playwright-report` と `test-results` をartifactとして保存する。

## E2Eのディレクトリ構成

specは画面のURLではなく、機能領域ごとに配置する。

```text
e2e/
├── public/      # LP・ログインなどの公開画面
├── auth/        # 認証状態と保護された導線
├── events/      # イベント一覧・作成・編集・確定
├── dashboard/   # ダッシュボード
├── account/     # アカウント・カレンダー設定
├── fixtures/    # 認証状態とテストデータの準備
└── helpers/     # 日時固定などの汎用的な補助処理
```

空のspecや将来利用するだけの共通処理は作らない。複数のspecで画面操作が重複するまではPage Objectを導入せず、テスト内に利用者の操作を明示する。

各テストのタイトルには `[AUTH-001]` のように、機能領域と3桁の連番からなるケースIDを付ける。specをテストケースの正とし、別のJSON台帳は作成しない。ケース一覧は `npm run test:e2e:list` で生成する。

## 自動化する主要導線

優先順位は次のとおりとする。

1. 公開ページと未認証時のredirect
2. イベントの作成とvalidation
3. イベントの編集
4. 候補日程の確定
5. 候補日程同期設定のON/OFF
6. APIエラーとGoogle再認可の表示

### 現在の到達点と追加候補

2026-07-13 時点で、上記の主要導線は mock backend を使用した E2E 53 ケースで一通り自動化済み。実装済みケースの正は各 spec とし、最新の一覧は `npm run test:e2e:list` で確認する。

次に補完する場合は、以下の順を目安とする。

1. アカウント: ログアウト成功・キャンセル、Google 再認可表示、プロフィール・連携状態、カレンダー取得失敗・空状態
2. イベント操作: 候補日程のドラッグ並び替え、手動入力による日程確定成功、候補日程のクリップボードコピー
3. 作成導線: 作成画面内での候補日程同期 ON/OFF と、設定失敗時の入力内容保持
4. ダッシュボード: イベント選択、popover / dialog、各セクションのエラー・再試行
5. 表示環境: モバイル幅、キーボード操作、必要性が確認できた場合の Chromium 以外のブラウザ
6. 統合テスト: mock backend ではなく実 backend・テスト DB を接続した主要導線。Google Calendar API は引き続き代替実装を使用し、実カレンダーを更新しない

これらは主要導線の完了を妨げるものではなく、実際の不具合傾向や UI 改修に合わせて個別タスクとして追加する。

## 認証と外部API

- Google OAuth の実画面や実アカウントを通常のE2Eでは操作しない。
- 認証後導線の自動化では、テスト環境に限定したセッション生成方法を別途設計する。
- Google Calendar API はテスト環境で代替実装またはfixtureを使用し、実カレンダーを更新しない。
- 本番で有効になる認証回避やテスト専用エンドポイントは追加しない。
- ブラウザの通信モックだけでは Server Component の `requireUser` を代替できないため、認証後導線を追加する前にbackendを含むテスト構成を決定する。
- 期限切れセッションの認証境界は、E2E専用mock backendから401を返して確認する。mockは `e2e/fixtures` に置き、本番コードから参照しない。

## テストデータ

- 各テストは単独で実行できる状態にする。
- テスト間でユーザーやイベントを共有しない。
- 日時に依存するテストでは時刻とタイムゾーンを固定する。
- テスト終了後に作成データを削除するか、テスト単位で破棄できるDBを使用する。

## CI

`.github/workflows/frontend-e2e.yml` で、frontend 関連ファイルを変更する PR と手動実行時に Chromium の E2E を実行する。Node.js 20、`npm ci`、`npx playwright install --with-deps chromium` を使用し、テストは CI 設定により 1 worker で動作する。

テスト失敗時は `frontend-e2e-results` artifact に HTML report、JUnit XML、trace、スクリーンショットを 7 日間保存する。複数ブラウザや画面幅の網羅は、実際の不具合傾向を確認してから追加する。
