import React, {useEffect, useState} from "react";
import {Outlet, useLocation, history, useSelectedRoutes} from "umi";
import {App, ConfigProvider, Spin} from "antd";
import {getAllMenuItems, historyPush} from "@/utils/route";
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
    ctx.req.options.headers = getHeaders();
    await next();
    if (ctx.res.code == 403 || ctx.res.code == 503) {
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
            margin: "0 auto",
            width: "100%",
            lineHeight: "80vh",
        },
    };
});

const ProLayout = () => {
    const web = useModel("web");
    const user = useModel("user");
    const routes = useSelectedRoutes();
    const {styles: {load}} = useStyles();

    const [menu, setMenu] = useState([]);
    const [loading, setLoading] = useState(true);

    // 初始化 iconfont 地址
    useEffect(() => {
        if (defaultSettings?.iconfontUrl) {
            setIconfontUrl(defaultSettings.iconfontUrl);
        }
    }, []);

    useEffect(() => {
        setMenu(getAllMenuItems(true));
        setLoading(true)
        user?.refreshWebUser(() => {
            setLoading(false)
        });
    }, []);

    const theme = useTheme();
    useEffect(() => {
        if (web?.info?.ui) {
            if (web.info.ui.color) {
                theme?.setColor(web.info.ui.color);
            }
            if (web.info.ui.radius >= 0) {
                theme?.setRadius(web.info.ui.radius <= 15 ? web.info.ui.radius : 0);
            }
            if (web.info.ui.compact) {
                theme?.setCompactTheme();
            } else {
                theme?.setDefaultTheme();
            }
        }
    }, [web?.info?.ui]);

    if (loading || !web?.info || !web?.info?.ui) {
        return <Spin spinning={true} size="large" className={load}/>;
    }
    if (routes?.pop()?.route?.auth === false) {
        return <>
            <Title/>
            <App>
                <Outlet/>
            </App>
        </>;
    }
    return <Layout menu={menu}/>;
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