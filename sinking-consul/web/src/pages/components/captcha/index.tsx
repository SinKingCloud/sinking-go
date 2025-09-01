import {forwardRef, useState, useImperativeHandle, useRef} from "react";
import {Modal, Spin, App} from "antd";
import GoCaptcha from "go-captcha-react";
import {getCaptcha} from "@/service/common/captcha";
import {getRandStr} from "@/utils/string";
import {createStyles} from "antd-style";

/**
 * 验证码组件
 */
export interface CaptchaRef {
    Show?: (onSuccess: (res: any) => void) => void;
}

/**
 * 样式配置
 */
const useStyles: any = createStyles(({isDarkMode}): any => {
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
            } : {},
        },
    };
});

const Captcha = forwardRef<CaptchaRef>((_, ref): any => {
    const {styles} = useStyles();

    const {message} = App.useApp();

    let slideRef = useRef<any>(null);
    const [visible, setVisible] = useState(false);  //弹窗状态
    const [loading, setLoading] = useState(true);//加载状态
    const [data, setData] = useState<any>(null); //验证码数据

    /**
     * 相关数据
     */
    let token = useRef<string>("");
    let success = useRef<(res: any) => void>(null);


    /**
     * 关闭验证码
     */
    const close = () => {
        setVisible(false);
        setData(null);
    };

    /**
     * 刷新验证码
     */
    const refresh = async () => {
        setLoading(true);
        setVisible(true);
        token.current = getRandStr(16);
        await getCaptcha({
            body: {token: token?.current},
            onSuccess: (res) => {
                const d = res?.data;
                setData({
                    key: d?.key,
                    image: d?.image_base64,
                    thumb: d?.tile_base64,
                    thumbWidth: d?.tile_width,
                    thumbHeight: d?.tile_height,
                    tileX: d?.tile_x,
                    tileY: d?.tile_y,
                    width: d?.width,
                    height: d?.height,
                });
            },
            onFail: () => {
                message.error("验证码加载失败");
                close()
            },
            onFinally: () => {
                setLoading(false);
            }
        });
    }

    /**
     * 显示验证码
     * @param onSuccess 回调函数
     */
    const show = (onSuccess: (res: any) => void = undefined) => {
        if (onSuccess) {
            success.current = onSuccess;
        }
        refresh().finally(() => {
            setVisible(true);
        });
    }

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
                        image: data?.image,
                        thumb: data?.thumb,
                        thumbX: 0,
                        thumbY: data?.tileY,
                        thumbWidth: data?.thumbWidth,
                        thumbHeight: data?.thumbHeight,
                    }}
                    config={{
                        width: data?.width,
                        height: data?.height,
                        scope: true,
                    }}
                    events={{
                        confirm: (point) => {
                            close()
                            if (success) {
                                success?.current?.({token: token?.current, x: point?.x || 0, y: point?.y || 0});
                            }
                        },
                        refresh: async () => {
                            await refresh()
                        },
                        close: close,
                    }}
                />
            </Spin>
        </Modal>
    );
});

export default Captcha;
