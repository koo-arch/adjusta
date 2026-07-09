import type { EventListParams, SearchParams } from './types';

export const buildDraftEventListQueryKey = (params: EventListParams = {}) =>
    ['draftEventList', params] as const;

export const buildEventDetailQueryKey = (eventID?: string) =>
    ['eventDetail', eventID] as const;

export const buildNeedsActionDraftsQueryKey = () => ['needsActionDrafts'] as const;

export const buildUpcomingEventsQueryKey = () => ['upcomingEvents'] as const;

export const buildDraftEventSearchQueryKey = (params: SearchParams) =>
    ['draftEventSearch', params] as const;
