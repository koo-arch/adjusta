import Link from 'next/link';
import { CalendarDays, Menu, Plus } from 'lucide-react';
import { Button } from '@/components/ui/button';
import UserMenuSkeleton from '@/features/auth/components/UserMenuSkeleton';

const navigation = [
    { name: 'ホーム', href: '/dashboard' },
    { name: 'イベント一覧', href: '/events' },
];

const HeaderSkeleton = () => (
    <header className="sticky top-0 z-40 border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/70">
        <div className="mx-auto flex h-16 w-full max-w-screen-2xl items-center gap-3 px-4 md:px-8">
            <Button
                variant="ghost"
                size="icon"
                className="-ml-2 md:hidden"
                disabled
                aria-label="メインメニューを読み込み中"
            >
                <Menu aria-hidden="true" />
            </Button>

            <Link
                href="/dashboard"
                className="inline-flex shrink-0 items-center gap-2 font-semibold tracking-tight"
            >
                <CalendarDays className="size-5 text-primary" aria-hidden="true" />
                <span>Adjusta</span>
            </Link>

            <nav className="hidden h-full flex-1 items-stretch gap-1 md:flex">
                {navigation.map((item) => (
                    <Button
                        key={item.href}
                        asChild
                        size="sm"
                        variant="ghost"
                        className="-mb-px h-full rounded-none border-b-2 border-transparent px-3 hover:border-primary hover:bg-transparent hover:text-primary"
                    >
                        <Link href={item.href}>{item.name}</Link>
                    </Button>
                ))}
            </nav>

            <div className="ml-auto flex items-center gap-2">
                <Button
                    asChild
                    variant="ghost"
                    size="icon"
                    aria-label="新規作成"
                    title="新規作成"
                >
                    <Link href="/events/new">
                        <Plus aria-hidden="true" />
                    </Link>
                </Button>
                <UserMenuSkeleton />
            </div>
        </div>
    </header>
);

export default HeaderSkeleton;
