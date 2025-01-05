import React from 'react';
import Image from 'next/image';
import ScheduleManage from '../..//public/images/schedule_manage.jpg';
import Button from '@/components/Button';

const Home: React.FC = () => {
  return (
    <main className="min-h-screen flex flex-col lg:flex-row items-center justify-center">
      {/* Left Section */}
      <div className="flex flex-col items-center lg:items-start text-center lg:text-left px-6 lg:px-16 py-10 lg:py-0 w-full lg:w-1/2">
        <h1 className="text-4xl font-extrabold text-gray-800 mb-4 break-keep">
          日程調整を<wbr />もっとシンプルに
        </h1>
        <p className="text-lg text-gray-600 mb-6">
          あなたのイベントの日程調整を<wbr />効率的にサポートします。
        </p>
        <Button
          to="/login"
          shape="full"
          size="xl"
         >
          今すぐ始める
        </Button>
      </div>

      {/* Right Section */}
      <div className="w-full lg:w-1/2 h-64 lg:h-screen relative">
        <Image
          src={ScheduleManage}
          alt="スケジュール管理のイメージ"
          fill
          sizes="(max-width: 768px) 100vw, (max-width: 1200px) 50vw, 33vw"
          objectFit='cover'
          objectPosition='center'
          priority
        />
      </div>
    </main>
  );
};

export default Home;