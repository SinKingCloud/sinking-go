import React, {useEffect, useState} from "react";
import {App, Spin} from "antd";
import {getAllMenuItems, historyPush} from "@/utils/route";
import {Layout} from "@/layouts/components";
import {Outlet, useModel, useSelectedRoutes} from "umi";
import {createStyles} from "antd-style";
import {Theme} from "@/components";
import Title from "./components/title";
import {deleteHeader, getHeaders} from "@/utils/auth";
import request from "@/utils/request";

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

export default () => {
    const web = useModel("web");
    const user = useModel("user");
    const routes = useSelectedRoutes();
    const {styles: {load}} = useStyles();

    const [menu, setMenu] = useState([]);
    const [loading, setLoading] = useState(true);
    useEffect(() => {
        setMenu(getAllMenuItems(true));
        setLoading(true)
        user?.refreshWebUser(() => {
            setLoading(false)
        });
    }, []);

    const getDom = () => {
        if (loading || !web?.info || !web?.info?.ui) {
            return <Spin spinning={true} size="large" className={load}/>;
        }
        if (routes?.pop()?.route?.auth === false) {
            return <Theme>
                <Title/>
                <App>
                    <Outlet/>
                </App>
            </Theme>;
        }
        return <Layout menu={menu}/>;
    }

    return getDom();
}