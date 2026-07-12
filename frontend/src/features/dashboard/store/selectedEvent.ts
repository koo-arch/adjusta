import { atom } from 'jotai';
import type { CalendarEvent } from '@/features/calendar/types';

// カレンダーでクリックされたイベント。右パネル(タブ ⇄ 詳細)と共有する
export const selectedDashboardEventAtom = atom<CalendarEvent | null>(null);

export type DashboardPanelTab = 'needs-action' | 'upcoming' | 'detail';

// 右パネルのアクティブタブ。カレンダー側からも切り替えるため atom に置く
export const dashboardPanelTabAtom = atom<DashboardPanelTab>('needs-action');
