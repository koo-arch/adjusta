import { APIClientError } from '@/lib/api/client';

export const isUnauthorizedAPIError = (error: unknown): error is APIClientError =>
    error instanceof APIClientError && error.status === 401;

// 存在しない ID・他ユーザーのリソースはバックエンドが KindNotFound → 404 に変換する
export const isNotFoundAPIError = (error: unknown): error is APIClientError =>
    error instanceof APIClientError && error.status === 404;

// バックエンドのエラーボディは {"error": "<message>"} のみで機械可読コードを持たない
// (APIError.Kind は json:"-")。現状 409 を返すのは KindGoogleReauth
// (google_reauthorization_required)だけなので、暫定的に HTTP ステータスで判定する。
// メッセージ文字列でのマッチは脆いため行わない。
// フォローアップ: バックエンドのエラーボディに code フィールドが追加されたら判定を厳密化する
export const isGoogleReauthorizationRequiredError = (error: unknown): error is APIClientError =>
    error instanceof APIClientError && error.status === 409;
