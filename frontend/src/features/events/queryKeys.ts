import type { SearchParams } from './types';

export const buildEventDetailQueryKey = (eventID?: string) =>
    ['eventDetail', eventID] as const;

export const buildNeedsActionDraftsQueryKey = () => ['needsActionDrafts'] as const;

export const buildUpcomingEventsQueryKey = () => ['upcomingEvents'] as const;

// params 省略時は prefix キーを返す(mutation 後の一括 invalidate 用)
export const buildDraftEventSearchQueryKey = (params?: SearchParams) =>
    params === undefined ? (['draftEventSearch'] as const) : (['draftEventSearch', params] as const);
