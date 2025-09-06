import React from 'react';
import { Form, Button, Spin } from 'antd';
import { FormInstance } from 'antd/lib/form';

interface BaseFormProps {
  form: FormInstance;
  loading?: boolean;
  submitLoading?: boolean;
  onFinish: (values: any) => void;
  onReset: () => void;
  children: React.ReactNode;
  hideLoading?: boolean;
}

const BaseForm: React.FC<BaseFormProps> = ({
  form,
  loading = false,
  submitLoading = false,
  onFinish,
  onReset,
  children,
  hideLoading = false
}) => {
  const content = (
    <div style={{ display: loading && !hideLoading ? 'none' : 'block' }}>
      <Form form={form} onFinish={onFinish} layout="vertical">
        {children}
        <Form.Item>
          <Button onClick={onReset} style={{ marginRight: 8 }}>
            重置
          </Button>
          <Button type="primary" htmlType="submit" loading={submitLoading}>
            提交
          </Button>
        </Form.Item>
      </Form>
    </div>
  );

  if (hideLoading) {
    return content;
  }

  return (
    <Spin spinning={loading} size="default">
      {content}
    </Spin>
  );
};

export default BaseForm;
