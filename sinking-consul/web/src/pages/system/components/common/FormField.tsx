// noinspection TypeScriptValidateTypes

import React from 'react';
import {Form, Input, Select, ColorPicker} from 'antd';
import {FormItemProps} from 'antd/lib/form';

interface BaseFieldProps extends Omit<FormItemProps, 'children'> {
    name: string;
    label: string;
    tooltip?: string;
    placeholder?: string;
    required?: boolean;
    minLength?: number;
    maxWidth?: string | number;
    allowClear?: boolean;
    options?: any;
    rows?: any;
    format?: any;
}

interface InputFieldProps extends BaseFieldProps {
    type: 'input' | 'password' | 'textarea';
    rows?: number;
}

interface SelectFieldProps extends BaseFieldProps {
    type: 'select';
    options: Array<{ value: string | number; label: string }>;
    allowClear?: boolean;
}

interface ColorFieldProps extends BaseFieldProps {
    type: 'color';
    format?: 'hex' | 'rgb';
}

type FormFieldProps = InputFieldProps | SelectFieldProps | ColorFieldProps;

const FormField: React.FC<FormFieldProps> = (props) => {
    const {
        type,
        name,
        label,
        tooltip,
        placeholder,
        required = false,
        minLength,
        maxWidth = '400px',
        ...restProps
    } = props;

    const baseRules = [];
    if (required) {
        baseRules.push({required: true, message: `请输入${label}`});
    }
    if (minLength) {
        baseRules.push({min: minLength, message: `${label}至少${minLength}位`});
    }

    const itemProps = {
        name,
        label,
        tooltip,
        rules: baseRules,
        style: {maxWidth, width: '100%'},
        ...restProps
    };

    const renderField = () => {
        switch (type) {
            case 'input':
                return <Input placeholder={placeholder}/>;

            case 'password':
                return <Input.Password placeholder={placeholder}/>;

            case 'textarea':
                const textareaProps = props as InputFieldProps;
                return <Input.TextArea placeholder={placeholder} rows={textareaProps.rows || 4}/>;

            case 'select':
                const selectProps = props as SelectFieldProps;
                return (
                    <Select
                        placeholder={placeholder}
                        allowClear={selectProps.allowClear}
                        options={selectProps.options}
                    />
                );

            case 'color':
                const colorProps = props as ColorFieldProps;
                return (
                    <ColorPicker
                        format={colorProps.format || 'rgb'}
                        defaultFormat={colorProps.format || 'rgb'}
                    />
                );

            default:
                return <Input placeholder={placeholder}/>;
        }
    };

    if (type === 'color') {
        return (
            <Form.Item
                {...itemProps}
                getValueFromEvent={(color) => color?.toRgbString()}
            >
                {renderField()}
            </Form.Item>
        );
    }

    return (
        <Form.Item {...itemProps}>
            {renderField()}
        </Form.Item>
    );
};

export default FormField;
