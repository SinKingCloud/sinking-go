import {forwardRef, useState, useImperativeHandle, useRef} from "react";
import {Modal, Spin, App} from "antd";
import GoCaptcha from "go-captcha-react";
import {getCaptcha} from "@/service/common/captcha";
import {getRandStr} from "@/utils/string";
import {createStyles} from "antd-style";

/**
 * 验证码组件接口
 */
export interface CaptchaRef {
    Show?: (onSuccess?: (res: any) => void, onClose?: () => void) => void;
}

/**
 * 验证码响应数据接口
 */
interface CaptchaResponse {
    token: string;
    x: number;
    y: number;
}

/**
 * 验证码数据接口
 */
interface CaptchaData {
    key: string;
    image: string;
    thumb: string;
    thumbWidth: number;
    thumbHeight: number;
    tileX: number;
    tileY: number;
    width: number;
    height: number;
}

/**
 * 样式配置
 */
const useStyles: any = createStyles(({isDarkMode, token}) => {
    return {
        modal: {
            ".ant-modal": {
                width: "326px !important",
            },
            ".ant-modal-content": {
                padding: "0 !important",
                backgroundColor: "transparent !important",
                boxShadow: "none !important",
            },
            ".ant-modal-body> div": isDarkMode ? {
                "--go-captcha-theme-text-color": "#8f94a7 !important",
                "--go-captcha-theme-bg-color": "#18181a !important",
                "--go-captcha-theme-btn-color": "#ffffff !important",
                "--go-captcha-theme-btn-bg-color": "#4e87ff !important",
                "--go-captcha-theme-btn-border-color": "#4e87ff !important",
                "--go-captcha-theme-active-color": "#3e7cff !important",
                "--go-captcha-theme-border-color": "#3c3f44 !important",
                "--go-captcha-theme-icon-color": "#696d7b !important",
                "--go-captcha-theme-drag-bar-color": "#3c3f44 !important",
                "--go-captcha-theme-drag-bg-color": "#3e7cff !important",
                "--go-captcha-theme-drag-icon-color": "#ffffff !important",
                "--go-captcha-theme-round-color": "#3c3f44 !important",
                "--go-captcha-theme-loading-icon-color": "#3e7cff !important",
                "--go-captcha-theme-body-bg-color": "#34383e !important",
                "--go-captcha-theme-dot-color": "#cedffe !important",
                "--go-captcha-theme-dot-bg-color": "#3e7cff !important",
                "--go-captcha-theme-dot-border-color": "#f7f9fb !important",
            } : {
                "--go-captcha-theme-btn-color": token?.colorPrimary + " !important",
                "--go-captcha-theme-btn-bg-color": token?.colorPrimary + " !important",
                "--go-captcha-theme-drag-bg-color": token?.colorPrimary + " !important",
            },
        },
    };
});

const Captcha = forwardRef<CaptchaRef>((_, ref): any => {
    const {styles} = useStyles();
    const {message} = App.useApp();

    const slideRef = useRef<any>(null);
    const [visible, setVisible] = useState(false);
    const [loading, setLoading] = useState(true);
    const [data, setData] = useState<CaptchaData | null>(null);

    /**
     * 回调函数引用
     */
    const token = useRef<string>("");
    const successCallback = useRef<((res: CaptchaResponse) => void) | null>(null);
    const closeCallback = useRef<(() => void) | null>(null);

    /**
     * 关闭验证码
     */
    const close = () => {
        setVisible(false);
        setData(null);
        setLoading(true);
        if (closeCallback.current) {
            closeCallback.current();
        }
        successCallback.current = null;
        closeCallback.current = null;
    };

    /**
     * 刷新验证码
     */
    const refresh = async () => {
        try {
            setLoading(true);
            token.current = getRandStr(16);
            await getCaptcha({
                body: {token: token.current},
                onSuccess: (res) => {
                    const d = res?.data;
                    if (d) {
                        setData({
                            key: d.key,
                            image: d.image_base64,
                            thumb: d.tile_base64,
                            thumbWidth: d.tile_width,
                            thumbHeight: d.tile_height,
                            tileX: d.tile_x,
                            tileY: d.tile_y,
                            width: d.width,
                            height: d.height,
                        });
                    }
                },
                onFail: (error) => {
                    message.error(error?.message || "验证码加载失败");
                    close();
                },
                onFinally: () => {
                    setLoading(false);
                }
            });
        } catch (error) {
            message.error("验证码加载失败");
            setLoading(false);
            close();
        }
    };

    /**
     * 显示验证码
     * @param onSuccess 成功回调函数
     * @param onClose 关闭回调函数
     */
    const show = async (onSuccess?: (res: CaptchaResponse) => void, onClose?: () => void) => {
        // 设置回调函数
        successCallback.current = onSuccess || null;
        closeCallback.current = onClose || null;
        // 显示模态框并刷新验证码
        setVisible(true);
        await refresh();
    };

    /**
     * 方法挂载
     */
    useImperativeHandle(ref, () => ({
        Show: show,
    }));

    return (
        <Modal
            open={visible}
            destroyOnHidden={true}
            onCancel={close}
            footer={null}
            closable={false}
            maskClosable={false}
            keyboard={false}
            rootClassName={styles?.modal}
        >
            <Spin spinning={loading} size={"large"}>
                <GoCaptcha.Slide
                    ref={slideRef}
                    data={{
                        image: data?.image || "",
                        thumb: data?.thumb || "",
                        thumbX: 0,
                        thumbY: data?.tileY || 0,
                        thumbWidth: data?.thumbWidth || 0,
                        thumbHeight: data?.thumbHeight || 0,
                    }}
                    config={{
                        width: data?.width,
                        height: data?.height,
                        scope: true,
                    }}
                    events={{
                        confirm: (point) => {
                            const result: CaptchaResponse = {
                                token: token.current,
                                x: point?.x || 0,
                                y: point?.y || 0
                            };
                            if (successCallback.current) {
                                successCallback.current(result);
                            }
                            close();
                        },
                        refresh: refresh,
                        close: close,
                    }}
                />
            </Spin>
        </Modal>
    );
});

export default Captcha;
