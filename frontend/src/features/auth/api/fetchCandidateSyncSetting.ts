import { apiClient } from '@/lib/api/client';
import type { CandidateSyncSetting } from '@/features/auth/types';

export const fetchCandidateSyncSetting = async () => {
    const response = await apiClient.get<CandidateSyncSetting>('/api/calendar-settings/candidate-sync');
    return response.data;
};
