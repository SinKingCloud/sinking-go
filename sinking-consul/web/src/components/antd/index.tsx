import { App } from 'antd';
import type { MessageInstance } from 'antd/es/message/interface';
import type { ModalStaticFunctions } from 'antd/es/modal/confirm';
import type { NotificationInstance } from 'antd/es/notification/interface';

let Message: MessageInstance;
let Notification: NotificationInstance;
let Modal: Omit<ModalStaticFunctions, 'warn'>;

export default () => {
    const staticFunction = App.useApp();
    Message = staticFunction.message;
    Modal = staticFunction.modal;
    Notification = staticFunction.notification;
    return null;
};

export { Message, Notification, Modal };