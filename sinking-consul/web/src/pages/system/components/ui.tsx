import React, {useEffect, useState} from 'react';
import {App, ColorPicker, Form, Spin, Select, Button} from "antd";
import {getConfig, setConfig} from "@/service/admin/system";
import {useModel} from "umi";

const UiView: React.FC = () => {
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
                group: "ui"
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
                message?.error(r?.message || "加载配置失败")
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
                group: "ui",
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
        getConfigs();
    }, []);

    return (
        <Spin spinning={dataLoading} size="default">
            <div style={{display: dataLoading ? 'none' : 'block'}}>
                <Form form={form} onFinish={onFinish} layout="vertical">
                    <Form.Item
                        name="compact"
                        label="紧凑模式"
                        tooltip="控制界面元素间距，开启后界面更紧凑"
                        style={{maxWidth: '400px', width: '100%'}}
                    >
                        <Select
                            placeholder="请选择紧凑模式是否开启"
                            allowClear
                            options={[
                                {value: '1', label: '开启'},
                                {value: '0', label: '关闭'}
                            ]}
                        />
                    </Form.Item>
                    <Form.Item
                        name="layout"
                        label="网站布局"
                        tooltip="选择网站的整体布局方式，上下布局或左右布局"
                        style={{maxWidth: '400px', width: '100%'}}
                    >
                        <Select
                            placeholder="请选择网站布局方式"
                            allowClear
                            options={[
                                {value: 'top', label: '上下布局'},
                                {value: 'left', label: '左右布局'}
                            ]}
                        />
                    </Form.Item>
                    <Form.Item
                        name="theme"
                        label="菜单主题"
                        tooltip="选择菜单的主题模式，亮色或暗色"
                        style={{maxWidth: '400px', width: '100%'}}
                    >
                        <Select
                            placeholder="请选择菜单主题模式"
                            allowClear
                            options={[
                                {value: 'light', label: '亮色模式'},
                                {value: 'dark', label: '暗色模式'}
                            ]}
                        />
                    </Form.Item>
                    <Form.Item
                        name="watermark"
                        label="界面水印"
                        tooltip="控制系统界面是否显示用户昵称水印效果"
                        style={{maxWidth: '400px', width: '100%'}}
                    >
                        <Select
                            placeholder="请选择是否显示界面水印"
                            allowClear
                            options={[
                                {value: '1', label: '开启'},
                                {value: '0', label: '关闭'}
                            ]}
                        />
                    </Form.Item>
                    <Form.Item
                        name="color"
                        label="主题颜色"
                        tooltip="选择网站的主题颜色，影响整体视觉风格"
                        style={{maxWidth: '400px', width: '100%'}}
                        getValueFromEvent={(color) => color?.toRgbString()}
                    >
                        <ColorPicker
                            format="rgb"
                            defaultFormat="rgb"
                        />
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

export default UiView;
