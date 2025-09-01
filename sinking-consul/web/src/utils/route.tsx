import {history} from 'umi';
import route, {user, userPath, index, indexPath, admin, adminPath} from '../../config/routes'
import {Icon} from "@/components"
import React from "react"

/**
 * 递归获取完整路径
 * @param routes
 * @param name
 * @param params
 * @param currentPath
 */
function getPathByName(routes: any, name: any, params: any = {}, currentPath = ''): any {
    for (const route of routes) {
        let finalPath = currentPath + (currentPath && !currentPath.endsWith('/') ? '/' : '') + (route.path.startsWith('/') ? route.path.slice(1) : route.path);
        if (params && Object.keys(params).length > 0) {
            for (const key in params) {
                if (route.path.includes(`:${key}`)) {
                    finalPath = finalPath.replace(`:${key}`, params[key]);
                }
            }
        }
        if (route.name === name) {
            return finalPath.startsWith('/') ? finalPath : `/${finalPath}`;
        }
        if (route.routes) {
            const foundPath = getPathByName(route.routes, name, params, finalPath);
            if (foundPath) {
                return foundPath;
            }
        }
    }
    return null;
}

/**
 * 路由缓存
 */
const routeCache: any = {}

/**
 * 根据name获取路径
 * @param name 标识
 * @param params 参数
 */
export function getFullPathByName(name: any, params: any = {}): string {
    if (Object.keys(params).length <= 0 && routeCache?.[name] != undefined) {
        return routeCache[name];
    }
    routeCache[name] = getPathByName(route, name, params, '');
    return routeCache[name];
}

/**
 * 跳转页面
 * @param name 标识
 * @param params 参数
 */
export function historyPush(name: any, params = {}) {
    history?.push(getFullPathByName(name, params) || "/");
}

/**
 * 递归获取菜单
 * @param routes
 * @param parentPath
 * @param hideMenu
 */
function getMenuItems(routes: any, parentPath = '/', hideMenu = true) {
    return routes.map((route: any) => {
        if (hideMenu && route?.hideInMenu) {
            return;
        }
        const menuItem: any = {
            label: route?.title,
            path: route?.path,
            name: route?.name
        };
        if (route?.icon) {
            menuItem.icon = <Icon type={route?.icon}/>;
        }
        if (parentPath != '/' && !route.path.startsWith('/')) {
            menuItem.key = `${parentPath}/${route.path}`;
        } else {
            if (parentPath == '/' && route.path != '/') {
                menuItem.key = `/${route.path}`;
            } else {
                menuItem.key = `${route.path}`;
            }
        }
        if (route?.routes) {
            const children = getMenuItems(route?.routes, menuItem.key, hideMenu);
            if (children.length > 0) {
                menuItem.children = children;
            }
        }
        return menuItem;
    });
}


const parentCache: any = {}

/**
 * 根据name查找父级菜单
 * @param data 菜单
 * @param name 标识
 */
export function getParentList(data: any[], name: string): any[] {
    if (Object.keys(name).length <= 0 && parentCache?.[name] != undefined) {
        return parentCache[name];
    }
    let parents: any[] = [];

    function findParent(data: any, name: string): boolean {
        if (data?.children && data?.children?.length > 0) {
            if (data?.name === name) {
                parents.push(data);
                return true;
            }
        } else {
            for (const item of data) {
                if (item?.name === name) {
                    parents.push(item);
                    return true;
                }
                if (item?.children) {
                    if (findParent(item?.children, name)) {
                        parents.push(item);
                        return true;
                    }
                }
            }

        }
        return false;
    }

    findParent(data, name);
    parentCache[name] = parents.reverse()
    return parentCache[name];
}


/**
 * 获取user菜单
 */
export function getUserMenuItems(hideMenu = true) {
    return getMenuItems(user, '/' + userPath, hideMenu);
}

/**
 * 获取index菜单
 */
export function getIndexMenuItems(hideMenu = true) {
    return getMenuItems(index, '/' + indexPath, hideMenu);
}

/**
 * 获取admin菜单
 */
export function getAdminMenuItems(hideMenu = true) {
    return getMenuItems(admin, '/' + adminPath, hideMenu);
}

/**
 * 获取菜单
 */
export function getAllMenuItems(hideMenu = true) {
    return getMenuItems(route, '/', hideMenu)
}

/**
 * 当前访问系统path
 */
export function getCurrentPath(pathName: any): any {
    if (pathName == "/" || pathName == "") {
        pathName = "/index/index";
    }
    let mode = "";
    const regex = /\/([^/]+)\//; // 正则表达式匹配 / 之间的内容
    const matches = pathName.match(regex); // 匹配结果数组
    if (matches && matches.length >= 2) {
        mode = matches[1];
    }
    if (mode == userPath) {
        return '/' + userPath;
    }
    if (mode == indexPath) {
        return '/' + indexPath;
    }
    if (mode == adminPath) {
        return '/' + adminPath;
    }
    return '/' + indexPath;
}

/**
 * 当前访问系统路由
 */
export function getCurrentMenus(pathName: any, hideMenu = true): any {
    if (pathName == "/" || pathName == "") {
        pathName = "/index/index";
    }
    let mode = "";
    const regex = /\/([^/]+)\//; // 正则表达式匹配 / 之间的内容
    const matches = pathName.match(regex); // 匹配结果数组
    if (matches && matches.length >= 2) {
        mode = matches[1];
    }
    if (mode == userPath) {
        return getUserMenuItems(hideMenu);
    }
    if (mode == indexPath) {
        return getIndexMenuItems(hideMenu);
    }
    if (mode == adminPath) {
        return getAdminMenuItems(hideMenu);
    }
    return [];
}


/**
 * 获取第一个菜单(排除children)
 * @param menu
 */
export function getFirstMenuWithoutChildren(menu: any[]): any | null {
    for (const item of menu) {
        if (!item.hasOwnProperty('children')) {
            return item; // 如果对象没有 children 属性，则返回该对象
        } else {
            const result = getFirstMenuWithoutChildren(item.children);
            if (result !== null) {
                return result;
            }
        }
    }
    return null;
}
