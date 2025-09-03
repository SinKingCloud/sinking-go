/**
 * 系统路由
 */
export default [
    /**
     * 业务
     */
    {
        path: "/",
        title: "数据概览",
        name: "user.index",
        icon: 'icon-home',
        hideInMenu: false,
        component: "@/pages/index",
    },
    {
        path: "test",
        title: "测试菜单",
        name: "user.index2",
        icon: 'icon-home',
        hideInMenu: false,
        component: "@/pages/index",
    },
    {
        path: "person",
        title: "账户管理",
        name: "user.person",
        icon: 'icon-user',
        hideInMenu: false,
        routes: [
            {
                path: "setting",
                component: "@/pages/index",
                title: "资料管理",
                name: "user.setting",
                icon: 'icon-dot',
                hideInMenu: false,
            },
            {
                path: "log",
                component: "@/pages/index",
                title: "操作日志",
                name: "user.log2",
                icon: 'icon-dot',
                hideInMenu: false,
            },
        ]
    },
    /**
     * 登录
     */
    {
        path: 'login',
        component: '@/pages/login',
        name: "login",
        title: "帐号登录",
        auth: false,
        hideInMenu: true,
    },
    /**
     * 500
     */
    {
        path: '500',
        component: '@/pages/500.tsx',
        name: "error",
        title: "服务器错误",
        auth: false,
        hideInMenu: true,
    },
    /**
     * 403
     */
    {
        path: '401',
        component: '@/pages/401.tsx',
        name: "notAllowed",
        title: "无权限",
        auth: false,
        hideInMenu: true,
    },
    /**
     * 404
     */
    {
        path: '*',
        component: '@/pages/404.tsx',
        name: "notFound",
        title: "页面不存在",
        auth: false,
        hideInMenu: true,
    }
];