import {LockOutlined, UserOutlined} from '@ant-design/icons';
import {App, Button, Checkbox, Col, Form, Input, Row, Tabs, TabsProps} from 'antd';
import React, {useRef, useState} from 'react';
import {Body, Footer} from '@/components';
import {useModel} from "umi";
import {login} from "@/service/user/login";
import {setLoginToken} from "@/utils/auth";
import Captcha, {CaptchaRef} from "@/pages/components/captcha";
import {sendSms} from "@/service/common/sms";
import {createStyles} from "antd-style";
import Settings from "@/../config/defaultSettings";
import {historyPush} from "@/utils/route";
import {NamePath} from "rc-field-form/es/interface";

const useStyles = createStyles(({css, responsive}): any => {
    return {
        container: {
            display: "flex",
            flexDirection: "column",
            height: "100vh",
            overflow: "auto",
            backgroundImage: "url('https://gw.alipayobjects.com/zos/rmsportal/TVYTbAXWheQpRcWDaDMu.svg')",
            backgroundRepeat: "no-repeat",
            backgroundPosition: "center 110px",
            backgroundSize: "100%",
        },
        content: {
            flex: 1,
            padding: "32px 0"
        },
        main: css`
            width: 328px;
            margin: 0 auto;

            ${responsive.md} {
                .ant-tabs .ant-tabs-tab {
                    font-size: 12px;
                }

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
        line: css`
            .ant-tabs-nav::before {
                border: none !important;
            }
        `
    };
});

const Index: React.FC = () => {
    const {
        styles: {
            container, content, top, header, logo, desc, main, line
        }
    } = useStyles();
    const {message, modal} = App.useApp()
    const captcha = useRef<CaptchaRef>({});
    /**
     * 登陆方式
     */
    const [type, setType] = useState<string>('account');
    const [isLoading, setIsLoading] = useState(false);
    const [isSmsSendLoading, setIsSmsSendLoading] = useState(false);
    const [isRead, setIsRead] = useState(true);
    /**
     * 表单
     */
    const [form] = Form.useForm();
    /**
     * 获取当前用户信息
     */
    const web = useModel("web")
    const [sendCodeDisabled, setSendCodeDisabled] = useState(false);
    const getCode = (e: any) => {
        let time = 60;
        const timer = setInterval(() => {
            setSendCodeDisabled(true);
            e.target.innerHTML = `${time}秒后重新获取`;
            time--;
            if (time <= 0) {
                setSendCodeDisabled(false);
                e.target.innerHTML = ' 获取验证码';
                time = 0;
                clearInterval(timer);
            }
        }, 1000);
    };

    const items: TabsProps['items'] = [
        {
            key: 'account',
            label: '密码登录',
        },
        {
            key: 'phone',
            label: '短信登录',
        },
    ];
    return (
        <Body breadCrumb={false}>
            <div className={container}>
                <Captcha ref={captcha}/>
                <div className={content}>
                    <div className={top}>
                        <div className={header}>
                            <img alt='logo' className={logo}
                                 src={web?.info?.logo || ((Settings?.basePath || '/') + 'logo.svg')}/>
                            <span style={{
                                fontSize: "30px",
                                fontWeight: "bolder"
                            }}>{web?.info?.name || Settings?.title}</span>
                        </div>
                        <div className={desc}>
                            {web?.info?.name || Settings?.title}欢迎您的使用
                        </div>
                    </div>
                    <div className={main}>
                        <Tabs items={items} centered className={line} style={{marginBottom: "10px"}}
                              activeKey={type}
                              onChange={(key: string) => {
                                  setType(key)
                              }}/>
                        <Form form={form} size="large" onFinish={async (values: any) => {
                            if (!isRead) {
                                message?.error("请阅读并同意《用户使用条款》协议")
                                return;
                            }
                            values.device = "web";
                            if (type == "account") {
                                if (values.phone == undefined || values.phone.length < 5) {
                                    message?.error("请输入正确的手机号");
                                    return;
                                }
                                if (values.password == undefined || values.password?.length <= 5) {
                                    message?.error("请输入正确的密码");
                                    return;
                                }
                            } else if (type == 'phone') {
                                if (!/^[1][3,4,5,6,7,8,9][0-9]{9}$/.test(values.phone)) {
                                    message?.error("请输入正确的手机号");
                                    return;
                                }
                                if (values.code == undefined || values.code?.length != 6) {
                                    message?.error("请输入正确的手机验证码");
                                    return;
                                }
                            } else {
                                message?.error("不支持该登录方式");
                                return;
                            }
                            setIsLoading(true);
                            if (type == "account") {
                                await login({
                                    body: {
                                        phone: values.phone,
                                        password: values.password,
                                        device: values.device
                                    },
                                    onSuccess: (r: any) => {
                                        setLoginToken(values.device, r?.data);
                                        message?.success(r?.message);
                                        historyPush("user.index");
                                    },
                                    onFail: (r: any) => {
                                        message?.error(r?.message || "请求失败")
                                    },
                                    onFinally: () => {
                                        setIsLoading(false);
                                    }
                                })
                            } else if (type == "phone") {
                                captcha?.current?.Show?.(async (res) => {
                                    await login({
                                        body: {
                                            phone: values.phone,
                                            code: values.code,
                                            device: values.device,
                                            captcha_id: res?.randstr,
                                            captcha_code: res?.ticket
                                        },
                                        onSuccess: (r) => {
                                            setLoginToken(values.device, r?.data);
                                            message?.success(r?.message);
                                            historyPush("user.index")
                                        },
                                        onFail: (r) => {
                                            modal.error({
                                                title: '登录失败',
                                                content: r?.message,
                                                okText: "确认"
                                            })
                                        },
                                        onFinally: () => {
                                            setIsLoading(false)
                                        }
                                    });
                                });
                            }
                        }}>
                            {type === 'account' && (
                                <>
                                    <Form.Item name='phone'>
                                        <Input prefix={<UserOutlined className='site-form-item-icon'/>}
                                               placeholder='请输入手机号' size={'large'}/>
                                    </Form.Item>
                                    <Form.Item name='password'>
                                        <Input type={"password"}
                                               prefix={<LockOutlined className='site-form-item-icon'/>}
                                               size={'large'}
                                               placeholder='请输入账户密码'/>
                                    </Form.Item></>
                            )}
                            {type == 'phone' && <>
                                <Form.Item name='phone'>
                                    <Input prefix={<UserOutlined className='site-form-item-icon'/>}
                                           size={'large'}
                                           placeholder='请输入手机号'/>
                                </Form.Item>
                                <Form.Item style={{marginBottom: "0px"}}>
                                    <Row gutter={4} wrap={false}>
                                        <Col flex='7'>
                                            <Form.Item name='code'>
                                                <Input prefix={<LockOutlined className='site-form-item-icon'/>}
                                                       placeholder='短信验证码'
                                                       size={'large'}/>
                                            </Form.Item>
                                        </Col>
                                        <Col>
                                            <Button size={'large'}
                                                    loading={isSmsSendLoading}
                                                    onClick={(e) => {
                                                        const phone = form.getFieldValue("phone" as NamePath);
                                                        if (!/^[1][3,4,5,6,7,8,9][0-9]{9}$/.test(phone)) {
                                                            message?.error("请输入正确的手机号")
                                                            return;
                                                        }
                                                        captcha?.current?.Show?.(async (res) => {
                                                            setIsSmsSendLoading(true);
                                                            await sendSms({
                                                                body: {
                                                                    token: res?.token,
                                                                    captcha_x: res?.x,
                                                                    captcha_y: res?.y,
                                                                    phone: phone,
                                                                },
                                                                onSuccess: (r) => {
                                                                    message?.success(r?.message)
                                                                    getCode(e)
                                                                },
                                                                onFail: (r) => {
                                                                    message?.error(r?.message)
                                                                },
                                                                onFinally: () => {
                                                                    setIsSmsSendLoading(false)
                                                                }
                                                            });
                                                        })
                                                    }}
                                                    disabled={sendCodeDisabled}>获取验证码</Button>
                                        </Col>
                                    </Row>
                                </Form.Item>
                            </>}

                            <Form.Item hidden={type != "account"} style={{marginTop: "-10px"}}>
                                <a onClick={() => {
                                    message?.warning("使用短信登陆自动创建账户");
                                    setType("phone");
                                }}>注册帐号</a>
                                <a style={{float: 'right',}} onClick={() => {
                                    message?.warning("请使用短信验证登陆后修改密码");
                                    setType("phone");
                                }}>忘记密码</a>
                            </Form.Item>
                            <Form.Item>
                                <Button type='primary' loading={isLoading} htmlType='submit' size={'large'} block>
                                    登 录
                                </Button>
                            </Form.Item>
                            <Row>
                                <Form.Item name="read">
                                    <Checkbox checked={isRead} onChange={(e) => {
                                        setIsRead(e.target.checked);
                                    }}>我已阅读并同意 <a href=''>《用户使用条款》</a>协议</Checkbox>
                                </Form.Item>
                            </Row>
                        </Form>
                    </div>
                </div>
                <Footer/>
            </div>
        </Body>

    );
};

export default Index;
