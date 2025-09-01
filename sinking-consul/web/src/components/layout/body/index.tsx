import React, {useEffect, useState} from "react";
import {useLocation, useRouteData} from "umi";
import {ConfigProvider, Layout} from "antd";
import {createStyles} from "antd-style";
import {Spin, Space, Breadcrumb, App} from "antd";
import {getCurrentMenus, getFirstMenuWithoutChildren, getParentList, historyPush} from "@/utils/route";
import zhCN from 'antd/locale/zh_CN';
import {Animation, Theme} from "@/components";
import {Animate} from "@/components/animation";

const useStyles = createStyles(({token, css}): any => {
    return {
        body: css`
            padding: 10px;
        `,
        load: {
            margin: "0 auto",
            width: "100%",
            lineHeight: "80vh",
        },
        gutter: {
            display: "flex"
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

export type BodyProps = {
    loading?: boolean;//是否加载状态
    space?: boolean;//是否开启间距
    breadCrumb?: boolean;//面包屑
    animation?: boolean;//动画
    style?: any;//样式
    className?: any;//样式名
    children?: any;//子内容
    themes?: any;//主题
    mode?: any;//模式
};
/**
 * 页面主体部分
 * @param props
 * @constructor
 */
const Body: React.FC<BodyProps> = (props) => {
    const {
        children,
        loading,
        style,
        className,
        space = true,
        breadCrumb = true,
        themes = undefined,
        mode = undefined,
        animation = true,
    } = props;
    const {styles: {body, load, gutter, bread, breadStyle}} = useStyles();

    /**
     * 初始化面包屑
     */
    const [breadCrumbData, setBreadCrumb] = useState<any>([]);
    const location = useLocation();
    const match = useRouteData();
    const initBreadCrumb = () => {
        const items = getParentList(getCurrentMenus(location?.pathname, false), match?.route?.name);
        let temp = [{
            title: '首页',
            onClick: () => {
                historyPush(getFirstMenuWithoutChildren(getCurrentMenus(location?.pathname))?.name || "");
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

    /**
     * 初始化
     */
    useEffect(() => {
        if (breadCrumb) {
            initBreadCrumb();
        }
    }, [breadCrumb]);

    /**
     * 内容
     */
    const content = <ConfigProvider locale={zhCN}>
        <App>
            {(loading && <Spin spinning={true} size="large" className={load}></Spin>) ||
                <Layout style={style}>
                    {breadCrumb && breadCrumbData?.length > 0 &&
                        <Breadcrumb className={bread} items={breadCrumbData}/>}
                    <div className={className ? className : body}>
                        <Animation animate={animation ? Animate.FadeUp : Animate.None}>
                            {(space && <Space direction="vertical" size="middle" className={gutter}>
                                {children}
                            </Space>) || children}
                        </Animation>
                    </div>
                </Layout>}
        </App>
    </ConfigProvider>;

    /**
     * 页面容器
     */
    return <Theme theme={themes} mode={mode}>
        {content}
    </Theme>;
}

export default Body