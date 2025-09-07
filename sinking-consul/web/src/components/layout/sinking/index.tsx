import Footer from "../footer";
import Header from "../header";
import Sider from "../sider";
import {createStyles, useResponsive, useTheme} from "antd-style";
import React, {forwardRef, useEffect, useImperativeHandle, useState} from "react";
import {Outlet, useLocation, useRouteData, useSelectedRoutes} from "umi";
import {MenuFoldOutlined, MenuUnfoldOutlined} from "@ant-design/icons";
import Loading from "@/loading"
import {App, Breadcrumb, Button, ConfigProvider, Drawer, Layout, Watermark} from "antd";
import zhCN from 'antd/locale/zh_CN';
import {getAllMenuItems, getFirstMenuWithoutChildren, getParentList, historyPush} from "@/utils/route";

const useLayoutStyles = createStyles(({isDarkMode, token, css, responsive}): any => {
    return {
        container: {
            display: "flex",
            flexDirection: "row",
        },
        sider: {
            zIndex: 2,
            boxShadow: "2px 0 8px 0 rgba(29,35,41,.05)",
            background: "none !important",
            height: '100vh',
            overflowY: "auto",
            position: "sticky",
            left: 0,
            top: 0,
            bottom: 0,
            transition: "all 0.3s,background 0s,height 0s",
            ".ant-layout-header": {
                backgroundColor: token?.colorBgContainer
            },
        },
        menuBtn: {
            width: "55px !important",
            height: "55px",
            lineHeight: "55px",
            fontSize: "15px !important",
            cursor: "pointer",
            float: "left",
            ":hover": {
                background: "none !important",
            }
        },
        drawMenu: {
            padding: "0px !important",
        },
        body: {
            transition: "margin-left 0.3s !important",
            backgroundColor: isDarkMode ? "black" : "transparent",
        },
        sticky: {
            position: "sticky",
            top: 0,
            left: 0,
        },
        header: {
            padding: 0,
            height: "55px",
            lineHeight: "55px",
            width: "100%",
            zIndex: 3,
            boxShadow: "0 2px 8px 0 " + (isDarkMode ? "rgba(0, 0, 0, 0.25)" : "rgba(29, 35, 41, 0.05)"),
            userSelect: "none",
            background: token?.colorBgContainer + " !important",
            ".ant-menu-item-icon": {
                color: isDarkMode ? "rgb(255,255,255,0.85)" : ""
            }
        },
        content: css`
            min-height: calc(100vh - 125px);
            width: 100%;
            height: 100%;

            > div > div > div > div:first-of-type {
                width: 80%;
                margin-left: 10%;
            }

            ${responsive.md} {
                > div > div > div > div:first-of-type {
                    width: 100%;
                    margin-left: 0;
                }
            }
        `,
        content1: {
            minHeight: "calc(100vh - 125px)",
            width: "100% !important",
            height: "100%",
        },
        footer: {
            textAlign: 'center',
        },
        flow: {
            display: "flex",
            justifyContent: "space-between",
            height: "55px",
        },
        logo: {
            display: "inline-flex",
            justifyContent: "center",
            alignItems: "center",
            width: "220px",
            height: "55px",
        },
        darkColor: {
            backgroundColor: "#001529 !important"
        },
        bread: {
            backgroundColor: token?.colorBgContainer,
            padding: "5px 15px 5px 15px",
            fontSize: "12px",
        },
        breadStyle: {
            color: "rgb(156, 156, 156)",
            cursor: "pointer"
        },
    };
});

export type LayoutProps = {
    ref?: React.Ref<SinKingRef>,
    loading?: boolean,
    breadCrumb?: boolean,
    menuCollapsedWidth?: Number,
    menuUnCollapsedWidth?: Number,
    menus?: any,
    onMenuClick?: (item: any) => void,
    onMenuBtnClick?: (state: boolean) => void,
    footer?: any,
    headerRight?: any,
    headerLeft?: any,
    headerHidden?: boolean,
    headerFixed?: boolean,
    onLogoClick?: () => void,
    collapsedLogo?: (isLight: boolean) => any,
    unCollapsedLogo?: (isLight: boolean) => any,
    menuBottomBtnIcon?: string,
    menuBottomBtnText?: any,
    onMenuBottomBtnClick?: () => void,
    layout?: string,
    menuTheme?: string,
    waterMark?: any,
};

/**
 * 验证码组件
 */
export interface SinKingRef {
    collapsedMenu?: () => void;//菜单折叠
    unCollapsedMenu?: () => void;//菜单展开
    toggleCollapsedMenu?: () => void;//菜单切换
}

const SinKing: React.FC<LayoutProps> = forwardRef<SinKingRef>((props: any, ref): any => {
    let {
        loading = false,
        breadCrumb = true,
        menus,
        onMenuClick,
        onMenuBtnClick,
        onLogoClick,
        collapsedLogo,
        unCollapsedLogo,
        headerRight,
        headerLeft,
        headerFixed,
        headerHidden = false,
        menuCollapsedWidth = 60,
        menuUnCollapsedWidth = 220,
        menuBottomBtnIcon = undefined,
        menuBottomBtnText = undefined,
        onMenuBottomBtnClick,
        layout = 'inline',
        menuTheme = "light",
        waterMark = undefined
    } = props;
    const systemTheme = useTheme();
    /*
     * 样式
     */
    const [collapsed, setCollapsed] = useState(false);
    const [open, setOpen] = useState(false);
    const {
        styles: {
            container,
            sider,
            content1,
            sticky,
            header,
            content,
            footer,
            body,
            drawMenu,
            menuBtn,
            flow,
            logo,
            darkColor,
            bread,
            breadStyle
        }
    } = useLayoutStyles();
    const {mobile, md} = useResponsive();
    const menuBtnOnClick = () => {
        let status: boolean;
        if ((layout == 'inline' && mobile) || (layout == 'horizontal' && !md)) {
            if (collapsed) {
                setCollapsed(false);
            }
            status = !open
            setOpen(status);
        } else {
            status = !collapsed
            setCollapsed(status);
        }
        onMenuBtnClick?.(status);
    }

    /**
     * 方法挂载
     */
    useImperativeHandle(ref, () => ({
        collapsedMenu: () => {
            setCollapsed(true);
            setOpen(false);
        },
        unCollapsedMenu: () => {
            setCollapsed(false);
            setOpen(false);
        },
        toggleCollapsedMenu: () => {
            const status = !collapsed;
            setCollapsed(status);
            setOpen(false);
            onMenuBtnClick?.(status);
        }
    }));

    /**
     * 面包屑
     */
    const [breadCrumbData, setBreadCrumb] = useState<any>([]);
    const location = useLocation();
    const match = useSelectedRoutes();
    const initBreadCrumb = () => {
        const items = getParentList(getAllMenuItems(false), match?.at(-1)?.route?.name);
        let temp = [{
            title: '首页',
            onClick: () => {
                historyPush(getFirstMenuWithoutChildren(getAllMenuItems(location?.pathname))?.name || "");
            },
            className: breadStyle,
        }];
        const onClick = (x: any) => {
            if (x?.children && x?.children?.length > 0) {
                historyPush(getFirstMenuWithoutChildren(x?.children)?.name || "");
            } else {
                historyPush(x?.name);
            }
        }
        items.map((x) => {
            temp.push({
                title: x?.label,
                onClick: () => {
                    onClick(x);
                },
                className: breadStyle,
            });
        });
        setBreadCrumb(temp);
    }
    useEffect(() => {
        if (breadCrumb) {
            initBreadCrumb();
        }
    }, [breadCrumb, location?.pathname]);

    /**
     * 获取面包屑
     */
    const getBreadCrumb = () => {
        if (match?.at(-1)?.route?.hideBreadCrumb) {
            return;
        }
        return breadCrumb && breadCrumbData?.length > 0 &&
            <Breadcrumb className={bread} items={breadCrumbData}/>
    }

    /**
     * 获取菜单主题颜色
     */
    const getColor = () => {
        const mode = systemTheme?.isDarkMode ? "light" : (menuTheme == "dark" ? menuTheme : "light");
        return !systemTheme?.isDarkMode && mode == "dark" ? ' ' + darkColor : '';
    }

    /**
     * 获取水印字体
     */
    const getWaterMaskFont = () => {
        return {
            color: systemTheme?.isDarkMode ? "rgba(255, 255, 255, 0.03)" : "rgba(0, 0, 0, 0.07)"
        };
    }

    /**
     * 获取菜单
     * @param mode 布局模式
     */
    const getSider = (mode) => {
        mode = mode == "horizontal" ? mode : "inline";
        menuTheme = menuTheme == "dark" ? menuTheme : "light";
        return <Sider layout={mode} theme={menuTheme} collapsed={collapsed}
                      onLogoClick={onLogoClick}
                      collapsedLogo={collapsedLogo}
                      unCollapsedLogo={unCollapsedLogo}
                      menuBottomBtnIcon={menuBottomBtnIcon}
                      menuBottomBtnText={menuBottomBtnText}
                      onMenuBottomBtnClick={onMenuBottomBtnClick}
                      menus={menus}
                      onMenuClick={(item) => {
                          setOpen(false);
                          onMenuClick?.(item);
                      }}/>;
    }

    /**
     * 获取outlet
     */
    const getOutlet = () => {
        return <Watermark gap={[200, 200]} font={getWaterMaskFont()} content={waterMark}>
            <Outlet/>
        </Watermark>;
    }

    /**
     * 移动端抽屉
     */
    const drawer = <Drawer placement="left" closable={false} open={open} width={menuUnCollapsedWidth}
                           classNames={{body: drawMenu}}
                           onClose={() => {
                               setOpen(false)
                           }}>
        {getSider("inline")}
    </Drawer>;

    /**
     * 左右模式
     */
    const LayoutNormal = <Layout className={container}>
        <Layout.Sider className={sider} trigger={null} collapsible collapsed={collapsed}
                      width={menuUnCollapsedWidth} collapsedWidth={menuCollapsedWidth}
                      hidden={mobile}>
            {(mobile && drawer)
                ||
                getSider(layout)}
        </Layout.Sider>
        <Layout className={body}>
            <Layout.Header hidden={headerHidden} className={header + " " + (headerFixed ? sticky : "")}>
                <Header left={<div><Button type="text" size={"large"}
                                           icon={collapsed ? <MenuUnfoldOutlined/> : <MenuFoldOutlined/>}
                                           onClick={menuBtnOnClick} className={menuBtn}/>{headerLeft}</div>}
                        right={headerRight}/>
            </Layout.Header>
            <Layout.Content className={content1}>
                {getBreadCrumb()}
                {getOutlet()}
            </Layout.Content>
            {props?.footer && <Layout.Footer className={footer}>
                <Footer> {props?.footer}</Footer>
            </Layout.Footer>}
        </Layout>
    </Layout>;

    /**
     * 上下模式
     */
    const LayoutFlow = <Layout>
        <Layout.Header className={header + " " + (headerFixed ? sticky : "")}>
            {!md && <Header left={<div>
                <Button type="text" size={"large"} icon={collapsed ? <MenuUnfoldOutlined/> : <MenuFoldOutlined/>}
                        onClick={menuBtnOnClick} className={menuBtn}/>
                {headerLeft}
            </div>} right={headerRight}/>}
            {(!md && drawer) ||
                <div className={flow + getColor()}>
                    <div className={logo}>
                        {unCollapsedLogo?.(!systemTheme?.isDarkMode)}
                    </div>
                    {getSider(layout)}
                    <div>{headerRight}</div>
                </div>
            }
        </Layout.Header>
        <Layout.Content className={content}>
            {getBreadCrumb()}
            {getOutlet()}
        </Layout.Content>
        {props?.footer && <Layout.Footer className={footer}>
            <Footer> {props?.footer}</Footer>
        </Layout.Footer>}
    </Layout>;

    return (
        <ConfigProvider locale={zhCN}>
            <App>
                {(loading && <Loading/>) || (layout != "horizontal" ? LayoutNormal : LayoutFlow)}
            </App>
        </ConfigProvider>
    );

});
export default SinKing;