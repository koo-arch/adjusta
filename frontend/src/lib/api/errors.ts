import { APIClientError } from '@/lib/api/client';

export const isUnauthorizedAPIError = (error: unknown): error is APIClientError =>
    error instanceof APIClientError && error.status === 401;

// 存在しない ID・他ユーザーのリソースはバックエンドが KindNotFound → 404 に変換する
export const isNotFoundAPIError = (error: unknown): error is APIClientError =>
    error instanceof APIClientError && error.status === 404;

export const isGoogleReauthorizationRequiredError = (error: unknown): error is APIClientError =>
    error instanceof APIClientError && error.code === 'google_reauthorization_required';
