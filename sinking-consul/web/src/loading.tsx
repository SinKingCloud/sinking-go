import React from "react";
import {createStyles} from "antd-style";
import {App, ConfigProvider, Spin} from "antd";
import zhCN from 'antd/locale/zh_CN';
import 'dayjs/locale/zh-cn';

const useStyles = createStyles((): any => {
    return {
        body: {
            width: "100%",
            margin: "0 auto",
            lineHeight: "80vh"
        },
    };
});

export default () => {
    const {styles: {body}} = useStyles();
    return (
        <ConfigProvider locale={zhCN}>
            <App>
                <Spin size="large" className={body}/>
            </App>
        </ConfigProvider>
    );
};