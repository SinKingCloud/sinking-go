import React, {useEffect, useState} from 'react';
import {getConfig, setConfig} from "@/service/admin/system";
import {App, Form, Spin, Input, Button} from "antd";
import {createStyles} from "antd-style";
import {useModel} from "@@/exports";

const useStyles = createStyles(({css}) => {
    return {
        box: css`
            .ant-form-item .ant-form-item-control {
                margin-bottom: 10px !important;
            }
        `
    }
})

const WebView: React.FC = () => {
    const {styles: {box}} = useStyles()

    const [isLoading, setIsLoading] = useState(false);
    const {message} = App.useApp()
    const [form] = Form.useForm();
    const web = useModel("web");

    /**
     * 初始化表单值
     */
    const getConfigs = async () => {
        setIsLoading(true);
        return await getConfig({
            body: {
                action: "get",
                group: "web"
            },
            onSuccess: (r: any) => {
                form.setFieldsValue(r?.data || {});
            },
            onFail: (r: any) => {
                message?.error(r?.message || "加载配置失败");
            },
            onFinally: () => {
                setIsLoading(false);
            }
        });
    }

    /**
     * 提交表单
     */
    const onFinish = async (values: any) => {
        setIsLoading(true);
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
                setIsLoading(false);
            }
        });
    }

    /**
     * 初始化数据
     */
    useEffect(() => {
        getConfigs()
    }, []);

    return (
        <Spin spinning={isLoading} size="default">
            <div style={{display: isLoading ? 'none' : 'block'}}>
                <Form form={form} onFinish={onFinish} className={box} layout="vertical">
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
                        <Button type="primary" htmlType="submit" loading={isLoading}>
                            保存配置
                        </Button>
                    </Form.Item>
                </Form>
            </div>
        </Spin>
    );
};

export default WebView;
