@ AGENTS.md
## ドキュメント参照ルール

- 機能の追加・変更時は、まず `docs/requirements.md` で対象機能の要件を確認すること
- DBスキーマ・entのschema定義に触れる変更では、必ず `docs/db-design.md` を先に読むこと
- レイヤー間の依存方向やパッケージ配置に迷ったら `docs/rearchitecture-memo.md` を参照すること。
  ここにDDD再設計の経緯と判断理由が書かれている
- ドキュメントと実装が食い違っている場合は、勝手にどちらかに合わせず報告すること