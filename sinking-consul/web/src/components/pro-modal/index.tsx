// noinspection TypeScriptValidateTypes

import React, {forwardRef, useImperativeHandle, useState} from "react";
import {createStyles} from "antd-style";
import {Modal, ModalProps} from "antd";

/**
 * ProModal的属性类型
 */
export type ProModalProps = {
    modalProps?: ModalProps & any & undefined;// modal属性
    title?: any; // 标题
    onOk?: (values: any) => void | any; // 确认按钮回调
    onCancel?: () => void | any; // 取消按钮回调
    afterClose?: () => void | any; // 关闭后回调
    okText?: string; // 确认按钮文本
    okType?: "primary" | "default" | "dashed" | "link" | "text"; // 确认按钮类型
    width?: number | string; // 模态框宽度
    afterOpenChange?: (open: boolean) => void | any; // 打开状态变化回调
    children?: React.ReactNode; // 子组件内容
};

/**
 * ProModal组件
 */
export interface ProModalRef {
    show?: () => void | any; // 显示模态框
    hide?: () => void | any; // 隐藏模态框
}

/**
 * 样式配置
 */
const useStyles: any = createStyles(({token, isDarkMode}: any): any => {
    return {};
});

const ProModal = forwardRef<ProModalRef, ProModalProps>((props, ref): any => {
    const {
        modalProps = {},
        title,
        onOk,
        onCancel,
        afterClose,
        okText,
        okType,
        width = 520,
        afterOpenChange,
        children,
    } = props;

    const {styles} = useStyles();
    const [open, setOpen] = useState<boolean>(false);

    // 暴露给父组件的方法
    useImperativeHandle(ref, () => ({
        show: () => {
            setOpen(true);
            afterOpenChange?.(true);
        },
        hide: () => {
            setOpen(false);
            afterOpenChange?.(false);
        }
    }));

    // 处理确认按钮点击
    const handleOk = async () => {
        try {
            if (onOk) {
                const result = await onOk({});
                // 如果 onOk 返回 false，则不关闭模态框
                if (result !== false) {
                    setOpen(false);
                    afterOpenChange?.(false);
                }
            } else {
                setOpen(false);
                afterOpenChange?.(false);
            }
        } catch (error) {
            console.error('Modal onOk error:', error);
        }
    };

    // 处理取消按钮点击
    const handleCancel = () => {
        if (onCancel) {
            const result = onCancel();
            // 如果 onCancel 返回 false，则不关闭模态框
            if (result !== false) {
                setOpen(false);
                afterOpenChange?.(false);
            }
        } else {
            setOpen(false);
            afterOpenChange?.(false);
        }
    };

    // 处理模态框关闭后的回调
    const handleAfterClose = () => {
        afterClose?.();
    };

    return (
        <Modal
            title={title}
            open={open}
            onOk={handleOk}
            onCancel={handleCancel}
            afterClose={handleAfterClose}
            okText={okText}
            okType={okType}
            width={width}
            maskClosable={false}
            {...modalProps}
        >
            {children}
        </Modal>
    );
});

export default ProModal;
