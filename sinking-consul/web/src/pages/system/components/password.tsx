import React, {useState} from 'react';
import {App, Form, Spin, Input, Button} from "antd";
import {createStyles} from "antd-style";
import {updatePassword} from "@/service/admin/person";

const useStyles = createStyles(({css}) => {
    return {
        box: css`
            .ant-form-item .ant-form-item-control {
                margin-bottom: 10px !important;
            }
        `
    }
})

const PasswordView: React.FC = () => {
    const {styles: {box}} = useStyles()

    const [isLoading, setIsLoading] = useState(false);
    const {message} = App.useApp()
    const [form] = Form.useForm();

    /**
     * 提交表单
     */
    const onFinish = async (values: any) => {
        setIsLoading(true);
        await updatePassword({
            body: {
                password: values.password
            },
            onSuccess: (r: any) => {
                form.resetFields();
                message?.success(r?.message || "密码修改成功");
            },
            onFail: (r: any) => {
                message?.error(r?.message || "密码修改失败");
            },
            onFinally: () => {
                setIsLoading(false);
            }
        });
    }

    return (
        <Spin spinning={isLoading} size="default">
            <div style={{display: isLoading ? 'none' : 'block'}}>
                <Form form={form} onFinish={onFinish} className={box} layout="vertical">
                    <Form.Item
                        name="password"
                        label="新密码"
                        tooltip="请输入新的登录密码"
                        rules={[
                            {required: true, message: "请输入新密码"},
                            {min: 6, message: "密码至少6位"}
                        ]}
                        style={{maxWidth: '400px', width: '100%'}}
                    >
                        <Input.Password placeholder="请输入新密码"/>
                    </Form.Item>
                    <Form.Item
                        name="confirmPassword"
                        label="确认密码"
                        tooltip="请再次输入密码进行确认"
                        dependencies={['password' as any]}
                        rules={[
                            {required: true, message: "请确认密码"},
                            ({getFieldValue}) => ({
                                validator(_, value) {
                                    if (!value || getFieldValue("password" as any) === value) {
                                        return Promise.resolve();
                                    }
                                    return Promise.reject(new Error('两次输入的密码不一致'));
                                },
                            }),
                        ]}
                        style={{maxWidth: '400px', width: '100%'}}
                    >
                        <Input.Password placeholder="请再次输入密码"/>
                    </Form.Item>
                    <Form.Item>
                        <Button type="primary" htmlType="submit" loading={isLoading}>
                            修改密码
                        </Button>
                    </Form.Item>
                </Form>
            </div>
        </Spin>
    );
};

export default PasswordView;
