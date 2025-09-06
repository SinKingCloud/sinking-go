import React, {useEffect, useState} from 'react';
import {App, ColorPicker, Form, Spin, Upload, Select, Input, Button} from "antd";
import {getConfig, setConfig} from "@/service/admin/system";
import {LoadingOutlined, PlusOutlined} from "@ant-design/icons";
import {createStyles} from "antd-style";
import {useModel} from "umi";


const UiView: React.FC = () => {
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
                group: "ui"
            },
            onSuccess: (r: any) => {
                form.setFieldsValue(r?.data || {});
            },
            onFail: (r: any) => {
                message?.error(r?.message || "加载配置失败")
            },
            onFinally: () => {
                setIsLoading(false)
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
                setIsLoading(false);
            }
        });
    }


    /**
     * 初始化数据
     */
    useEffect(() => {
        getConfigs();
    }, []);

    return (
        <Spin spinning={isLoading} size="default">
            <div style={{display: isLoading ? 'none' : 'block'}}>
                <Form form={form} onFinish={onFinish} layout="vertical">
                    <Form.Item
                        name="compact"
                        label="紧凑模式"
                        tooltip="网站界面紧凑模式"
                        style={{maxWidth: '400px', width: '100%'}}
                    >
                        <Select
                            placeholder="请选择紧凑模式是否开启"
                            options={[
                                {value: '1', label: '开启'},
                                {value: '0', label: '关闭'}
                            ]}
                        />
                    </Form.Item>
                    <Form.Item
                        name="layout"
                        label="网站布局"
                        tooltip="网站的整体布局"
                        style={{maxWidth: '400px', width: '100%'}}
                    >
                        <Select
                            placeholder="请选择网站布局"
                            options={[
                                {value: 'top', label: '上下布局'},
                                {value: 'left', label: '左右布局'}
                            ]}
                        />
                    </Form.Item>
                    <Form.Item
                        name="theme"
                        label="菜单主题"
                        tooltip="网站的主题颜色"
                        style={{maxWidth: '400px', width: '100%'}}
                    >
                        <Select
                            placeholder="请选择菜单主题"
                            options={[
                                {value: 'light', label: '亮色模式'},
                                {value: 'dark', label: '暗色模式'}
                            ]}
                        />
                    </Form.Item>
                    <Form.Item
                        name="watermark"
                        label="界面水印"
                        tooltip="系统界面是否显示用户昵称水印"
                        style={{maxWidth: '400px', width: '100%'}}
                    >
                        <Select
                            placeholder="请选择是否打开界面水印"
                            options={[
                                {value: '1', label: '开启'},
                                {value: '0', label: '关闭'}
                            ]}
                        />
                    </Form.Item>
                    <Form.Item
                        name="color"
                        label="主题颜色"
                        tooltip="网站的主题颜色"
                        style={{maxWidth: '400px', width: '100%'}}
                        getValueFromEvent={(color) => color?.toRgbString()}
                    >
                        <ColorPicker
                            format="rgb"
                            defaultFormat="rgb"
                        />
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

export default UiView;
