import Footer from "../footer";
import Header from "../header";
import Sider from "../sider";
import BreadCrumb from "../../breadcrumb";
import {createStyles, useResponsive, useTheme} from "antd-style";
import React, {forwardRef, useImperativeHandle, useState} from "react";
import {Outlet} from "umi";
import {MenuFoldOutlined, MenuUnfoldOutlined} from "@ant-design/icons";
import Loading from "@/loading"
import {App, Button, ConfigProvider, Drawer, Layout, Watermark} from "antd";
import zhCN from 'antd/locale/zh_CN';

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
        body: css`
            transition: margin-left 0.3s !important;
            background-color: ${isDarkMode ? "black" : "transparent"};
            display: flex;
            flex-direction: column;
            min-height: 100vh;
        `,
        sticky: {
            position: "sticky",
            top: 0,
            left: 0,
        },
        header: css`
            padding: 0;
            height: 55px;
            line-height: 55px;
            width: 100%;
            z-index: 3;
            box-shadow: 0 2px 8px 0 ${isDarkMode ? "rgba(0, 0, 0, 0.25)" : "rgba(29, 35, 41, 0.05)"};
            user-select: none;
            background: ${token?.colorBgContainer} !important;
            flex-shrink: 0;

            .ant-menu-item-icon {
                color: ${isDarkMode ? "rgb(255,255,255,0.85)" : ""};
            }
        `,
        horizontalContent: css`
            width: 100%;
            flex: 1 0 auto;
        `,
        inlineContent: css`
            width: 100%;
            flex: 1 0 auto;
        `,
        flowContent: css`
            .ant-layout-body {
                width: 80%;
                margin-left: 10%;
            }

            ${responsive.md} {
                .ant-layout-body {
                    width: 100%;
                    margin-left: 0;
                }
            }
        `,
        footer: css`
            text-align: center;
            flex-shrink: 0;
        `,
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
        layoutNormal: css`
            min-height: 100vh;
            display: flex;
            flex-direction: column;
        `,
    };
});

export type LayoutProps = {
    ref?: React.Ref<SinKingRef>,
    loading?: boolean,
    breadCrumb?: boolean,
    menuCollapsedWidth?: number,
    menuUnCollapsedWidth?: number,
    menus?: any,
    menuClassNames?: any,
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
    menuBottomBtnIcon?: any,
    menuBottomBtnText?: any,
    onMenuBottomBtnClick?: () => void,
    layout?: string,
    flowLayout?: boolean,
    menuTheme?: string,
    waterMark?: any,
};

/**
 * 布局组件引用接口
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
        menuClassNames,
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
        flowLayout = false,
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
            inlineContent,
            sticky,
            header,
            horizontalContent,
            footer,
            body,
            drawMenu,
            menuBtn,
            flow,
            logo,
            darkColor,
            flowContent,
            layoutNormal
        },
        cx
    } = useLayoutStyles();
    const {mobile, md} = useResponsive();
    const menuBtnOnClick = () => {
        let status: boolean;
        if ((layout === 'inline' && mobile) || (layout === 'horizontal' && !md)) {
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
     * 获取菜单主题颜色
     */
    const getColor = () => {
        const mode = systemTheme?.isDarkMode ? "light" : (menuTheme === "dark" ? menuTheme : "light");
        return !systemTheme?.isDarkMode && mode === "dark" ? ' ' + darkColor : '';
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
    const getSider = (mode: string) => {
        const layoutMode = mode === "horizontal" ? mode : "inline";
        const theme = menuTheme === "dark" ? menuTheme : "light";
        return <Sider classNames={menuClassNames} layout={layoutMode} theme={theme} collapsed={collapsed}
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
        {mobile && drawer}
        <Layout.Sider className={sider} trigger={null} collapsible collapsed={collapsed}
                      width={menuUnCollapsedWidth} collapsedWidth={menuCollapsedWidth}
                      hidden={mobile}>
            {getSider(layout)}
        </Layout.Sider>
        <Layout className={body}>
            {!headerHidden && <Layout.Header className={cx(header, headerFixed && sticky)}>
                <Header left={<div><Button type="text" size={"large"}
                                           icon={collapsed ? <MenuUnfoldOutlined/> : <MenuFoldOutlined/>}
                                           onClick={menuBtnOnClick} className={menuBtn}/>{headerLeft}</div>}
                        right={headerRight}/>
            </Layout.Header>}
            <Layout.Content className={cx(inlineContent, flowLayout && flowContent)}>
                <BreadCrumb enabled={breadCrumb}/>
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
    const LayoutFlow = <Layout className={layoutNormal}>
        <Layout.Header className={cx(header, headerFixed && sticky)}>
            {!md && <Header left={<div>
                <Button type="text" size={"large"} icon={collapsed ? <MenuUnfoldOutlined/> : <MenuFoldOutlined/>}
                        onClick={menuBtnOnClick} className={menuBtn}/>
                {headerLeft}
            </div>} right={headerRight}/>}
            {(!md && drawer) ||
                <div className={cx(flow, getColor())}>
                    <div className={logo}>
                        {unCollapsedLogo?.(!systemTheme?.isDarkMode)}
                    </div>
                    {getSider(layout)}
                    <div>{headerRight}</div>
                </div>
            }
        </Layout.Header>
        <Layout.Content className={cx(horizontalContent, flowLayout && flowContent)}>
            <BreadCrumb enabled={breadCrumb}/>
            {getOutlet()}
        </Layout.Content>
        {props?.footer && <Layout.Footer className={footer}>
            <Footer> {props?.footer}</Footer>
        </Layout.Footer>}
    </Layout>;

    return (
        <ConfigProvider locale={zhCN}>
            <App>
                {(loading && <Loading/>) || (layout !== "horizontal" ? LayoutNormal : LayoutFlow)}
            </App>
        </ConfigProvider>
    );
});
export default SinKing;