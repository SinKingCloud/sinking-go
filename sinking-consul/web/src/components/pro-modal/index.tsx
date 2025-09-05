// noinspection TypeScriptValidateTypes
import React, {forwardRef, useImperativeHandle, useState, useMemo} from "react";
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
        modalProps = {} as ModalProps & any & undefined,
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

    const isControlled = useMemo(() => modalProps.open !== undefined, [modalProps.open]);// 是否受控
    const [internalOpen, setInternalOpen] = useState<boolean>(false);// 内部状态管理模态框的打开与关闭
    const open = isControlled ? modalProps.open : internalOpen;// 是否打开模态框

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
        } else {
            if (!isControlled) {
                setInternalOpen(false);
                afterOpenChange?.(false);
            }
        }
    };

    /**
     * 处理模态框关闭后的事件
     */
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

ProModal.displayName = "ProModal";

export default ProModal;