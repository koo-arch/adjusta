import type { AuthUser } from '@/features/auth/types';
import { isGoogleReauthorizationRequiredError } from '@/lib/api/errors';

export type GoogleConnectionStatus = 'connected' | 'reauthorization_required';

export type GoogleConnectionState =
    | { kind: 'loading' }
    | { kind: 'connected'; account: AuthUser }
    | { kind: 'reauthorization_required' }
    | { kind: 'load_failed' };

export const GOOGLE_CONNECTION_STATUS: Record<
    GoogleConnectionStatus,
    { label: string; color: 'green' | 'yellow' }
> = {
    connected: {
        label: '正常',
        color: 'green',
    },
    reauthorization_required: {
        label: '要再認可',
        color: 'yellow',
    },
};

interface ResolveGoogleConnectionStateParams {
    account?: AuthUser;
    isLoading: boolean;
    error: unknown;
}

export const resolveGoogleConnectionState = ({
    account,
    isLoading,
    error,
}: ResolveGoogleConnectionStateParams): GoogleConnectionState => {
    if (isLoading) {
        return { kind: 'loading' };
    }
    if (isGoogleReauthorizationRequiredError(error)) {
        return { kind: 'reauthorization_required' };
    }
    if (error || !account) {
        return { kind: 'load_failed' };
    }
    return { kind: 'connected', account };
};
