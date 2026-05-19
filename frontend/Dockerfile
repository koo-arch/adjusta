FROM node:25.9-bookworm

# コンテナ内の作業ディレクトリを指定
WORKDIR /frontend

# package.jsonとpackage-lock.jsonをコピー
COPY package*.json ./

# ローカルのファイルをコンテナにコピー
COPY . .

# モジュールのダウンロード
RUN npm install

# ポートの公開
EXPOSE 3000

# コンテナ起動時のコマンド
CMD ["npm", "run", "dev"]
