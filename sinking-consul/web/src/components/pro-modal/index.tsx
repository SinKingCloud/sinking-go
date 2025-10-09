// noinspection TypeScriptValidateTypes
import React, {forwardRef, useImperativeHandle, useState, useMemo, useCallback} from "react";
import {Modal, ModalProps} from "antd";

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

    const isControlled = useMemo(() => (modalProps as any)?.open !== undefined, [(modalProps as any)?.open]);
    const [internalOpen, setInternalOpen] = useState<boolean>(false);
    const open = isControlled ? (modalProps as any)?.open : internalOpen;

    /**
     * 处理确认按钮点击事件
     */
    const handleOk = useCallback(async () => {
        await onOk?.();
    }, [onOk]);

    /**
     * 处理取消按钮点击事件
     */
    const handleCancel = useCallback(() => {
        if (onCancel) {
            onCancel();
        } else if (!isControlled) {
            setInternalOpen(false);
            afterOpenChange?.(false);
        }
    }, [onCancel, isControlled, afterOpenChange]);

    /**
     * 处理模态框关闭后的事件
     */
    const handleAfterClose = useCallback(() => {
        afterClose?.();
    }, [afterClose]);

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
    }), [isControlled, afterOpenChange]);

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
            {...(modalProps || {})}
        >
            {children}
        </Modal>
    );
});

export default ProModal;