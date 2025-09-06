import React, {useState} from 'react';
import {App, Form, Input, Button} from "antd";
import {updatePassword} from "@/service/admin/person";

const PasswordView: React.FC = () => {
    const [submitLoading, setSubmitLoading] = useState(false);
    const {message} = App.useApp();
    const [form] = Form.useForm();

    /**
     * 提交表单
     */
    const onFinish = async (values: any) => {
        setSubmitLoading(true);
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
                setSubmitLoading(false);
            }
        });
    }

    /**
     * 重置表单
     */
    const onReset = () => {
        form.resetFields();
    }

    return (
        <div>
            <Form form={form} onFinish={onFinish} layout="vertical">
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
                    <Button onClick={onReset} style={{marginRight: 8}}>
                        重置
                    </Button>
                    <Button type="primary" htmlType="submit" loading={submitLoading}>
                        提交
                    </Button>
                </Form.Item>
            </Form>
        </div>
    );
};

export default PasswordView;
