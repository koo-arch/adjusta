'use client'
import React, { useSyncExternalStore } from 'react';
import Link from 'next/link';
import DraftRegisterButton from '@/features/events/draft/components/DraftRegisterButton';
import { usePathname } from 'next/navigation';
import { Button } from '@/components/ui/button';
import {
    Sheet,
    SheetClose,
    SheetContent,
    SheetHeader,
    SheetTitle,
    SheetTrigger,
} from '@/components/ui/sheet';
import { Menu } from 'lucide-react';

const navigation = [
    { name: 'ホーム', href: '/dashboard' },
    { name: 'イベント一覧', href: '/events' },
]

const classNames = (...classes: string[]) => {
    return classes.filter(Boolean).join(' ')
}

const subscribePageScroll = (onStoreChange: () => void) => {
    if (typeof window === 'undefined') {
        return () => {};
    }

    let frameId: number | null = null;
    const onScroll = () => {
        if (frameId != null) {
            return;
        }

        frameId = window.requestAnimationFrame(() => {
            frameId = null;
            onStoreChange();
        });
    };

    window.addEventListener('scroll', onScroll);
    return () => {
        window.removeEventListener('scroll', onScroll);
        if (frameId != null) {
            window.cancelAnimationFrame(frameId);
        }
    };
};

const getHasPageShadow = () => {
    if (typeof window === 'undefined') {
        return false;
    }

    return window.scrollY > 0;
};

interface HeaderProps {
    userMenu: React.ReactNode;
}

const Header: React.FC<HeaderProps> = ({ userMenu }) => {
    const pathname = usePathname();
    const hasShadow = useSyncExternalStore(
        subscribePageScroll,
        getHasPageShadow,
        () => false,
    );

    const isActived = (href: string) => pathname === href;

    return (
        <nav className={`sticky top-0 z-10 bg-white transition-shadow ${hasShadow ? 'shadow-md' : ''}`}>
            <div className="mx-auto max-w-screen-2xl px-4 md:px-8">
                <div className='relative flex items-center justify-between py-4'>
                    <div className="absolute inset-y-0 left-0 flex items-center md:hidden">
                        <Sheet>
                            <SheetTrigger asChild>
                                <Button variant="ghost" size="icon" aria-label="メインメニューを開く">
                                    <Menu aria-hidden="true" />
                                </Button>
                            </SheetTrigger>
                            <SheetContent side="left" className="w-72">
                                <SheetHeader className="text-left">
                                    <SheetTitle>Adjusta</SheetTitle>
                                </SheetHeader>
                                <div className="mt-8 space-y-2">
                                    {navigation.map((item) => (
                                        <SheetClose key={item.name} asChild>
                                            <Link
                                                href={item.href}
                                                className={classNames(
                                                    isActived(item.href) ? 'text-primary' : 'hover:bg-accent',
                                                    'block rounded-md px-3 py-2 text-base font-semibold transition-colors',
                                                )}
                                            >
                                                {item.name}
                                            </Link>
                                        </SheetClose>
                                    ))}
                                </div>
                            </SheetContent>
                        </Sheet>
                    </div>
                            <div className="flex flex-1 items-center justify-center md:items-stretch md:justify-start">
                                <div className="flex flex-shrink-0 items-center">
                                    <Link href="/dashboard">
                                        <div className="cursor-pointer text-xl font-extrabold">
                                            Adjusta
                                        </div>
                                    </Link>
                                </div>
                                <div className="hidden md:ml-6 md:block">
                                    <div className="flex gap-12">
                                        {navigation.map((item) => (
                                            <Link key={item.name} href={item.href} className={classNames(
                                                isActived(item.href) ? 'text-indigo-500' : 'hover:text-indigo-500 active:text-indigo-700',
                                                "text-lg font-semibold transition duration-100"
                                            )}
                                            >
                                                    {item.name}
                                            </Link>
                                        ))}
                                    </div>
                                </div>
                            </div>
                            <div className="flex items-center space-x-4">
                                <DraftRegisterButton />
                                {userMenu}
                            </div>
                </div>
            </div>
        </nav>
    )
}

export default Header;
