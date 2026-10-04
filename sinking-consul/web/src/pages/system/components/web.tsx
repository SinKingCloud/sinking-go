import React, {useMemo} from "react";
import {Form, Input} from "antd";
import {useSystemConfig} from "../hooks";
import {actions as Actions, error as Error, loading as Loading} from "./common";

export default ({styles, info}: any): React.ReactNode => {
    const defaults = useMemo(() => ({
        name: info?.name || "",
        title: info?.title || "",
        keywords: info?.keywords || "",
        describe: info?.describe || "",
    }), [info?.name, info?.title, info?.keywords, info?.describe]);
    const config = useSystemConfig({group: "web", defaults});

    return (
        <Form form={config.form} layout="vertical" variant="filled" onFinish={config.saveConfig}>
            {config.dataLoading ? <Loading styles={styles}/> : config.error ? (
                <Error styles={styles} message={config.error} onRetry={config.loadConfig}/>
            ) : <>
                <Form.Item name="name" label="网站名称" tooltip="网站的名称" className={styles.formField}>
                    <Input placeholder="请输入网站名称"/>
                </Form.Item>
                <Form.Item name="title" label="网站标题" tooltip="网站页面的标题，显示在浏览器标签页"
                           className={styles.formField}>
                    <Input placeholder="请输入网站标题"/>
                </Form.Item>
                <Form.Item name="keywords" label="网站关键字" tooltip="用于 SEO 优化的关键字，多个关键字用逗号分隔"
                           className={styles.formField}>
                    <Input placeholder="请输入网站关键字"/>
                </Form.Item>
                <Form.Item name="describe" label="网站描述" tooltip="网站的描述信息，用于 SEO 优化"
                           className={styles.formField}>
                    <Input.TextArea placeholder="请输入网站描述" rows={4}/>
                </Form.Item>
                <Actions styles={styles} saving={config.submitLoading} onReset={config.resetForm}/>
            </>}
        </Form>
    );
};
