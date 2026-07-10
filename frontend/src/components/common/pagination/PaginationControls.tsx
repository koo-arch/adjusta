'use client'
import type { MouseEvent } from 'react';
import { cn } from '@/lib/utils';
import { ChevronLeft, ChevronRight } from 'lucide-react';
import {
    Pagination,
    PaginationContent,
    PaginationEllipsis,
    PaginationItem,
    PaginationLink,
} from '@/components/ui/pagination';

interface PaginationControlsProps {
    page: number;
    total: number;
    limit: number;
    onPageChange: (page: number) => void;
}

// 先頭・末尾と現在ページの前後 1 ページを表示し、間は省略記号にする
const buildPageNumbers = (current: number, totalPages: number): Array<number | string> => {
    if (totalPages <= 5) {
        return Array.from({ length: totalPages }, (_, index) => index + 1);
    }

    const pages = new Set<number>();
    pages.add(1);
    pages.add(totalPages);

    for (let i = current - 1; i <= current + 1; i++) {
        if (i > 1 && i < totalPages) {
            pages.add(i);
        }
    }

    return Array.from(pages)
        .sort((a, b) => a - b)
        .reduce<Array<number | string>>((acc, value, index, array) => {
            acc.push(value);
            const next = array[index + 1];
            if (typeof next === 'number' && next - value > 1) {
                acc.push('…');
            }
            return acc;
        }, []);
};

export const PaginationControls = ({
    page,
    total,
    limit,
    onPageChange,
}: PaginationControlsProps) => {
    const totalPages = Math.max(1, Math.ceil(total / limit));

    const pages = buildPageNumbers(page, totalPages);
    const rangeStart = total === 0 ? 0 : (page - 1) * limit + 1;
    const rangeEnd = Math.min(page * limit, total);

    const handleNavigation = (nextPage: number) => {
        if (nextPage < 1 || nextPage > totalPages || nextPage === page) {
            return;
        }
        onPageChange(nextPage);
    };

    const buildClickHandler = (nextPage: number, disabled = false) => (event: MouseEvent) => {
        event.preventDefault();
        if (disabled) return;
        handleNavigation(nextPage);
    };

    return (
        <div className="mt-6 grid gap-4 border-t pt-4 text-sm text-muted-foreground">
            <div className="text-center sm:text-left">
                {rangeStart.toLocaleString()}–{rangeEnd.toLocaleString()} / {total.toLocaleString()} 件
            </div>
            <Pagination className="w-full justify-center">
                <PaginationContent className="flex-wrap">
                    <PaginationItem>
                        <PaginationLink
                            href="#"
                            size="default"
                            aria-label="前のページへ"
                            onClick={buildClickHandler(page - 1, page === 1)}
                            className={cn('gap-1 px-3', page === 1 && 'pointer-events-none opacity-50')}
                        >
                            <ChevronLeft className="size-4" />
                            <span className="hidden sm:inline">前へ</span>
                        </PaginationLink>
                    </PaginationItem>
                    {pages.map((value, index) => (
                        <PaginationItem key={typeof value === 'number' ? value : `ellipsis-${index}`}>
                            {typeof value === 'number' ? (
                                <PaginationLink
                                    href="#"
                                    isActive={value === page}
                                    onClick={buildClickHandler(value)}
                                >
                                    {value}
                                </PaginationLink>
                            ) : (
                                <PaginationEllipsis />
                            )}
                        </PaginationItem>
                    ))}
                    <PaginationItem>
                        <PaginationLink
                            href="#"
                            size="default"
                            aria-label="次のページへ"
                            onClick={buildClickHandler(page + 1, page === totalPages)}
                            className={cn('gap-1 px-3', page === totalPages && 'pointer-events-none opacity-50')}
                        >
                            <span className="hidden sm:inline">次へ</span>
                            <ChevronRight className="size-4" />
                        </PaginationLink>
                    </PaginationItem>
                </PaginationContent>
            </Pagination>
        </div>
    );
};
