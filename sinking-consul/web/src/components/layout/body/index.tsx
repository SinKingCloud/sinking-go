import React from "react";
import {ConfigProvider, Layout} from "antd";
import {createStyles} from "antd-style";
import {Spin, Space, App} from "antd";
import zhCN from 'antd/locale/zh_CN';
import {Animation, Theme} from "@/components";
import {Animate} from "@/components/animation";

const useStyles = createStyles(({css}): any => {
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
    };
});

export type BodyProps = {
    loading?: boolean;//是否加载状态
    space?: boolean;//是否开启间距
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
const Body: React.FC<BodyProps> = React.memo((props) => {
    const {
        children,
        loading,
        style,
        className,
        space = true,
        themes = undefined,
        mode = undefined,
        animation = true,
    } = props;
    const {styles: {body, load, gutter}, cx} = useStyles();

    /**
     * 内容
     */
    const content = <ConfigProvider locale={zhCN}>
        <App>
            {loading ? (
                <Spin spinning={true} size="large" className={load}></Spin>
            ) : (
                <Layout style={style}>
                    <div className={cx("ant-layout-body", className || body)}>
                        <Animation animate={animation ? Animate.FadeUp : Animate.None}>
                            {space ? (
                                <Space direction="vertical" size="middle" className={gutter}>
                                    {children}
                                </Space>
                            ) : children}
                        </Animation>
                    </div>
                </Layout>
            )}
        </App>
    </ConfigProvider>;

    /**
     * 页面容器
     */
    return <Theme theme={themes} mode={mode}>
        {content}
    </Theme>;
});

export default Body;