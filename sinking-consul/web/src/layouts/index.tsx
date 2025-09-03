import React, {useEffect, useState} from "react";
import {App, Spin} from "antd";
import {getAllMenuItems} from "@/utils/route";
import {Layout} from "@/layouts/components";
import {Outlet, useModel, useSelectedRoutes} from "umi";
import {createStyles} from "antd-style";
import {Theme} from "@/components";
import Title from "./components/title";

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
    const routes = useSelectedRoutes();
    const {styles: {load}} = useStyles();

    const [menu, setMenu] = useState([]);
    useEffect(() => {
        setMenu(getAllMenuItems(true));
    }, []);

    const getDom = () => {
        if (!web?.info || !web?.info?.ui) {
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