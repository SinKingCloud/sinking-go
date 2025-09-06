import React from 'react';
import {useSystemConfig} from "../hooks/useSystemConfig";
import BaseForm from "./common/BaseForm";
import FormField from "./common/FormField";

const WebView: React.FC = () => {
    const {form, dataLoading, submitLoading, saveConfig, resetForm} = useSystemConfig({
        group: "web"
    });

    return (
        <BaseForm
            form={form}
            loading={dataLoading}
            submitLoading={submitLoading}
            onFinish={saveConfig}
            onReset={resetForm}
        >
            <FormField
                type="input"
                name="name"
                label="网站名称"
                tooltip="网站的名称"
                placeholder="请输入网站名称"
            />
            <FormField
                type="input"
                name="title"
                label="网站标题"
                tooltip="网站页面的标题，显示在浏览器标签页"
                placeholder="请输入网站标题"
            />
            <FormField
                type="input"
                name="keywords"
                label="网站关键字"
                tooltip="用于SEO优化的关键字，多个关键字用逗号分隔"
                placeholder="请输入网站关键字"
            />
            <FormField
                type="textarea"
                name="describe"
                label="网站描述"
                tooltip="网站的描述信息，用于SEO优化"
                placeholder="请输入网站描述"
                rows={4}
            />
        </BaseForm>
    );
};

export default WebView;
