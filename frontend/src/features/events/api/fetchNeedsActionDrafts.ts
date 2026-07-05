import { apiClient } from '@/lib/api/client';
import type { NeedsActionDraft } from '@/features/events/types';

export const fetchNeedsActionDrafts = async () => {
    const response = await apiClient.get<NeedsActionDraft[]>('/api/event/draft/needs-action');
    return response.data;
};
