'use client'
import { usePathname, useRouter, useSearchParams } from 'next/navigation';
import { useSearchEvents } from './useSearchEvents';

export const STATUS_TABS = ['all', 'active', 'confirmed', 'draft', 'cancelled'] as const;

export type StatusTab = (typeof STATUS_TABS)[number];

const parseStatusTab = (value: string | null): StatusTab =>
    STATUS_TABS.includes(value as StatusTab) ? (value as StatusTab) : 'all';

const parsePage = (value: string | null): number => {
    const page = Number(value);
    return Number.isInteger(page) && page >= 1 ? page : 1;
};

// 一覧の絞り込み状態(タブ・タイトル検索・ページ)を URL クエリと同期し、
// その条件で TanStack Query の検索を実行する
export const useEventListSearch = () => {
    const router = useRouter();
    const pathname = usePathname();
    const searchParams = useSearchParams();
    const statusTab = parseStatusTab(searchParams.get('status'));
    const page = parsePage(searchParams.get('page'));
    const title = (searchParams.get('title') ?? '').trim();

    const query = useSearchEvents({
        ...(statusTab === 'all' ? {} : { status: statusTab }),
        ...(title !== '' ? { title } : {}),
        page,
    });

    // タブ・検索テキストは直交フィルタとして相互に保持する
    const buildHref = ({ tab, targetPage, searchTitle }: { tab: StatusTab; targetPage: number; searchTitle: string }) => {
        const params = new URLSearchParams();
        if (tab !== 'all') {
            params.set('status', tab);
        }
        if (searchTitle !== '') {
            params.set('title', searchTitle);
        }
        if (targetPage > 1) {
            params.set('page', String(targetPage));
        }
        const queryString = params.toString();
        return queryString ? `${pathname}?${queryString}` : pathname;
    };

    const selectTab = (value: string) => {
        // タブ切替はページを 1 に戻す。履歴は絞り込み条件で汚さない
        router.replace(buildHref({ tab: parseStatusTab(value), targetPage: 1, searchTitle: title }), { scroll: false });
    };

    const search = (value: string) => {
        // 検索実行もページを 1 に戻す
        router.replace(buildHref({ tab: statusTab, targetPage: 1, searchTitle: value }), { scroll: false });
    };

    const goToPage = (nextPage: number) => {
        router.push(buildHref({ tab: statusTab, targetPage: nextPage, searchTitle: title }));
    };

    return {
        statusTab,
        title,
        page,
        selectTab,
        search,
        goToPage,
        ...query,
    };
};
