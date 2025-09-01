/**
 * 用户系统路由
 */
export const userPath = "user";
export const user = [
    {
        path: "index",
        title: "数据概览",
        name: "user.index",
        icon: 'icon-home',
        hideInMenu: false,
        component: "@/pages/user/index",
    },
    {
        path: "pay",
        title: "财务管理",
        name: "user.pay",
        icon: 'icon-wallet',
        hideInMenu: false,
        routes: [
            {
                path: "recharge",
                component: "@/pages/user/pay/recharge",
                title: "账户充值",
                name: "user.recharge",
                icon: 'icon-dot',
                hideInMenu: false,
            },
            {
                path: "order",
                component: "@/pages/user/pay/order",
                title: "订单记录",
                name: "user.order",
                icon: 'icon-dot',
                hideInMenu: false,
            },
            {
                path: "log",
                component: "@/pages/user/pay/log",
                title: "余额明细",
                name: "user.log",
                icon: 'icon-dot',
                hideInMenu: false,
            },
        ]
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
                component: "@/pages/user/person/setting",
                title: "资料管理",
                name: "user.setting",
                icon: 'icon-dot',
                hideInMenu: false,
            },
            {
                path: "log",
                component: "@/pages/user/person/log",
                title: "操作日志",
                name: "user.log2",
                icon: 'icon-dot',
                hideInMenu: false,
            },
        ]
    },
];
export const adminPath = "admin";
export const admin = [
    {
        path: "index",
        component: "@/pages/admin/index",
        title: "网站概览",
        name: "admin.index",
        icon: 'icon-home',
        hideInMenu: false,
    },
    {
        path: "user",
        component: "@/pages/admin/index",
        title: "用户管理",
        name: "admin.user",
        icon: 'icon-users',
        hideInMenu: false,
    },
    {
        path: "notice",
        component: "@/pages/admin/index",
        title: "公告管理",
        name: "admin.notice",
        icon: 'icon-message',
        hideInMenu: false,
    },
    {
        path: "order",
        component: "@/pages/admin/index",
        title: "订单管理",
        name: "admin.order",
        icon: 'icon-order',
        hideInMenu: false,
    },
    {
        path: "setting",
        component: "@/pages/admin/index",
        title: "系统管理",
        name: "admin.setting",
        icon: 'icon-setting',
        hideInMenu: false,
    },
]
/**
 * 首页系统路由
 */
export const indexPath = "index";
export const index = [
    {
        path: "index",
        component: "@/pages/index/index",
        title: "系统首页",
        name: "index.index",
        icon: 'icon-doc',
        hideInMenu: false,
    },
];
/**
 * 系统路由
 */
export default [
    /**
     * 首页路由
     */
    {
        path: "/",
        name: "redirect." + indexPath,
        title: index[0]?.title,
        routes: [
            {
                path: "/",
                component: index[0]?.component,
                title: index[0]?.title,
                name: index[0]?.name,
                hideInMenu: true,
            },
            {
                path: '/login',
                component: '@/pages/user/login',
                name: "login",
                title: "帐号登录",
                hideInMenu: false,
            },
        ],
    },
    {
        path: indexPath,
        name: indexPath,
        title: "系统首页",
        routes: index,
    },
    /**
     * 用户路由
     */
    {
        path: "/" + userPath,
        name: "redirect." + userPath,
        redirect: '/' + userPath + '/' + (user[0]?.path || 'index'),
        hideInMenu: true,
    },
    {
        path: userPath,
        name: userPath,
        title: "用户系统",
        routes: user,
    },
    /**
     * admin路由
     */
    {
        path: "/" + adminPath,
        name: "redirect." + adminPath,
        redirect: '/' + adminPath + '/' + (admin[0]?.path || 'index'),
        hideInMenu: true,
    },
    {
        path: adminPath,
        name: adminPath,
        title: "后台网站",
        icon: 'icon-set',
        routes: admin,
    },
    /**
     * 支付结果
     */
    {
        path: "pay",
        component: "@/pages/pay.tsx",
        name: "pay",
        title: "支付结果",
        layout: false,
        hideInMenu: true,
    },
    /**
     * 500
     */
    {
        path: '/500',
        component: '@/pages/500.tsx',
        name: "error",
        title: "服务器错误",
        layout: false,
        hideInMenu: true,
    },
    /**
     * 403
     */
    {
        path: '/401',
        component: '@/pages/401.tsx',
        name: "notAllowed",
        title: "无权限",
        layout: false,
        hideInMenu: true,
    },
    /**
     * 404
     */
    {
        path: '/*',
        component: '@/pages/404.tsx',
        name: "notFound",
        title: "页面不存在",
        layout: false,
        hideInMenu: true,
    }
];