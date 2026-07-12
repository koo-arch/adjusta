import { apiClient } from '@/lib/api/client';
import type { CandidateSyncSetting } from '@/features/auth/types';

export const updateCandidateSyncSetting = async (enabled: boolean) => {
    const response = await apiClient.put<CandidateSyncSetting, { enabled: boolean }>(
        '/api/calendar-settings/candidate-sync',
        { enabled },
    );
    return response.data;
};
