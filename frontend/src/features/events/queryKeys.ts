import type { SearchParams } from './types';

export const buildDraftEventListQueryKey = () => ['draftEventList'] as const;

export const buildEventDetailQueryKey = (eventID?: string) =>
    ['eventDetail', eventID] as const;

export const buildNeedsActionDraftsQueryKey = () => ['needsActionDrafts'] as const;

export const buildUpcomingEventsQueryKey = () => ['upcomingEvents'] as const;

export const buildDraftEventSearchQueryKey = (params: SearchParams) =>
    ['draftEventSearch', params] as const;
