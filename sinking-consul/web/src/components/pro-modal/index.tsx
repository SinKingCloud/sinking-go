// noinspection TypeScriptValidateTypes
import React, {forwardRef, useImperativeHandle, useState} from "react";
import {Modal, ModalProps} from "antd";
import {useTheme} from "@/components/theme";

/**
 * ProModal的属性类型
 */
export type ProModalProps = {
    modalProps?: ModalProps & any & undefined;// modal属性
    title?: any; // 标题
    onOk?: () => void | any; // 确认按钮回调
    onCancel?: () => void | any; // 取消按钮回调
    afterClose?: () => void | any; // 关闭后回调
    okText?: string; // 确认按钮文本
    okType?: "primary" | "default" | "dashed" | "link" | "text"; // 确认按钮类型
    width?: number | string; // 模态框宽度
    afterOpenChange?: (open: boolean) => void | any; // 打开状态变化回调（可选）
    children?: React.ReactNode; // 子组件内容
};

/**
 * ProModal组件
 */
export interface ProModalRef {
    show: () => void; // 显示模态框
    hide: () => void; // 隐藏模态框
}

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

    const theme = useTheme();
    const [internalOpen, setInternalOpen] = useState<boolean>(false);

    // 判断是否为受控模式
    const isControlled = (modalProps as any)?.open !== undefined;
    const open = isControlled ? (modalProps as any)?.open : internalOpen;

    // 根据紧凑模式调整宽度（紧凑模式下减少 20% 宽度）
    const adjustedWidth = theme?.isCompactTheme?.()
        ? (typeof width === 'number' ? Math.floor(width * 0.87) : width)
        : width;

    /**
     * 处理确认按钮点击事件
     */
    const handleOk = async () => {
        await onOk?.();
    };

    /**
     * 处理取消按钮点击事件
     */
    const handleCancel = () => {
        if (onCancel) {
            onCancel();
        } else if (!isControlled) {
            setInternalOpen(false);
            afterOpenChange?.(false);
        }
    };

    /**
     * 处理模态框关闭后的事件
     */
    const handleAfterClose = () => {
        afterClose?.();
    };

    /**
     * 暴露给父组件的方法
     */
    useImperativeHandle(ref, () => ({
        show: () => {
            if (!isControlled) {
                setInternalOpen(true);
            }
            afterOpenChange?.(true);
        },
        hide: () => {
            if (!isControlled) {
                setInternalOpen(false);
            }
            afterOpenChange?.(false);
        }
    }));

    return (
        <Modal
            title={title}
            open={open}
            onOk={handleOk}
            onCancel={handleCancel}
            afterClose={handleAfterClose}
            okText={okText}
            okType={okType}
            width={adjustedWidth}
            maskClosable={false}
            {...modalProps}
        >
            {children}
        </Modal>
    );
});

export default ProModal;