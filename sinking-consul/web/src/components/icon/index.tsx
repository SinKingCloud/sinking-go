// noinspection TypeScriptValidateTypes,HtmlUnknownAttribute,JSUnusedGlobalSymbols

import React from 'react';
import {createFromIconfontCN} from '@ant-design/icons';
import * as AntdIcons from '@ant-design/icons';
import {createStyles} from 'antd-style';
import defaultSettings from '../../../config/defaultSettings';

const ICONFONT_URL = '//at.alicdn.com/t/c/font_5039718_a026bszsrgd.js';

const IconfontIcon = createFromIconfontCN({
    scriptUrl: ICONFONT_URL,
});

const useStyles: any = createStyles(() => ({
    iconWrapper: {
        display: 'inline-flex',
        alignItems: 'center',
        justifyContent: 'center',
    }
}));

interface IconProps extends React.HTMLAttributes<HTMLSpanElement> {
    type: string;
    ref?: any;
    className?: any;
    style?: any;
    onClick?: any;
}

const Icon = React.forwardRef<HTMLSpanElement, IconProps>(
    ({type, className, style, onClick, ...rest}, ref) => {
        const {styles, cx} = useStyles();
        const AntdIconComponent: any = (AntdIcons as any)[type];
        if (AntdIconComponent) {
            return (
                <AntdIconComponent ref={ref} className={cx(styles.iconWrapper, className)} style={style}
                                   onClick={onClick} {...rest}/>
            );
        }
        if (type.startsWith('icon-')) {
            return (
                <IconfontIcon type={type} ref={ref} className={cx(styles.iconWrapper, className)} style={style}
                              onClick={onClick} {...rest}/>
            );
        }
        return (
            <span ref={ref} className={cx(styles.iconWrapper, className)} style={style} onClick={onClick} {...rest}>
                {type}
            </span>
        );
    }
);

// 获取所有 Antd 图标名称
export const getAntdIconNames = (): string[] => {
    return Object.keys(AntdIcons).filter(key => {
        const component = (AntdIcons as any)[key];
        return typeof component === 'function' &&
            !key.startsWith('create') &&
            !key.startsWith('set') &&
            !key.startsWith('get') &&
            key !== 'default';
    });
};

// 获取所有 Iconfont 图标名称
let iconfontCache: string[] | null = null;

const normalizeUrl = (url: string): string => {
    if (url.startsWith('//')) {
        return `https:${url}`;
    }
    if (url.startsWith('http://') || url.startsWith('https://')) {
        return url;
    }
    const basePath = (defaultSettings?.basePath || '').replace(/\/$/, '');
    const srcPath = url.replace(/^\//, '');
    return basePath ? `${basePath}/${srcPath}` : `/${srcPath}`;
};

export const getIconfontIconNames = async (): Promise<string[]> => {
    if (iconfontCache) return iconfontCache;

    try {
        const response = await fetch(normalizeUrl(ICONFONT_URL));
        const scriptText = await response.text();
        const iconNames = [...scriptText.matchAll(/id="(icon-[^"]+)"/g)].map(m => m[1]);
        iconfontCache = iconNames;
        return iconNames;
    } catch (_) {
        return [];
    }
};

// 获取所有图标名称（Antd + Iconfont）
export const getAllIconNames = async (): Promise<string[]> => {
    const [antdIcons, iconfontIcons] = await Promise.all([
        Promise.resolve(getAntdIconNames()),
        getIconfontIconNames()
    ]);
    return [...antdIcons, ...iconfontIcons];
};
export {Icon};

/**
 * 图标数据
 */
const Exit = "icon-exit";
const Right = "icon-right";
const Log = "icon-log";
const Set = "icon-set";
const Bottom = "icon-bottom";
const Light = "icon-light";//亮色
const Dark = "icon-dark";
const Auto = "icon-auto";
const Home = "icon-home";
const Logo = "icon-logo";
export {
    Exit,
    Right,
    Log,
    Set,
    Bottom,
    Light,
    Dark,
    Auto,
    Home,
    Logo,
}