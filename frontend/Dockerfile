FROM node:18-alpine

# コンテナ内の作業ディレクトリを指定
WORKDIR /frontend

# Gitなどのツールのインストール
RUN apk update && apk add --no-cache git

# package.jsonとpackage-lock.jsonをコピー
COPY package*.json ./

# ローカルのファイルをコンテナにコピー
COPY . .

# モジュールのダウンロード
RUN npm install

# ビルド
RUN npm run build

# ポートの公開
EXPOSE 3000

# コンテナ起動時のコマンド
CMD ["npm", "run", "dev"]
