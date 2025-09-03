import {LockOutlined, UserOutlined} from '@ant-design/icons';
import {App, Button, Form, Input} from 'antd';
import React, {useRef, useState} from 'react';
import {Body} from '@/components';
import {useModel} from "umi";
import {login} from "@/service/auth/login";
import {setLoginToken} from "@/utils/auth";
import Captcha, {CaptchaRef} from "@/pages/components/captcha";
import {createStyles} from "antd-style";
import Settings from "@/../config/defaultSettings";
import {historyPush} from "@/utils/route";

const useStyles = createStyles(({css, responsive}): any => {
    return {
        container: {
            display: "flex",
            flexDirection: "column",
            height: "100vh",
            backgroundImage: "url('https://gw.alipayobjects.com/zos/rmsportal/TVYTbAXWheQpRcWDaDMu.svg')",
            backgroundRepeat: "no-repeat",
            backgroundPosition: "center 110px",
            backgroundSize: "100%",
        },
        content: {
            flex: 1,
            padding: "120px 0 32px 0"
        },
        main: css`
            width: 328px;
            margin: 0 auto;

            ${responsive.md} {
                width: 95%;
                max-width: 300px;
            }
        `,
        top: {
            textAlign: "center"
        },
        header: {
            height: 40,
            lineHeight: "40px",
            a: {
                textDecoration: "none"
            },
            span: {
                fontSize: "30px",
                fontWeight: "bolder"
            }
        },
        logo: {
            height: 40,
            marginRight: 16,
            verticalAlign: "top",
        },
        desc: {
            marginTop: 12,
            marginBottom: 20,
            color: "@text-color-secondary",
            fontSize: 14,
        },
    };
});

const Login: React.FC = () => {
    const {
        styles: {
            container, content, top, header, logo, desc, main
        }
    } = useStyles();
    const {message} = App.useApp();
    const captcha = useRef<CaptchaRef>({});
    const [isLoading, setIsLoading] = useState(false);
    /**
     * 表单
     */
    const [form] = Form.useForm();
    /**
     * 获取当前用户信息
     */
    const web = useModel("web");

    return (
        <Body breadCrumb={false}>
            <div className={container}>
                <Captcha ref={captcha}/>
                <div className={content}>
                    <div className={top}>
                        <div className={header}>
                            <img alt='logo' className={logo}
                                 src={web?.info?.logo || ((Settings?.basePath || '/') + 'logo.svg')}/>
                            <span>{web?.info?.name || Settings?.title}</span>
                        </div>
                        <div className={desc}>
                            {web?.info?.name || Settings?.title}欢迎您的使用
                        </div>
                    </div>
                    <div className={main}>
                        <Form form={form} size="large" onFinish={async (values: any) => {
                            if (values.account == undefined || values.account.length < 3) {
                                message?.error("请输入正确的账户");
                                return;
                            }
                            if (values.password == undefined || values.password?.length < 4) {
                                message?.error("请输入正确的密码");
                                return;
                            }
                            setIsLoading(true);
                            captcha?.current?.Show?.(
                                async (res) => {
                                    await login({
                                        body: {
                                            account: values.account,
                                            password: values.password,
                                            token: res.token,
                                            captcha_x: res.x,
                                            captcha_y: res.y
                                        },
                                        onSuccess: (r: any) => {
                                            setLoginToken(r?.data);
                                            message?.success(r?.message);
                                            historyPush("/");
                                        },
                                        onFail: (r: any) => {
                                            message?.error(r?.message || "登录失败")
                                        },
                                        onFinally: () => {
                                            setIsLoading(false);
                                        }
                                    });
                                },
                                () => {
                                    setIsLoading(false);
                                }
                            );
                        }}>
                            <Form.Item name='account'>
                                <Input prefix={<UserOutlined className='site-form-item-icon'/>}
                                       placeholder='请输入账户' size={'large'}/>
                            </Form.Item>
                            <Form.Item name='password'>
                                <Input type={"password"}
                                       prefix={<LockOutlined className='site-form-item-icon'/>}
                                       size={'large'}
                                       placeholder='请输入账户密码'/>
                            </Form.Item>
                            <Form.Item>
                                <Button type='primary' loading={isLoading} htmlType='submit' size={'large'} block>
                                    登 录
                                </Button>
                            </Form.Item>
                        </Form>
                    </div>
                </div>
            </div>
        </Body>
    );
};

export default Login;
