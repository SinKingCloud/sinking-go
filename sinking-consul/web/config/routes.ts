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
        name: "index",
        icon: 'icon-home',
        hideInMenu: false,
        hideBreadCrumb: true,
        component: "@/pages/index",
    },
    {
        path: "cluster",
        title: "集群列表",
        name: "cluster",
        icon: 'icon-cluster',
        hideInMenu: false,
        component: "@/pages/cluster",
    },
    {
        path: "node",
        title: "服务节点",
        name: "node",
        icon: 'icon-node',
        hideInMenu: false,
        component: "@/pages/node",
    },
    {
        path: "config",
        title: "配置管理",
        name: "config",
        icon: 'icon-config',
        hideInMenu: false,
        component: "@/pages/config",
    },
    {
        path: "log",
        title: "操作日志",
        name: "log",
        icon: 'icon-log-blue',
        hideInMenu: false,
        component: "@/pages/log",
    },
    {
        path: "system",
        title: "系统管理",
        name: "system",
        icon: 'icon-setting',
        hideInMenu: false,
        component: "@/pages/system",
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
        path: '403',
        component: '@/pages/403.tsx',
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