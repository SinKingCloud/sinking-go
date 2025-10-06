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
            <FormField
                type="select"
                name="radius"
                label="主题圆角"
                tooltip="调整界面元素的圆角大小，数值越大圆角越明显"
                placeholder="请选择主题圆角大小"
                options={[
                    {value: '0', label: '0px (无圆角)'},
                    {value: '1', label: '1px'},
                    {value: '2', label: '2px'},
                    {value: '3', label: '3px'},
                    {value: '4', label: '4px'},
                    {value: '5', label: '5px'},
                    {value: '6', label: '6px'},
                    {value: '7', label: '7px'},
                    {value: '8', label: '8px'},
                    {value: '9', label: '9px'},
                    {value: '10', label: '10px'},
                    {value: '11', label: '11px'},
                    {value: '12', label: '12px'},
                    {value: '13', label: '13px'},
                    {value: '14', label: '14px'},
                    {value: '15', label: '15px'}
                ]}
            />
        </BaseForm>
    );
};

export default UiView;
