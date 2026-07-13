# バックエンドテスト方針

## 目的

認証・認可、業務ルール、transaction、Google Calendar 同期の回帰を自動テストで早期に検出する。実装の内部構造ではなく、入力に対する状態遷移、永続化、外部サービス呼び出し、API応答を優先して確認する。

## テストの役割

| 対象 | 方針 |
| --- | --- |
| domain | 外部依存なしで業務ルールと状態遷移を確認する |
| usecase | repository・transaction・外部APIをfakeに置き換え、処理順序と失敗時の状態を確認する |
| handler | validation、認証境界、HTTP status、機械可読なエラーcodeを確認する |
| repository | PostgreSQL互換のテストDBを使用し、query、制約、Soft Deleteを確認する |
| infrastructure | Google APIなどの通信先をローカルHTTPサーバーに置き換え、requestとresponse変換を確認する |

生成されたentコードや単純なDTO変換を網羅するためだけのテストは追加しない。repositoryの動作をSQLiteだけで保証せず、PostgreSQL固有のindexや制約はPostgreSQLで確認する。

## 実行方法

バックエンド全体を実行する。

```bash
cd backend
go test ./...
```

対象パッケージだけを実行する。

```bash
cd backend
go test ./internal/usecase/events
```

競合し得る共有状態や並行処理を追加した場合はrace detectorも実行する。

```bash
cd backend
go test -race ./...
```

各テストは単独・任意順で実行できるようにし、外部サービス、開発用DB、実ユーザーのデータに依存させない。

## Google Calendar同期

通常の自動テストでは実Google Calendar APIを呼ばない。usecaseではfake gateway、infrastructureではローカルHTTPサーバーを使用する。

優先して保証するケースは次のとおり。

1. 同期成功時にGoogle Event ID、`synced`、最終同期日時を保存し、直前のエラーを消去する
2. 同期失敗時にAdjusta側のデータを保持し、`sync_failed`とエラー内容を保存する
3. 次回の同期操作で失敗状態から復旧できる
4. 保存済みGoogle Event IDがある場合は新規作成せず更新する
5. Google側で予定が削除されていた場合は再作成し、新しいEvent IDを保存する
6. 同期済みの候補日程を詳細画面で再取得してもGoogle Calendar APIを呼ばない
7. 同じ同期操作を繰り返しても予定を重複作成しない
8. 候補日程同期がOFFの場合はGoogle Calendar APIを呼ばない
9. token失効時は再認可が必要なエラーとして扱う

Google側への作成成功後、DBへのEvent ID保存に失敗する境界は、Google APIとDBを単一transactionにできないため別途対策が必要になる。安定した冪等性キーまたは照合方法を設計するまでは、既知の未保証ケースとして扱う。

実Googleアカウントを使う確認は自動テストに含めず、デプロイ前後の手動確認として行う。専用のテストアカウントとカレンダーを使い、作成した予定を確認後に削除する。

## DBテスト

domainとusecaseのテストはDBなしで実行する。repository、migration、PostgreSQL固有制約のテストを追加する場合は、テストごとに破棄可能なPostgreSQLを用意し、以下を守る。

- productionや開発用DATABASE_URLを使用しない
- テスト間でrecordを共有しない
- transaction rollbackまたはschema/databaseの破棄で後処理する
- partial unique index、foreign key、Soft Delete条件は実PostgreSQLで確認する

## CI

バックエンド関連ファイルを変更するPRと手動実行時に、GitHub Actionsで`go test ./...`を実行する。通常のunit testとadapter testはDBやGoogle認証情報を必要としない構成を維持する。

CIではAtlasがmigration履歴を一時PostgreSQLへ適用し、ent schemaとの差分がないことを確認する。repositoryテストを追加する場合も破棄可能なPostgreSQLへ分離し、実Google Calendar APIはCIから呼ばない。
