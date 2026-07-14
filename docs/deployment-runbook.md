# Productionデプロイ手順書

最終更新: 2026-07-14

Adjustaのproduction環境を新規に構築し、その後も自分でデプロイ・確認・復旧できるようにするための操作手順書。

構成を選んだ理由や未決事項は[`deployment.md`](./deployment.md)、DB migrationの設計原則は[`db-design.md`](./db-design.md)を参照する。本書は実際の操作を扱う。

### 本書の使い方

本書は次の2つの用途を区別する。

- **現在のproductionを更新・復旧する場合**: 表に記載した`adjusta`の既存値を使用し、作成済みリソースの作成手順は確認だけ行う。
- **別のproduction環境をゼロから作る場合**: project ID、project number、Cloud Run URL、Vercel URLを新環境で取得した値に置き換える。`adjusta`の固定値をそのままコピーしない。

各節の「完了条件」を満たしてから次へ進む。管理画面の名称は変更されることがあるため、画面操作で見つからない場合は併記した公式ドキュメントを確認する。

## 1. 完成後の構成

```text
Browser
  └─ https://adjusta.vercel.app
       ├─ Next.js pages
       └─ /api/* route handlers
            └─ Cloud Run: adjusta-api
                 ├─ Neon PostgreSQL
                 ├─ Google OAuth / Calendar API
                 └─ Secret Manager

GitHub Actions
  ├─ GitHub OIDC → Workload Identity Federation
  ├─ Docker image → Artifact Registry
  ├─ Atlas migration → Neon
  └─ Backend revision → Cloud Run
```

productionで使用している値は次のとおり。

Cloud Runに独立した「Cloud Runプロジェクト」があるわけではない。Cloud Run service、Artifact Registry、Secret Manager、service account、Workload Identity Federation、Google OAuth clientは、すべて表のGCP project `adjusta`内のリソースとして管理する。別環境を作る場合も、特別な理由がなければ同じ1つのGCP projectへまとめる。

| 項目 | 値 |
| --- | --- |
| GCP project ID | `adjusta` |
| GCP project number | `1001878278191` |
| GCP region | `asia-northeast1` |
| Artifact Registry repository | `adjusta` |
| Cloud Run service | `adjusta-api` |
| Cloud Run URL | `https://adjusta-api-1001878278191.asia-northeast1.run.app` |
| Runtime service account | `adjusta-runtime@adjusta.iam.gserviceaccount.com` |
| Deploy service account | `adjusta-deploy@adjusta.iam.gserviceaccount.com` |
| Vercel URL | `https://adjusta.vercel.app` |
| Application DB schema | `adjusta` |
| Atlas revision table | `public.atlas_schema_revisions` |

## 2. 事前準備とGCP project作成

必要なもの:

- GCP billing account
- Vercel、Neon、GitHubのアカウント
- Google OAuth同意画面を管理できる権限
- GitHub repositoryのActions・Environmentsを変更できる権限
- ローカルの`gcloud` CLI。`gh` CLIは任意

### 2.1 操作者の権限

個人の新規GCP projectでは、作成者のGoogle accountで初回構築を行うのが簡単。組織配下では管理者に次を確認する。

- project作成元のorganizationまたはfolderに`resourcemanager.projects.create`権限。`roles/resourcemanager.projectCreator`に含まれる
- 対象billing accountにprojectを関連付ける権限。通常は`roles/billing.user`
- 作成後のprojectで、API、IAM、service account、Workload Identity Pool、Artifact Registry、Secret Manager、Cloud Runを作成・変更できる権限
- Cloud Runを公開するための`run.services.setIamPolicy`、またはInvoker IAM checkを無効化する権限

個人projectでは初回構築中だけProject Ownerで進められるが、通常のデプロイは後述する`adjusta-deploy`へ限定する。組織環境ではOwnerを恒常付与せず、管理者に必要な管理操作を依頼する。

project作成には[Create projects](https://docs.cloud.google.com/resource-manager/docs/creating-managing-projects)、billing accountの権限には[Cloud Billing access control](https://cloud.google.com/billing/docs/how-to/billing-access)を参照する。

### 2.2 gcloudへログインする

```bash
gcloud auth login
```

現在のproductionを操作する場合は、既存projectを選択する。

```bash
export PROJECT_ID="adjusta"
gcloud config set project "$PROJECT_ID"
gcloud projects describe "$PROJECT_ID"
```

### 2.3 新規projectを作成する場合

この節は別環境をゼロから作る場合だけ実行する。project IDは全世界で一意であり、作成後は変更できない。`adjusta`が利用できない場合は、組織名や環境名を含む一意なIDを選ぶ。

```bash
export PROJECT_ID="YOUR_UNIQUE_PROJECT_ID"
gcloud projects create "$PROJECT_ID" --name="Adjusta"
```

利用できるbilling accountを確認し、対象IDを設定する。billing account IDは秘密値ではないが、リポジトリへ固定する必要はない。

```bash
gcloud billing accounts list
export BILLING_ACCOUNT_ID="YOUR_BILLING_ACCOUNT_ID"
gcloud billing projects link "$PROJECT_ID" \
  --billing-account="$BILLING_ACCOUNT_ID"
```

projectをgcloudの既定値に設定する。作成直後にprojectが見つからない場合は、作成処理の反映を数分待って再試行する。

```bash
gcloud config set project "$PROJECT_ID"
```

project作成コマンド自体はbillingを自動で有効化しない。詳細は[Create projects](https://docs.cloud.google.com/resource-manager/docs/creating-managing-projects)と[`gcloud billing projects`](https://docs.cloud.google.com/sdk/gcloud/reference/billing/projects)を参照する。

費用の想定外増加を検知できるよう、GCP Consoleの`Billing > Budgets & alerts`で少額の予算通知も設定する。予算は利用を自動停止する機能ではない。

### 2.4 共通変数を設定する

以降のコマンド用に、秘密ではない値を設定する。別環境ではservice名やrepository名も必要に応じて変更する。

```bash
export PROJECT_NUMBER="$(gcloud projects describe "$PROJECT_ID" \
  --format='value(projectNumber)')"
export REGION="asia-northeast1"
export REPOSITORY="adjusta"
export CLOUD_RUN_SERVICE="adjusta-api"
export RUNTIME_SA="adjusta-runtime"
export DEPLOY_SA="adjusta-deploy"
export RUNTIME_SA_EMAIL="${RUNTIME_SA}@${PROJECT_ID}.iam.gserviceaccount.com"
export DEPLOY_SA_EMAIL="${DEPLOY_SA}@${PROJECT_ID}.iam.gserviceaccount.com"
export GITHUB_REPOSITORY="koo-arch/adjusta"
export WIF_POOL="github"
export WIF_PROVIDER="adjusta"
```

完了条件:

```bash
test -n "$PROJECT_NUMBER"
test "$(gcloud config get-value project)" = "$PROJECT_ID"
gcloud projects describe "$PROJECT_ID" \
  --format='table(projectId,projectNumber,state)'
gcloud billing projects describe "$PROJECT_ID"
```

- projectのstateが`ACTIVE`
- project numberが空でない
- billingが有効で、意図したbilling accountへ紐付いている
- `gcloud config get-value project`が対象project IDと一致する

### 秘密値の扱い

- DB URL、OAuth client secret、session secretをGit、Markdown、issue、PR、チャットへ貼らない。
- `--data-file=-`で標準入力から登録するか、各サービスの管理画面を使う。
- `echo SECRET`や、秘密値を引数へ直接書くコマンドはshell historyやログへ残り得るため使わない。
- 秘密値を誤って表示・共有した場合は、削除だけで済ませず発行元でローテーションする。
- OAuthの認可codeも一時的な資格情報なので共有しない。

## 3. Neonを作成する

1. [Neon Console](https://console.neon.tech/)で新しいprojectを作成する。
2. 名前を`adjusta`とし、Cloud Runに近いregionを選択する。
3. `Connect`からproduction branchのconnection stringを取得する。
4. 初回はpoolerを介さないdirect connection stringを使用し、`sslmode=require`が含まれることを確認する。
5. URLへ`search_path`は追加しない。entのqueryとmigration側で`adjusta` schemaを明示する。
6. connection stringはpassword manager等へ一時保管し、後のSecret Manager登録後は平文ファイルを残さない。

Neonはproject作成時に`public` schemaを作る。Adjustaのmigrationはアプリケーションテーブルを`adjusta`へ作り、Atlasの履歴だけを`public.atlas_schema_revisions`へ保存する。初回migration前に手動でテーブルを作成しない。

Neonのprojectと接続方法は[Projects](https://neon.com/docs/manage/projects)、[Connection pooling](https://neon.com/docs/connect/connection-pooling)を参照する。poolerへ切り替える場合は、backend接続とAtlas migrationの両方を別途検証する。

完了条件:

- production branch、database、roleが作成されている
- direct connection stringを取得できている
- SQL Editorへ接続できる
- migration前なのでAdjustaのアプリケーションテーブルはまだ存在しない

## 4. GCPの基盤を作成する

### 4.1 APIを有効化する

```bash
gcloud services enable \
  artifactregistry.googleapis.com \
  iam.googleapis.com \
  iamcredentials.googleapis.com \
  run.googleapis.com \
  secretmanager.googleapis.com \
  sts.googleapis.com
```

Google Calendar APIも、GCP Consoleの`APIs & Services > Library`から有効化する。

完了条件:

```bash
gcloud services list --enabled \
  --filter='config.name:(artifactregistry.googleapis.com OR iamcredentials.googleapis.com OR run.googleapis.com OR secretmanager.googleapis.com OR sts.googleapis.com)'
```

表示に列挙したAPIがすべて含まれ、GCP Console上でGoogle Calendar APIも有効になっている。

### 4.2 Artifact Registryを作成する

```bash
gcloud artifacts repositories create "$REPOSITORY" \
  --repository-format=docker \
  --location="$REGION" \
  --description="Adjusta backend images"
```

既に存在する場合は作成不要。確認する。

```bash
gcloud artifacts repositories describe "$REPOSITORY" \
  --location="$REGION"
```

詳細は[Create standard repositories](https://cloud.google.com/artifact-registry/docs/repositories/create-repos)と[Docker quickstart](https://docs.cloud.google.com/artifact-registry/docs/docker/store-docker-container-images)を参照する。

### 4.3 service accountを作成する

Cloud Runでbackendを実行するidentityと、GitHub Actionsがデプロイに使うidentityを分ける。

```bash
gcloud iam service-accounts create "$RUNTIME_SA" \
  --display-name="Adjusta runtime"

gcloud iam service-accounts create "$DEPLOY_SA" \
  --display-name="Adjusta deploy"
```

deploy identityへ必要な権限を付与する。

```bash
gcloud projects add-iam-policy-binding "$PROJECT_ID" \
  --member="serviceAccount:${DEPLOY_SA_EMAIL}" \
  --role="roles/run.developer"

gcloud artifacts repositories add-iam-policy-binding "$REPOSITORY" \
  --location="$REGION" \
  --member="serviceAccount:${DEPLOY_SA_EMAIL}" \
  --role="roles/artifactregistry.writer"

gcloud iam service-accounts add-iam-policy-binding "$RUNTIME_SA_EMAIL" \
  --member="serviceAccount:${DEPLOY_SA_EMAIL}" \
  --role="roles/iam.serviceAccountUser"
```

Cloud Runのidentityと権限は[Service identity](https://docs.cloud.google.com/run/docs/securing/service-identity)と[IAM roles](https://docs.cloud.google.com/run/docs/reference/iam/roles)を参照する。

### 4.4 Cloud Run serviceを仮作成する

VercelとOAuthに設定するURLを先に確定するため、Googleのhello imageでserviceを作成する。

```bash
gcloud run deploy "$CLOUD_RUN_SERVICE" \
  --image="us-docker.pkg.dev/cloudrun/container/hello" \
  --region="$REGION" \
  --service-account="$RUNTIME_SA_EMAIL" \
  --no-invoker-iam-check \
  --min-instances=0
```

URLを確認する。

```bash
export BACKEND_URL="$(gcloud run services describe "$CLOUD_RUN_SERVICE" \
  --region="$REGION" \
  --format='value(status.url)')"
printf '%s\n' "$BACKEND_URL"
```

AdjustaではVercelのserver-side proxyから呼び出すため、Cloud Run serviceは未認証呼び出しを許可する。アプリケーションの認証はsession cookieで行う。`--no-invoker-iam-check`は、domain restricted sharingで`allUsers`をIAMへ追加できない環境でも使用できる現在の推奨方法。Cloud Runの公開設定は[Allowing public access](https://docs.cloud.google.com/run/docs/authenticating/public)、deploy操作は[Deploying container images](https://docs.cloud.google.com/run/docs/deploying)を参照する。

完了条件:

- `$BACKEND_URL`が空でなく、`https://`から始まる
- ブラウザまたは`curl "$BACKEND_URL"`でhello responseを取得できる
- Cloud Run serviceのregionとruntime service accountが想定どおり

別環境では、以後の`https://adjusta-api-1001878278191.asia-northeast1.run.app`をすべて`$BACKEND_URL`の実値へ置き換える。

## 5. Vercel frontendを作成する

1. Vercel Dashboardで`Add New > Project`を開く。
2. GitHubの`koo-arch/adjusta` repositoryをimportする。
3. Framework PresetはNext.js、Root Directoryは`frontend`にする。
4. Production Branchは`main`にする。
5. project名を`adjusta`にし、production domainを確認する。名前が利用済みならVercelが割り当てた実際のdomainを使用する。
6. `Settings > Environment Variables`で、Productionに次を登録する。

| 名前 | 値 |
| --- | --- |
| `INTERNAL_BACKEND_URL` | 手順4.4で取得した`$BACKEND_URL`。現在のproductionは`https://adjusta-api-1001878278191.asia-northeast1.run.app` |

`NEXT_PUBLIC_API_BASE_URL`はproductionへ登録しない。ブラウザはVercelのsame-origin `/api/*`を呼び、Next.js Route Handlerだけが`INTERNAL_BACKEND_URL`を使ってCloud Runへ転送する。

環境変数を変更しただけでは既存deploymentに反映されないため、変更後はRedeployする。Previewは環境変数なしでもbuildできるが、認証済み画面を動かすにはPreview専用backend・OAuth callbackを用意する必要がある。production backendを安易にPreview OAuthへ流用しない。

設定方法は[Vercel Git deployments](https://vercel.com/docs/git)、[Project settings](https://vercel.com/docs/project-configuration/project-settings)、[Environment variables](https://vercel.com/docs/environment-variables)を参照する。

Vercelが割り当てたproduction URLを控える。現在のproduction以外では、以後の`https://adjusta.vercel.app`をすべてこの値へ置き換える。

```bash
export FRONTEND_URL="https://YOUR_VERCEL_PRODUCTION_DOMAIN"
```

完了条件:

- production deploymentのbuildが成功している
- production domainをHTTPSで開ける
- Production環境に`INTERNAL_BACKEND_URL`が登録されている
- Production環境に`NEXT_PUBLIC_API_BASE_URL`が登録されていない

## 6. Google OAuthを設定する

GCP Consoleの`Google Auth Platform`または`APIs & Services`から設定する。

1. OAuth consent screenへアプリ名、support email、developer contactを登録する。
2. Data Accessで、実装が要求する次のscopeを確認する。
   - `openid`
   - `https://www.googleapis.com/auth/userinfo.email`
   - `https://www.googleapis.com/auth/userinfo.profile`
   - `https://www.googleapis.com/auth/calendar`
3. Web applicationのOAuth clientを作成する。
4. Authorized redirect URIへ、次の値を完全一致で登録する。

```text
https://adjusta.vercel.app/api/auth/google/callback
```

5. 必要ならAuthorized JavaScript originへ`https://adjusta.vercel.app`を登録する。
6. OAuth client IDとclient secretを取得する。
7. testing状態なら利用するGoogle accountをtest userへ追加する。本番公開前に、Calendar scopeを含む同意画面の公開状態とGoogleのverification要否を確認する。

redirect URIはscheme、host、path、末尾slashまで一致する必要がある。詳細は[Using OAuth 2.0 for web server applications](https://developers.google.com/identity/protocols/oauth2/web-server)を参照する。

別環境ではredirect URIを`${FRONTEND_URL}/api/auth/google/callback`、ログイン後URLを`${FRONTEND_URL}`として扱う。

完了条件:

- Google Calendar APIが、OAuth clientと同じGCP projectで有効
- OAuth consent screenに必要なscopeとtest userまたは公開設定がある
- Web application clientが作成され、client IDとclient secretを取得済み
- Authorized redirect URIがVercelのproduction callbackと完全一致

## 7. Secret Managerへ秘密値を登録する

次の3 secretを作る。

```bash
gcloud secrets create adjusta-database-url \
  --replication-policy=automatic

gcloud secrets create adjusta-session-secret \
  --replication-policy=automatic

gcloud secrets create adjusta-google-client-secret \
  --replication-policy=automatic
```

各値は標準入力から追加する。コマンド実行後に値を貼り付け、`Ctrl-D`で終了する。

```bash
gcloud secrets versions add adjusta-database-url --data-file=-
gcloud secrets versions add adjusta-session-secret --data-file=-
gcloud secrets versions add adjusta-google-client-secret --data-file=-
```

登録内容:

| Secret | 内容 |
| --- | --- |
| `adjusta-database-url` | Neonのdirect connection string |
| `adjusta-session-secret` | password manager等で生成した十分に長いランダム値 |
| `adjusta-google-client-secret` | Google OAuth client secret |

runtime identityへ、3 secretそれぞれのaccessor権限を付ける。

```bash
for SECRET in \
  adjusta-database-url \
  adjusta-session-secret \
  adjusta-google-client-secret
do
  gcloud secrets add-iam-policy-binding "$SECRET" \
    --member="serviceAccount:${RUNTIME_SA_EMAIL}" \
    --role="roles/secretmanager.secretAccessor"
done
```

workflowはmigration時にDB URLを読み、Cloud Run deploy時にsecretを関連付ける。deploy identityにもsecret単位でaccessor権限を付ける。

```bash
for SECRET in \
  adjusta-database-url \
  adjusta-session-secret \
  adjusta-google-client-secret
do
  gcloud secrets add-iam-policy-binding "$SECRET" \
    --member="serviceAccount:${DEPLOY_SA_EMAIL}" \
    --role="roles/secretmanager.secretAccessor"
done
```

Secret ManagerとCloud Runの関連付けは[Create and access secrets](https://docs.cloud.google.com/secret-manager/docs/creating-and-accessing-secrets)と[Configure secrets for Cloud Run](https://docs.cloud.google.com/run/docs/configuring/services/secrets)を参照する。

完了条件:

```bash
for SECRET in \
  adjusta-database-url \
  adjusta-session-secret \
  adjusta-google-client-secret
do
  gcloud secrets describe "$SECRET" --format='value(name)'
  gcloud secrets versions list "$SECRET" \
    --filter='state=ENABLED' \
    --format='value(name)'
done
```

- 3つのsecretに有効なversionが1つ以上ある
- runtime service accountとdeploy service accountにsecret単位のaccessor権限がある
- secretの値を確認コマンドやログへ出していない

## 8. GitHub ActionsのWorkload Identity Federationを作成する

service account key JSONは発行せず、GitHubのOIDC tokenをGCPの短期credentialへ交換する。

### 8.1 poolとproviderを作成する

```bash
gcloud iam workload-identity-pools create "$WIF_POOL" \
  --location=global \
  --display-name="GitHub Actions"

gcloud iam workload-identity-pools providers create-oidc "$WIF_PROVIDER" \
  --location=global \
  --workload-identity-pool="$WIF_POOL" \
  --display-name="GitHub koo-arch/adjusta" \
  --issuer-uri="https://token.actions.githubusercontent.com" \
  --attribute-mapping="google.subject=assertion.sub,attribute.repository=assertion.repository,attribute.ref=assertion.ref" \
  --attribute-condition="assertion.repository=='koo-arch/adjusta' && assertion.ref=='refs/heads/main'"
```

providerの完全なresource名を取得する。後でGitHub secretへ登録する。

```bash
gcloud iam workload-identity-pools providers describe "$WIF_PROVIDER" \
  --location=global \
  --workload-identity-pool="$WIF_POOL" \
  --format='value(name)'
```

poolの完全なresource名を取得し、対象repositoryからdeploy identityをimpersonateできるようにする。

```bash
WIF_POOL_NAME="$(gcloud iam workload-identity-pools describe "$WIF_POOL" \
  --location=global \
  --format='value(name)')"

gcloud iam service-accounts add-iam-policy-binding "$DEPLOY_SA_EMAIL" \
  --role="roles/iam.workloadIdentityUser" \
  --member="principalSet://iam.googleapis.com/${WIF_POOL_NAME}/attribute.repository/${GITHUB_REPOSITORY}"
```

GitHub側のworkflowには`permissions: id-token: write`が必要で、現在の`.github/workflows/backend-deploy.yml`には設定済み。詳細は[Workload Identity Federation with deployment pipelines](https://docs.cloud.google.com/iam/docs/workload-identity-federation-with-deployment-pipelines)と[google-github-actions/auth](https://github.com/google-github-actions/auth)を参照する。

IAMとWorkload Identity Federationの変更は反映まで数分かかることがある。作成直後の認証失敗では、値を変更し続ける前にprovider状態とbindingを確認して再試行する。

完了条件:

- providerのstateが`ACTIVE`
- providerの完全なresource名を取得できる
- provider conditionが`koo-arch/adjusta`の`main`だけを許可している
- deploy service accountに対象repositoryの`roles/iam.workloadIdentityUser` bindingがある

## 9. GitHub production environmentを設定する

GitHub repositoryの`Settings > Environments > New environment`で`production`を作る。

推奨設定:

- Deployment branches and tags: `Selected branches and tags`で`main`のみ
- Required reviewers: 初回確認中は自分を設定
- Prevent self-review: 個人運用で承認できなくなる場合は無効

Environment variablesへ登録する。

| 名前 | 値 |
| --- | --- |
| `GCP_PROJECT_ID` | `adjusta` |
| `GCP_REGION` | `asia-northeast1` |
| `GCP_ARTIFACT_REPOSITORY` | `adjusta` |
| `CLOUD_RUN_SERVICE` | `adjusta-api` |
| `CLOUD_RUN_SERVICE_ACCOUNT` | `adjusta-runtime@adjusta.iam.gserviceaccount.com` |
| `DATABASE_URL_SECRET` | `adjusta-database-url` |
| `SESSION_SECRET_SECRET` | `adjusta-session-secret` |
| `GOOGLE_CLIENT_SECRET_SECRET` | `adjusta-google-client-secret` |
| `GOOGLE_CLIENT_ID` | Google OAuth client ID |
| `GOOGLE_REDIRECT_URI` | `https://adjusta.vercel.app/api/auth/google/callback` |
| `REDIRECT_URL_AFTER_LOGIN` | `https://adjusta.vercel.app` |
| `CORS_ALLOW_ORIGINS` | `https://adjusta.vercel.app` |
| `COOKIE_DOMAIN` | 空文字。未登録でもよい |

Environment secretsへ登録する。

| 名前 | 値 |
| --- | --- |
| `GCP_WORKLOAD_IDENTITY_PROVIDER` | 手順8で取得したproviderの完全なresource名 |
| `GCP_DEPLOY_SERVICE_ACCOUNT` | `adjusta-deploy@adjusta.iam.gserviceaccount.com` |

OAuth client IDは秘密値ではないが、client secretは必ずSecret Managerだけへ置く。GitHub environmentsの保護規則は[Deployments and environments](https://docs.github.com/en/actions/reference/workflows-and-actions/deployments-and-environments)を参照する。

別環境では、`GOOGLE_REDIRECT_URI`、`REDIRECT_URL_AFTER_LOGIN`、`CORS_ALLOW_ORIGINS`へその環境の`$FRONTEND_URL`を使用する。Cloud Run URLはGitHub environmentへ直接登録せず、Vercelの`INTERNAL_BACKEND_URL`へ登録する。

完了条件:

- `production` environmentが存在する
- environment variablesとsecretsが表の名前で登録されている
- deployment branch policyが`main`だけを許可する
- required reviewerを設定した場合、自分で承認できない設定になっていない
- GitHubへservice account key JSONを登録していない

## 10. 初回backend deployを実行する

`.github/workflows/backend-deploy.yml`は手動実行で、次の順序で処理する。

1. GitHub OIDCでGCPへ認証
2. `backend/Dockerfile`からimageをbuild
3. Artifact Registryへpush
4. AtlasでNeonのpending migrationを適用
5. 新しいCloud Run revisionをdeploy

GitHubの`Actions > Backend Deploy > Run workflow`から`main`を選んで実行し、production environmentの承認待ちになったら内容を確認して承認する。

`gh` CLIを使う場合:

```bash
gh workflow run backend-deploy.yml --ref main
gh run list --workflow backend-deploy.yml --limit 5
gh run watch RUN_ID --exit-status
```

workflowのmigrationは次の指定を使用する。

```bash
atlas migrate apply \
  --dir "file://backend/migrations" \
  --url "${DATABASE_URL}" \
  --revisions-schema "public"
```

`backend/migrations`のSQLが変更履歴の正本で、適用状態は`public.atlas_schema_revisions`へ記録される。`--allow-dirty`や初回DBへのbaselineは使用しない。[Atlas migrate apply](https://atlasgo.io/versioned/apply)も参照する。

完了条件:

- `Backend Deploy`の全stepが成功
- GitHub Actionsのcommit SHAとCloud Run revisionのimage SHAが一致
- Neonに`public.atlas_schema_revisions`と`adjusta` schemaが存在
- Cloud Runのlatest ready revisionがhello imageではなく`adjusta-backend` imageを使用
- Cloud Run URLへ未認証で到達でき、`/api/users/me`が`401`を返す

## 11. Vercelを再デプロイする

`INTERNAL_BACKEND_URL`をproductionへ追加した後、Vercel DashboardのDeploymentsからmainの最新deploymentをRedeployする。以後はmainへのpush/mergeでfrontendが自動デプロイされる。

Vercelのproduction URLで画面を開き、Cloud Runのhello画面ではなくAdjusta APIへ転送されることを確認する。

## 12. 初回動作確認

### 12.1 インフラ確認

Cloud RunのURLとrevisionを確認する。

```bash
gcloud run services describe "$CLOUD_RUN_SERVICE" \
  --region="$REGION" \
  --format='yaml(status.url,status.latestReadyRevisionName)'
```

未認証のユーザーAPIは`401`を返せば、Cloud Runへ到達できている。

```bash
curl --silent --output /dev/null --write-out '%{http_code}\n' \
  "$BACKEND_URL/api/users/me"
```

OAuth開始endpointはredirectを返す。

```bash
curl --silent --output /dev/null --write-out '%{http_code}\n' \
  "$FRONTEND_URL/api/auth/google/login"
```

期待値はそれぞれ`401`と`307`。

Neon SQL Editorでmigrationを確認する。

```sql
SELECT version, applied
FROM public.atlas_schema_revisions
ORDER BY version;

SELECT table_schema, table_name
FROM information_schema.tables
WHERE table_schema IN ('adjusta', 'public')
ORDER BY table_schema, table_name;
```

`public`にはAtlasの履歴以外のアプリケーションテーブルを置かない。

### 12.2 OAuthとcookie確認

ブラウザで次を確認する。

1. `${FRONTEND_URL}/login`を開く。現在のproductionは`https://adjusta.vercel.app/login`。
2. Googleログインを完了する。
3. `${FRONTEND_URL}`へ戻り、ログイン済み画面が表示される。
4. DevToolsのApplication/Storageからsession cookieを確認する。

期待するcookie属性:

- `HttpOnly`
- `Secure`
- `SameSite=Lax`
- `Path=/`
- Domain属性なしのhost-only cookie

### 12.3 主要導線確認

- イベントを下書き作成できる
- 候補日程を編集・保存できる
- イベント詳細を表示できる
- 確定日時を登録・変更できる
- イベントを削除できる
- アカウント画面で候補日程同期をON/OFFできる
- ONの場合はGoogle Calendarへ候補日程が一度だけ同期される
- ログアウト後に認証必須画面へ戻れない

## 13. 通常のデプロイ

### Frontend

1. frontend変更のPRでFrontend E2Eを成功させる。
2. mainへmergeする。
3. Vercelの自動deploymentが成功したことを確認する。
4. productionで変更箇所をsmoke testする。

### Backend

1. backend変更のPRでBackend CIを成功させる。
2. DB変更がある場合はmigration SQLとexpand/contract互換性をreviewする。
3. mainへmergeする。
4. `Backend Deploy` workflowをmainから手動実行する。
5. production approval時にcommit SHAとmigrationを再確認する。
6. workflow成功後、Cloud Run logsと主要APIを確認する。

アプリケーションdeployより前にmigrationが実行されるが、それだけで安全になるわけではない。旧Cloud Run revisionがmigration後も動作できる変更にする。

## 14. ログと状態確認

GitHub Actionsの失敗stepを見る。

```bash
gh run view RUN_ID --log-failed
```

Cloud Run logsを見る。

```bash
gcloud run services logs read "$CLOUD_RUN_SERVICE" \
  --region="$REGION" \
  --limit=100
```

revision一覧を見る。

```bash
gcloud run revisions list \
  --service="$CLOUD_RUN_SERVICE" \
  --region="$REGION"
```

Vercelは対象deploymentの`Build Logs`と`Runtime Logs`、Neonは`Monitoring`とSQL Editorを確認する。ログへcookie、Authorization header、DB URL、OAuth codeを出さない。

## 15. ロールバック

### Cloud Run

問題のないrevision名を確認し、trafficを戻す。

```bash
gcloud run services update-traffic "$CLOUD_RUN_SERVICE" \
  --region="$REGION" \
  --to-revisions=REVISION_NAME=100
```

DB migrationはCloud Runのrollbackでは戻らない。破壊的migrationを同じreleaseに含めず、旧revisionが新schemaで動作できる状態を維持する。データを戻す必要がある場合は、影響を確認した専用migrationとして扱う。

### Vercel

Vercel DashboardのDeploymentsから正常なproduction deploymentを選び、RollbackまたはPromoteを行う。frontendだけを戻してもbackend APIとの互換性が必要になる。

## 16. secretのローテーション

新しい値を発行し、既存secretへ新versionとして登録する。

```bash
gcloud secrets versions add SECRET_NAME --data-file=-
```

その後`Backend Deploy`を再実行し、`latest`を参照する新Cloud Run revisionを作る。

- session secretを変えると既存sessionは無効になる。
- OAuth client secretはGoogle側で再発行してからSecret Managerを更新する。
- Neon passwordを変えたら、新しいconnection stringを登録して接続確認後に旧credentialを失効させる。
- 古いsecret versionは新revisionの確認後に無効化する。

## 17. よくあるエラー

### Atlas: connected database is not clean

初期migration前のDBに既存tableやschemaがある。新規production DBではbaselineや`--allow-dirty`を使わず、意図せず作られたobjectの由来を確認する。

破棄可能なローカルDBならvolumeを作り直す。保持が必要なDBは`public`から`adjusta`への明示的な移行計画を作る。

### Atlas: atlas_schema_revisions関連のrelationエラー

Neon productionではworkflowの`atlas migrate apply`に次が必要。

```text
--revisions-schema public
```

過去の失敗で空の`atlas_schema_revisions` schemaだけが残っている場合は、中身が空であることを確認してからNeon SQL Editorで削除する。

```sql
DROP SCHEMA atlas_schema_revisions;
```

既存の正しいrevision tableを削除してはいけない。`public.atlas_schema_revisions`が正しい保存先である。

### Vercel buildでINTERNAL_BACKEND_URLがない

現在はserver routeをdynamicにしているため、Preview build時に値がなくてもbuildできる。production runtimeには`INTERNAL_BACKEND_URL`が必須。Vercelへ追加後にRedeployする。

### OAuth後に/loginへ戻る

次を順に確認する。

1. `INTERNAL_BACKEND_URL`がCloud Run URLになっているか
2. `GOOGLE_REDIRECT_URI`とGoogle Consoleのredirect URIが完全一致するか
3. `REDIRECT_URL_AFTER_LOGIN`が実際の`$FRONTEND_URL`か。現在のproductionは`https://adjusta.vercel.app`
4. callback responseのcookieがVercel側で中継されているか
5. cookieが`Secure`、`SameSite=Lax`、host-onlyか
6. Cloud Run logsにDB・OAuthエラーがないか

### redirect_uri_mismatch

Google ConsoleとGitHub environmentの`GOOGLE_REDIRECT_URI`を比較する。`http/https`、末尾slash、pathの違いも不一致になる。変更後はbackendを再デプロイする。

### Cloud Runが403を返す

public invoker設定を確認する。

```bash
gcloud run services get-iam-policy "$CLOUD_RUN_SERVICE" \
  --region="$REGION"
```

まずInvoker IAM checkの状態を確認し、無効化する。これはdomain restricted sharing下でも利用できる推奨方法。

```bash
gcloud run services update "$CLOUD_RUN_SERVICE" \
  --region="$REGION" \
  --no-invoker-iam-check
```

組織ポリシーでこの操作も制限されている場合は、組織管理者による許可が必要。`allUsers`への`roles/run.invoker`付与は代替手段だが、domain restricted sharingでは拒否されることがある。

### GitHub ActionsのGCP認証に失敗する

- workflowに`id-token: write`があるか
- `GCP_WORKLOAD_IDENTITY_PROVIDER`が完全なresource名か
- providerのrepository・branch conditionが実行元と一致するか
- `production` environmentのbranch policyがmainを許可しているか
- deploy service accountに`roles/iam.workloadIdentityUser` bindingがあるか
- GitHub Actions上で承認待ちになっていないか

### Secret Managerのpermission denied

runtime service accountとdeploy service accountを混同していないか確認する。runtimeはアプリ起動時、deployはmigrationとrevision作成時にsecretを扱う。権限はproject全体ではなく、対象secretごとに付与する。

## 18. 構成変更時の更新箇所

domain、service名、secret名などを変更した場合は、最低限次を同期する。

- `.github/workflows/backend-deploy.yml`
- GitHub `production` environment variables/secrets
- Vercel `INTERNAL_BACKEND_URL`
- Google OAuth authorized redirect URI
- Cloud Runの環境変数・secret mapping・IAM
- Secret Manager IAM
- `docs/deployment.md`
- 本書のproduction値

設定とドキュメントが食い違った場合は、実環境を無条件に正とせず、Git履歴と直近の成功workflowを確認して差分の理由を明確にする。
