import React, {useEffect, useState} from 'react';
import {getConfig, setConfig} from "@/service/admin/system";
import {App, Form, Spin, Input, Button} from "antd";
import {useModel} from "umi";

const WebView: React.FC = () => {
    const [dataLoading, setDataLoading] = useState(true);
    const [submitLoading, setSubmitLoading] = useState(false);
    const {message} = App.useApp()
    const [form] = Form.useForm();
    const web = useModel("web");

    /**
     * 初始化表单值
     */
    const getConfigs = () => {
        setDataLoading(true);
        return getConfig({
            body: {
                action: "get",
                group: "web"
            },
            onSuccess: (r: any) => {
                const data = r?.data || {};
                // 处理空值，确保空字符串被转换为undefined以显示placeholder
                const processedData = Object.keys(data).reduce((acc, key) => {
                    const value = data[key];
                    acc[key] = value === '' || value === null ? undefined : value;
                    return acc;
                }, {} as any);
                form.setFieldsValue(processedData);
            },
            onFail: (r: any) => {
                message?.error(r?.message || "加载配置失败");
            },
            onFinally: () => {
                setDataLoading(false);
            }
        });
    }

    /**
     * 提交表单
     */
    const onFinish = async (values: any) => {
        setSubmitLoading(true);
        const configs = Object.entries(values).map(([key, value]) => ({key, value}));
        await setConfig({
            body: {
                action: "set",
                group: "web",
                configs
            },
            onSuccess: (r: any) => {
                message?.success(r?.message || "配置保存成功");
                web?.refreshInfo();
            },
            onFail: (r: any) => {
                message?.error(r?.message || "配置保存失败");
            },
            onFinally: () => {
                setSubmitLoading(false);
            }
        });
    }

    /**
     * 重置表单
     */
    const onReset = () => {
        form.resetFields();
    }

    /**
     * 初始化数据
     */
    useEffect(() => {
        getConfigs()
    }, []);

    return (
        <Spin spinning={dataLoading} size="default">
            <div style={{display: dataLoading ? 'none' : 'block'}}>
                <Form form={form} onFinish={onFinish} layout="vertical">
                    <Form.Item
                        name="name"
                        label="网站名称"
                        tooltip="网站的名称"
                        style={{maxWidth: '400px', width: '100%'}}
                    >
                        <Input placeholder="请输入网站名称"/>
                    </Form.Item>
                    <Form.Item
                        name="title"
                        label="网站标题"
                        tooltip="网站页面的标题，显示在浏览器标签页"
                        style={{maxWidth: '400px', width: '100%'}}
                    >
                        <Input placeholder="请输入网站标题"/>
                    </Form.Item>
                    <Form.Item
                        name="keywords"
                        label="网站关键字"
                        tooltip="用于SEO优化的关键字，多个关键字用逗号分隔"
                        style={{maxWidth: '400px', width: '100%'}}
                    >
                        <Input placeholder="请输入网站关键字"/>
                    </Form.Item>
                    <Form.Item
                        name="describe"
                        label="网站描述"
                        tooltip="网站的描述信息，用于SEO优化"
                        style={{maxWidth: '400px', width: '100%'}}
                    >
                        <Input.TextArea placeholder="请输入网站描述" rows={4}/>
                    </Form.Item>
                    <Form.Item>
                        <Button onClick={onReset} style={{marginRight: 8}}>
                            重置
                        </Button>
                        <Button type="primary" htmlType="submit" loading={submitLoading}>
                            提交
                        </Button>
                    </Form.Item>
                </Form>
            </div>
        </Spin>
    );
};

export default WebView;
