import { useState, useEffect } from 'react';
import { Form, App } from 'antd';
import { getConfig, setConfig } from "@/service/admin/system";
import { useModel } from "umi";

interface UseSystemConfigOptions {
  group: string;
  autoLoad?: boolean;
  onSuccess?: (data: any) => void;
  onError?: (error: any) => void;
}

export const useSystemConfig = (options: UseSystemConfigOptions) => {
  const { group, autoLoad = true, onSuccess, onError } = options;
  const [dataLoading, setDataLoading] = useState(false);
  const [submitLoading, setSubmitLoading] = useState(false);
  const { message } = App.useApp();
  const [form] = Form.useForm();
  const web = useModel("web");

  /**
   * 加载配置数据
   */
  const loadConfig = async () => {
    setDataLoading(true);
    try {
      await getConfig({
        body: {
          action: "get",
          group
        },
        onSuccess: (r: any) => {
          const data = r?.data || {};
          // 处理空值，确保空字符串被转换为undefined以显示placeholder
          const processedData = Object.keys(data).reduce((acc, key) => {
            const value = data[key];
            acc[key] = value === '' || value === null ? undefined : value;
            return acc;
          }, {} as any);
          form.setFieldsValue(processedData);
          onSuccess?.(processedData);
        },
        onFail: (r: any) => {
          const errorMsg = r?.message || "加载配置失败";
          message?.error(errorMsg);
          onError?.(r);
        },
        onFinally: () => {
          setDataLoading(false);
        }
      });
    } catch (error) {
      setDataLoading(false);
      onError?.(error);
    }
  };

  /**
   * 保存配置
   */
  const saveConfig = async (values: any) => {
    setSubmitLoading(true);
    const configs = Object.entries(values).map(([key, value]) => ({ key, value }));
    
    try {
      await setConfig({
        body: {
          action: "set",
          group,
          configs
        },
        onSuccess: (r: any) => {
          message?.success(r?.message || "配置保存成功");
          web?.refreshInfo();
        },
        onFail: (r: any) => {
          message?.error(r?.message || "配置保存失败");
        },
        onFinally: () => {
          setSubmitLoading(false);
        }
      });
    } catch (error) {
      setSubmitLoading(false);
      message?.error("配置保存失败");
    }
  };

  /**
   * 重置表单
   */
  const resetForm = () => {
    form.resetFields();
  };

  // 自动加载数据
  useEffect(() => {
    if (autoLoad) {
      loadConfig();
    }
  }, [group, autoLoad]);

  return {
    form,
    dataLoading,
    submitLoading,
    loadConfig,
    saveConfig,
    resetForm
  };
};

