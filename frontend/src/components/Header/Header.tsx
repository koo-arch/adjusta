'use client'
import React from 'react';
import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { Button } from '@/components/ui/button';
import {
    Sheet,
    SheetClose,
    SheetContent,
    SheetDescription,
    SheetHeader,
    SheetTitle,
    SheetTrigger,
} from '@/components/ui/sheet';
import { CalendarDays, Menu } from 'lucide-react';
import { cn } from '@/lib/utils';

const navigation = [
    { name: 'ホーム', href: '/dashboard' },
    { name: 'イベント一覧', href: '/events' },
    { name: '新規作成', href: '/events/new' },
];

interface HeaderProps {
    userMenu: React.ReactNode;
}

const Header: React.FC<HeaderProps> = ({ userMenu }) => {
    const pathname = usePathname();
    const isActive = (href: string) => pathname === href;

    return (
        <header className="sticky top-0 z-40 border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/70">
            <div className="mx-auto flex h-16 w-full max-w-screen-2xl items-center gap-3 px-4 md:px-8">
                <Sheet>
                    <SheetTrigger asChild>
                        <Button
                            variant="ghost"
                            size="icon"
                            className="-ml-2 md:hidden"
                            aria-label="メインメニューを開く"
                        >
                            <Menu aria-hidden="true" />
                        </Button>
                    </SheetTrigger>
                    <SheetContent side="left" className="w-[86vw] max-w-sm">
                        <SheetHeader className="text-left">
                            <SheetTitle className="flex items-center gap-2">
                                <CalendarDays className="size-4 text-primary" aria-hidden="true" />
                                Adjusta
                            </SheetTitle>
                            <SheetDescription>主要ページへ移動できます。</SheetDescription>
                        </SheetHeader>
                        <nav className="mt-6 flex flex-col gap-1">
                            {navigation.map((item) => (
                                <SheetClose key={item.href} asChild>
                                    <Link
                                        href={item.href}
                                        className={cn(
                                            'border-l-2 border-transparent px-3 py-2 text-sm font-medium transition-colors hover:text-primary',
                                            isActive(item.href) && 'border-primary font-semibold text-primary',
                                        )}
                                    >
                                        {item.name}
                                    </Link>
                                </SheetClose>
                            ))}
                        </nav>
                    </SheetContent>
                </Sheet>

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
                            className={cn(
                                '-mb-px h-full rounded-none border-b-2 border-transparent px-3 hover:border-primary hover:bg-transparent hover:text-primary',
                                isActive(item.href) && 'border-primary font-semibold text-primary',
                            )}
                        >
                            <Link href={item.href}>{item.name}</Link>
                        </Button>
                    ))}
                </nav>

                <div className="ml-auto flex items-center gap-2">
                    {userMenu}
                </div>
            </div>
        </header>
    );
};

export default Header;
