import React, {useEffect, useLayoutEffect} from "react";
import {Outlet, useSelectedRoutes} from "umi";
import {App, ConfigProvider, Spin} from "antd";
import {historyPush} from "@/utils/route";
import {Layout} from "@/layouts/components";
import {useModel} from "umi";
import Title from "./components/title";
import {deleteHeader, getHeaders} from "@/utils/auth";
import request from "@/utils/request";
import defaultSettings from "../../config/defaultSettings";
import {setIconfontUrl, Theme, useTheme} from "sinking-antd";
import {createStyles} from "antd-style";
import zhCN from 'antd/locale/zh_CN';

/**
 * 中间件
 * @param ctx context
 * @param next 执行函数
 */
const check = async (ctx: any, next: any) => {
    ctx.req.options.headers = {...ctx.req.options.headers, ...getHeaders()};
    await next();
    if (ctx.res?.code == 403 || ctx.res?.code == 503) {
        historyPush("login");//登陆页面;
        deleteHeader()
    }
}
request.use(check);

/**
 * 样式信息
 */
const useStyles = createStyles((): any => {
    return {
        load: {
            width: "100%",
            display: "flex",
            justifyContent: "center",
            alignItems: "center",
            height: "80vh",
        },
    };
});

const ProLayout = () => {
    const web = useModel("web");
    const user = useModel("user");
    const routes = useSelectedRoutes();
    const {styles: {load}} = useStyles();

    const isPublic = routes?.at(-1)?.route?.auth === false;

    // 初始化 iconfont 地址
    useEffect(() => {
        if (defaultSettings?.iconfontUrl) {
            setIconfontUrl(defaultSettings.iconfontUrl);
        }
    }, []);

    useEffect(() => {
        if (!isPublic && !user?.web) user?.refreshWebUser();
    }, [isPublic]);

    const theme = useTheme();
    useLayoutEffect(() => {
        if (web?.info?.ui) {
            if (web?.info?.ui?.color) {
                theme?.setColor(web?.info?.ui?.color);
            }
            if (web?.info?.ui?.radius >= 0) {
                theme?.setRadius(web?.info?.ui?.radius <= 15 ? web?.info?.ui?.radius : 0);
            }
            if (web?.info?.ui?.compact) {
                theme?.setCompactTheme();
            } else {
                theme?.setDefaultTheme();
            }
        }
    }, [web?.info?.ui]);

    if (!web?.info?.ui) {
        return <Spin spinning={true} size="large" className={load}/>;
    }
    if (isPublic) {
        return <>
            <Title/>
            <App>
                <Outlet/>
            </App>
        </>;
    }
    if (!user?.web) return <Spin spinning size="large" className={load}/>;
    return <Layout/>;
}

export default () => {
    return (
        <ConfigProvider locale={zhCN}>
            <Theme>
                <ProLayout/>
            </Theme>
        </ConfigProvider>
    );
}