import React, {useEffect, useRef, useState} from "react";
import {App, Form, Input} from "antd";
import {updatePassword} from "@/service/admin/person";
import {actions as Actions} from "./common";

export default ({styles}: any): React.ReactNode => {
    const [saving, setSaving] = useState(false);
    const {message} = App.useApp();
    const [form] = Form.useForm();
    const mountedRef = useRef(true);
    const savingRef = useRef(false);

    useEffect(() => {
        mountedRef.current = true;
        return () => {
            mountedRef.current = false;
        };
    }, []);

    const save = async (values: {password: string}) => {
        if (savingRef.current) {
            return;
        }
        savingRef.current = true;
        setSaving(true);
        try {
            const response = await updatePassword({body: {password: values.password}});
            if (!mountedRef.current) {
                return;
            }
            if (response?.code !== 200) {
                message.error(response?.message || "密码修改失败");
                return;
            }
            form.resetFields();
            message.success(response?.message || "密码修改成功");
        } catch {
            if (mountedRef.current) {
                message.error("密码修改失败");
            }
        } finally {
            savingRef.current = false;
            if (mountedRef.current) {
                setSaving(false);
            }
        }
    };

    return (
        <Form form={form} layout="vertical" variant="filled" onFinish={save}>
            <Form.Item name="password" label="新密码" tooltip="请输入新的登录密码" className={styles.formField}
                       rules={[
                           {required: true, message: "请输入新密码"},
                           {min: 6, message: "密码至少6位"},
                       ]}>
                <Input.Password placeholder="请输入新密码" autoComplete="new-password"/>
            </Form.Item>
            <Form.Item name="confirmPassword" label="确认密码" tooltip="请再次输入密码进行确认"
                       className={styles.formField} dependencies={["password"]}
                       rules={[
                           {required: true, message: "请确认密码"},
                           ({getFieldValue}) => ({
                               validator(_, value) {
                                   if (!value || getFieldValue("password") === value) {
                                       return Promise.resolve();
                                   }
                                   return Promise.reject(new Error("两次输入的密码不一致"));
                               },
                           }),
                       ]}>
                <Input.Password placeholder="请再次输入密码" autoComplete="new-password"/>
            </Form.Item>
            <Actions styles={styles} saving={saving} onReset={() => form.resetFields()}/>
        </Form>
    );
};
