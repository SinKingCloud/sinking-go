import React, {useMemo, useState} from "react";
import {Layout, Icon, useTheme} from "sinking-antd";
import {useModel, useSelectedRoutes, useLocation, history, Outlet} from "umi";
import {deleteHeader} from "@/utils/auth";
import {getAllMenuItems, getFirstMenuWithoutChildren, getParentList, historyPush} from "@/utils/route";
import {App, Popover, Tooltip} from "antd";
import {createStyles} from "antd-style";
import Settings from "@/../config/defaultSettings";
import {logout} from "@/service/auth/login";
import request from "@/utils/request";
import Title from "../title";
import defaultSettings from "@/../config/defaultSettings";

/**
 * 中间件
 * @param ctx context
 * @param next 执行函数
 */
const check = async (ctx: any, next: any) => {
    await next();
    if (ctx.res.code == 401) {
        historyPush("notAllowed");
    }
}
request.use(check);

/**
 * 样式
 */
const useRightTopStyles = createStyles(({css, token, isDarkMode}): any => {
    return {
        nickname: {
            marginLeft: "3px",
            fontSize: "13px",
            color: isDarkMode ? token.colorTextSecondary : "rgb(150,150,150)",
            fontWeight: "bold",
        },
        bottomIconDark: {
            color: isDarkMode ? token.colorTextSecondary : "rgb(150,150,150)",
        },
        pop: css`
            margin-left: 10px;
            margin-right: 20px;
            display: inline-flex;
            align-items: center;
            padding: 11px 5px;
            border: 0;
            border-radius: 10px;
            background: transparent;
            font: inherit;
            line-height: 20px;
            white-space: nowrap;
            vertical-align: middle;
            transition: background-color 0.3s ease;
            cursor: pointer;

            &:focus-visible {
                outline: 2px solid ${token.colorPrimary};
                outline-offset: 2px;
            }

            .anticon {
                margin-left: 2px;
                font-size: 10px;
            }
        `,
        box: css`
            &.ant-popover {
                filter: none !important;
            }

            .ant-popover-container {
                padding: 0 !important;
                overflow: hidden;
                border-radius: ${token.borderRadiusLG}px;
                box-shadow: 0 0 5px 0 rgba(0, 0, 0, 0.15) !important;
                width: 210px;
            }

            .ant-popover-arrow:before {
                background-color: ${token?.colorPrimary} !important;
            }
        `,
        content_top: css`
            height: 64px;
            box-sizing: border-box;
            background-color: ${token?.colorPrimary};
            overflow: hidden;
            background-image: url(data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAIgAAACGBAMAAAD0nt8RAAAAD1BMVEVHcEz///////////////8T4DEaAAAABXRSTlMADAYJA8T7L0gAAALASURBVGjezVpbcoMwDCS2DxBhDgC0BwhNDoDb3v9MfaQQMLb1cqfVVz6YjbS7ks2IptnFCOD7RhcDfIafVRgGvqPVJ/IZmoLcD4YqFbuAgAIkrCCKeqYaICsGnP8WxP0bkE05F71hK4EoWvC0YHh9EwN0v93Fr9frs3aefP/PjdA9Pfo3F8xuxRk74LlaTBpLaPOATaRA0A+lBPAB6rBUNy2KVXwmDOF8YwSs3g1Ij8zYVgPi0PYjlGOppJVATmi9BgcZUJDH1PKokzIamwdGVkGDPGEeB2S+jU9lT2/4KFASip4edxgtfpxDidJiIpvmOsoTYfQEI8UPmR3GOFJuOLHGe0o97YYT8fYKNE4jSnYP7mUpH2x292SW0vKtIyQlNAeM4pEzpTQ0E3BAXOpJA4mYqZTcCx+BCRKOOg5JDE+m5AskjVECcfGzJoNRsokFanRkSgpxJlNSiAudEgmIoYPM+AVWAzKQMTzlzUAOwqCkq0FJV8FqBa9NFWzCsBrt2BLbhEEJ1KDE16DE16Ckq0FJV4GSrNc4lGS9xqEk6zUOJVmvcSgB/UAq2GSqAWJq2IQxpQsgpsZcY4h89Jpbf0xSr5kngBeuaSOv3f/9xuQ2LUnP4tanFWlZqfiMID1nHnQZry/Kv/FB3PF9giTz9foyju/vc6Spl7TQHW5IqDaAPBphKhn/hBogbqoAwpv7WRNKue2kkzI/ZoSpzPIjtTWZ1+03BkgfzRPJKdQupmjlo98v1hoVp9AyiGbFKdRSN0Wy1x56C6FLJKsuhiYzYZtleL0iu2zQtoROXQwqM3nlMmmLQWRm7BkHfSJ5mXkL6aAuJiszd5FsBc1Lkpm/ATbqYpIyS1bRcSo35U5dseLfT0rpXt3qE6GtUTjcKj4SWFNRfcPh9Kvs9bLRN7qY1B+k/HCrTeSL24saozF4MR+gwScpYNNOxQAAAABJRU5ErkJggg==);
            background-repeat: no-repeat;
            background-position: right;
            display: flex;
            align-items: center;
            padding: 0 15px;
        `,
        top_text: {
            color: "#fff",
            width: "100%",
            minWidth: 0,
        },
        userName: {
            fontSize: "14px",
            fontWeight: 600,
            lineHeight: "20px",
            whiteSpace: "nowrap",
            overflow: "hidden",
            textOverflow: "ellipsis",
        },
        userIp: {
            marginTop: "3px",
            fontSize: "12px",
            lineHeight: "18px",
            opacity: 0.8,
            whiteSpace: "nowrap",
            overflow: "hidden",
            textOverflow: "ellipsis",
        },
        menuItemLabel: {
            display: "flex",
            alignItems: "center",
        },
        menu: {
            listStyle: "none",
            padding: 0,
            margin: 0,
            userSelect: "none",
            "li:last-of-type": {
                borderTop: "0.5px solid rgba(189, 189, 189, 0.2)",
                button: {
                    height: "45px",
                    lineHeight: "45px",
                },
            }
        },
        menuItem: {
            width: "100%",
            border: 0,
            background: "transparent",
            fontFamily: "inherit",
            textAlign: "left",
            cursor: "pointer",
            letterSpacing: "1px",
            height: "40px",
            lineHeight: "40px",
            fontSize: "12px",
            padding: "0px 15px",
            transition: "background-color 0.3s ease",
            color: isDarkMode ? token.colorTextSecondary : "rgba(0,0,0,0.65)",
            display: "flex",
            alignItems: "center",
            justifyContent: "space-between",
            ":hover": {
                backgroundColor: "rgba(0, 0, 0, 0.03)",
            },
            ":focus-visible": {
                outline: "2px solid " + token.colorPrimary,
                outlineOffset: "-2px",
            },
            ".anticon": {
                fontSize: "11px",
            },
            "span>.anticon": {
                fontSize: "12.5px",
                marginRight: "7px"
            }
        },
        icon: {
            fontSize: "17px",
            padding: "7px",
            marginRight: "5px",
            cursor: "pointer",
            verticalAlign: "middle",
            borderRadius: "5px",
            transition: "background-color 0.3s ease",
            ":hover": {
                backgroundColor: "rgba(0, 0, 0, 0.1)",
            },
            color: isDarkMode ? token.colorTextSecondary : "rgb(150,150,150)"
        },
    };
});

/**
 * 右侧部分组件
 * @constructor
 */
const RightTop: React.FC = () => {
    /**
     * 全局数据
     */
    const user = useModel("user");//用户信息
    const theme = useTheme();//主题信息
    const {message, modal} = App.useApp();
    const [userOpen, setUserOpen] = useState(false);

    /**
     * 退出登录
     */
    const outLogin = async () => {
        message?.loading({content: "正在退出登录", duration: 600000, key: "outLogin"});
        await logout({
            onSuccess: (r) => {
                message?.success(r?.message || "退出登录成功");
                deleteHeader();
                user?.setWeb(undefined);
                historyPush("login");
            },
            onFail: (r) => {
                message?.error(r?.message || "退出登录失败");
            },
            onFinally: () => {
                message?.destroy("outLogin");
            },
        });
    };

    /**
     * 退出登录确认
     */
    const confirmOutLogin = () => {
        modal.confirm({
            title: "退出登录",
            content: "确定退出当前账号吗？",
            okText: "确定",
            okButtonProps: {danger: true, type: "default"},
            cancelText: "取消",
            mask: {
                closable: true,
            },
            onOk: outLogin,
        } as any);
    };

    const {
        styles: {
            nickname,
            bottomIconDark,
            pop,
            content_top,
            top_text,
            userName,
            userIp,
            box,
            menu,
            menuItem,
            menuItemLabel,
            icon
        }
    } = useRightTopStyles();
    return <>
        <Tooltip title={theme?.getModeName(theme?.mode as any)}>
            <Icon type={theme?.isDarkMode() ? "icon-dark" : (theme?.isAutoMode() ? "icon-auto" : "icon-light")}
                  className={icon}
                  onClick={() => {
                      theme?.toggle?.();
                  }}/>
        </Tooltip>
        <Popover rootClassName={box} autoAdjustOverflow={false}
                 open={userOpen}
                 onOpenChange={setUserOpen}
                 trigger={["hover", "click"]}
                 placement="bottomRight"
                 content={<>
                     <div className={content_top}>
                         <div className={top_text}>
                             <div className={userName}>{user?.web?.account || "未登录"}</div>
                             <div className={userIp}>{user?.web?.login_ip || "未知IP"}</div>
                         </div>
                     </div>
                     <ul className={menu} onKeyDown={(event) => {
                         if (event.key === "Escape") setUserOpen(false);
                     }}>
                         <li>
                             <button type="button" className={menuItem} onClick={() => {
                                 historyPush("system");
                             }}>
                                 <span className={menuItemLabel}><Icon type="SettingOutlined"/>系统管理</span>
                                 <Icon type="RightOutlined"/>
                             </button>
                         </li>
                         <li>
                             <button type="button" className={menuItem} onClick={() => {
                                 historyPush("log");
                             }}>
                                 <span className={menuItemLabel}><Icon type="FileTextOutlined"/>操作日志</span>
                                 <Icon type="RightOutlined"/>
                             </button>
                         </li>
                         <li>
                             <button type="button" className={menuItem} onClick={confirmOutLogin}>
                                 <span className={menuItemLabel}><Icon type="LogoutOutlined"/>退出登录</span>
                                 <Icon type="RightOutlined"/>
                             </button>
                         </li>
                     </ul>
                 </>}>
            <button type="button" className={pop} aria-label="账户菜单" aria-expanded={userOpen}
                    onKeyDown={(event) => {
                        if (event.key === "Escape") setUserOpen(false);
                    }}>
                <span className={nickname}>{user?.web?.account || "未登录"}</span>
                <Icon className={bottomIconDark} type="DownOutlined"/>
            </button>
        </Popover>
    </>
}

export type slide = {
    menu?: any,
}

/**
 * 样式信息
 */
const useSKLayoutStyles = createStyles((): any => {
    return {
        collapsedImg: {
            fontSize: "27px",
        },
        unCollapsed: {
            overflow: "hidden",
            position: "absolute",
            display: "inline-flex",
            ">span": {
                fontSize: "27px",
            },
            ">div": {
                fontSize: "25px",
                marginLeft: "5px",
                fontWeight: "bolder",
                float: "left",
                lineHeight: "30px",
                whiteSpace: "nowrap",
            }
        },
    };
});

/**
 * 用户系统
 */
const SKLayout: React.FC<slide> = ({...props}) => {
    /**
     * 初始化用户信息
     */
    const {menu} = props;
    /**
     * 全局信息
     */
    const user = useModel("user");
    const web = useModel("web");
    const location = useLocation();
    const match = useSelectedRoutes();
    const {styles: {collapsedImg, unCollapsed}} = useSKLayoutStyles();

    /**
     * 计算面包屑数据
     */
    const breadCrumbItems = useMemo(() => {
        if (!location?.pathname) return [];

        const items = getParentList(getAllMenuItems(false), match?.at(-1)?.route?.name);
        const temp = [{
            title: "首页",
            onClick: () => {
                historyPush(getFirstMenuWithoutChildren(getAllMenuItems(location?.pathname))?.name || "");
            },
        }];

        const handleItemClick = (x: any) => {
            if (x?.children && x?.children?.length > 0) {
                historyPush(getFirstMenuWithoutChildren(x?.children)?.name || "");
            } else {
                historyPush(x?.name);
            }
        };

        items.forEach((x) => {
            temp.push({
                title: x?.label,
                onClick: () => handleItemClick(x),
            });
        });

        return temp;
    }, [location?.pathname, match]);
    return (
        <>
            <Title/>
            <Layout loading={!user?.web}
                    pathname={location?.pathname}
                    matchedRoutes={match || []}
                    onNavigate={(path) => history.push(path)}
                    breadCrumbItems={breadCrumbItems}
                    hideBreadCrumb={match?.at(-1)?.route?.hideBreadCrumb}
                    waterMark={web?.info?.ui?.watermark ? [web?.info?.name, user?.web?.account] : ""}
                    menus={menu}
                    layout={web?.info?.ui?.layout != "left" ? "horizontal" : "inline"}
                    flowLayout={web?.info?.ui?.layout != "left"}
                    menuTheme={web?.info?.ui?.theme == "dark" ? "dark" : "light"}
                    footer={<>©{new Date().getFullYear()} All Right
                        Revered {web?.info?.name || Settings?.title}</>}
                    headerHidden={false}
                    headerFixed={true}
                    headerRight={<RightTop/>}
                    menuCollapsedWidth={60}
                    menuUnCollapsedWidth={210}
                    collapsedLogo={() => {
                        return <Icon type={"icon-logo"} style={{color: web?.info?.ui?.color || defaultSettings?.color}}
                                     className={collapsedImg}/>;
                    }}
                    unCollapsedLogo={() => {
                        return (
                            <div className={unCollapsed}>
                                <Icon type={"icon-logo"}
                                      style={{color: web?.info?.ui?.color || defaultSettings?.color}}/>
                                <div style={{color: web?.info?.ui?.color || defaultSettings?.color}}>
                                    {web?.info?.name || Settings?.title}
                                </div>
                            </div>)
                    }}>
                <Outlet/>
            </Layout>
        </>
    );
}
export default SKLayout;
