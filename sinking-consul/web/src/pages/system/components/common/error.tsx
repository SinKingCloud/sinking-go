import React from "react";
import {Button} from "antd";
import {Icon} from "sinking-antd";

export default ({styles, message, onRetry}: any): React.ReactNode => (
    <div className={styles.errorState}>
        <Icon type="ExclamationCircleOutlined"/>
        <span>{message || "配置加载失败"}</span>
        <Button size="small" icon={<Icon type="ReloadOutlined"/>} onClick={onRetry}>
            重新加载
        </Button>
    </div>
);
