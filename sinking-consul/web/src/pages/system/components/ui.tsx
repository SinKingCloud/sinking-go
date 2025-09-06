import React from 'react';
import {useSystemConfig} from "@/pages/system/hooks";
import BaseForm from "./common/BaseForm";
import FormField from "./common/FormField";

const UiView: React.FC = () => {
    const {form, dataLoading, submitLoading, saveConfig, resetForm} = useSystemConfig({
        group: "ui"
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
                type="select"
                name="compact"
                label="紧凑模式"
                tooltip="控制界面元素间距，开启后界面更紧凑"
                placeholder="请选择紧凑模式是否开启"
                allowClear={true}
                options={[
                    {value: '1', label: '开启'},
                    {value: '0', label: '关闭'}
                ]}
            />
            <FormField
                type="select"
                name="layout"
                label="网站布局"
                tooltip="选择网站的整体布局方式，上下布局或左右布局"
                placeholder="请选择网站布局方式"
                allowClear={true}
                options={[
                    {value: 'top', label: '上下布局'},
                    {value: 'left', label: '左右布局'}
                ]}
            />
            <FormField
                type="select"
                name="theme"
                label="菜单主题"
                tooltip="选择菜单的主题模式，亮色或暗色"
                placeholder="请选择菜单主题模式"
                allowClear={true}
                options={[
                    {value: 'light', label: '亮色模式'},
                    {value: 'dark', label: '暗色模式'}
                ]}
            />
            <FormField
                type="select"
                name="watermark"
                label="界面水印"
                tooltip="控制系统界面是否显示用户昵称水印效果"
                placeholder="请选择是否显示界面水印"
                allowClear={true}
                options={[
                    {value: '1', label: '开启'},
                    {value: '0', label: '关闭'}
                ]}
            />
            <FormField
                type="color"
                name="color"
                label="主题颜色"
                tooltip="选择网站的主题颜色，影响整体视觉风格"
                format="rgb"
            />
        </BaseForm>
    );
};

export default UiView;
