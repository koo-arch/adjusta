import React from 'react';
import DashboardCalendar from '@/features/dashboard/components/DashboardCalendar';
import DashboardSidePanel from '@/features/dashboard/components/DashboardSidePanel';

const DashboardPageContainer = () => {
    return (
        <main className="mx-auto max-w-screen-2xl space-y-6 px-4 py-8 md:px-8">
            <h1 className="text-2xl font-bold leading-snug tracking-normal text-gray-900">ホーム</h1>
            {/* カレンダーを主役に、右にコンテキストパネル(作成フォームと同じ構図)。
                768〜1024px はパネルを引くとカレンダーが潰れるため lg から 2 カラム化する */}
            <div className="grid grid-cols-1 gap-8 lg:grid-cols-[minmax(0,1fr)_24rem] lg:gap-6">
                <section className="lg:order-2">
                    <DashboardSidePanel />
                </section>
                <section className="lg:order-1">
                    <DashboardCalendar />
                </section>
            </div>
        </main>
    );
};

export default DashboardPageContainer;
