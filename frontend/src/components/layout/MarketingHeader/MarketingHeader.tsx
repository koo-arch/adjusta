import React from 'react';
import Link from 'next/link';
import { Button } from '@/components/ui/button';

const MarketingHeader: React.FC = () => {
    return (
        <header className="sticky top-0 z-10 bg-white shadow-sm">
            <div className="mx-auto flex max-w-screen-2xl items-center justify-between px-4 py-4 md:px-8">
                <Link href="/">
                    <div className="cursor-pointer text-xl font-extrabold">
                        Adjusta
                    </div>
                </Link>
                <Button asChild>
                    <Link href="/login">ログイン</Link>
                </Button>
            </div>
        </header>
    );
};

export default MarketingHeader;
