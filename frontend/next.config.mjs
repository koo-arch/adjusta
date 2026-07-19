/** @type {import('next').NextConfig} */
const nextConfig = {
    cacheComponents: true,
    // E2Eでは通常の開発サーバーとlockや生成物を共有しない。
    distDir: process.env.NEXT_DIST_DIR ?? '.next',
    images: {
        remotePatterns: [
            {
                protocol: 'https',
                hostname: 'tailwindui.com',
                port: '',
            },
            {
                protocol: 'https',
                hostname: 'images.unsplash.com',
                port: '',
            },
            {
                protocol: 'https',
                hostname: 'lh3.googleusercontent.com',
                port: '',
            }
        ]
    },
    // イベント系 URL の再設計(screen-design 9.1/9.2)に伴う旧 URL からの恒久リダイレクト。
    // 先に定義したものが優先されるため、/register は :id パターンより前に置く
    async redirects() {
        return [
            {
                source: '/schedule/draft/register',
                destination: '/events/new',
                permanent: true,
            },
            {
                source: '/schedule/draft/:id/edit',
                destination: '/events/:id/edit',
                permanent: true,
            },
            {
                source: '/schedule/draft/:id',
                destination: '/events/:id',
                permanent: true,
            },
            {
                source: '/schedule/draft',
                destination: '/events',
                permanent: true,
            },
            {
                // /schedule(カレンダー単体ページ)は廃止し、ダッシュボードのカレンダー表示に統合
                source: '/schedule',
                destination: '/dashboard',
                permanent: true,
            },
        ];
    },
};

export default nextConfig;
