import {useCallback, useEffect, useRef, useState} from "react";
import {App, Form} from "antd";
import {useModel} from "umi";
import {getConfig, setConfig} from "@/service/admin/system";

type ConfigValues = Record<string, any>;

interface UseSystemConfigOptions {
    group: string;
    defaults: ConfigValues;
    normalize?: (values: ConfigValues) => ConfigValues;
}

const serialize = (value: any) => typeof value === "boolean" ? (value ? "1" : "0") : String(value ?? "");
const identity = (values: ConfigValues) => values;

export const useSystemConfig = ({group, defaults, normalize = identity}: UseSystemConfigOptions) => {
    const {message} = App.useApp();
    const web = useModel("web");
    const [form] = Form.useForm();
    const [dataLoading, setDataLoading] = useState(true);
    const [submitLoading, setSubmitLoading] = useState(false);
    const [error, setError] = useState("");
    const mountedRef = useRef(true);
    const savingRef = useRef(false);
    const requestRef = useRef(0);
    const initialRef = useRef<ConfigValues>({});
    const defaultsRef = useRef(defaults);
    const normalizeRef = useRef(normalize);
    defaultsRef.current = defaults;
    normalizeRef.current = normalize;

    const loadConfig = useCallback(async () => {
        const requestId = ++requestRef.current;
        setDataLoading(true);
        setError("");
        try {
            const response = await getConfig({body: {action: "get", group}});
            if (!mountedRef.current || requestId !== requestRef.current) {
                return;
            }
            if (response?.code !== 200) {
                setError(response?.message || "配置加载失败");
                return;
            }
            const source = response?.data || {};
            const values = Object.keys(defaultsRef.current).reduce<ConfigValues>((result, key) => {
                result[key] = source[key] === undefined ? defaultsRef.current[key] : source[key];
                return result;
            }, {});
            const nextValues = normalizeRef.current(values);
            initialRef.current = nextValues;
            form.setFieldsValue(nextValues);
        } catch {
            if (mountedRef.current && requestId === requestRef.current) {
                setError("配置加载失败");
            }
        } finally {
            if (mountedRef.current && requestId === requestRef.current) {
                setDataLoading(false);
            }
        }
    }, [form, group]);

    useEffect(() => {
        mountedRef.current = true;
        loadConfig();
        return () => {
            mountedRef.current = false;
            requestRef.current += 1;
        };
    }, [loadConfig]);

    const saveConfig = useCallback(async (values: ConfigValues) => {
        if (savingRef.current) {
            return;
        }
        savingRef.current = true;
        setSubmitLoading(true);
        try {
            const nextValues = normalizeRef.current(values);
            const response = await setConfig({
                body: {
                    action: "set",
                    group,
                    configs: Object.entries(nextValues).map(([key, value]) => ({key, value: serialize(value)})),
                },
            });
            if (response?.code === 200) {
                web?.refreshInfo();
            }
            if (!mountedRef.current) {
                return;
            }
            if (response?.code !== 200) {
                message.error(response?.message || "配置保存失败");
                return;
            }
            initialRef.current = nextValues;
            form.setFieldsValue(nextValues);
            message.success(response?.message || "配置保存成功");
        } catch {
            if (mountedRef.current) {
                message.error("配置保存失败");
            }
        } finally {
            savingRef.current = false;
            if (mountedRef.current) {
                setSubmitLoading(false);
            }
        }
    }, [form, group, message, web]);

    const resetForm = useCallback(() => {
        form.resetFields();
        form.setFieldsValue(initialRef.current);
    }, [form]);

    return {form, dataLoading, submitLoading, error, loadConfig, saveConfig, resetForm};
};
