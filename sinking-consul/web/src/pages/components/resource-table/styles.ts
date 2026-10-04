import {createStyles} from "antd-style";

const useStyles = createStyles(({css, token}) => ({
    table: css`
        min-width: 0;

        .resource-copy {
            display: block;
            min-width: 0;
            max-width: 100%;
            margin: 0;
            color: inherit;
            font-size: inherit;
            line-height: 20px;
            overflow-wrap: anywhere;
        }

        .resource-copy .ant-typography-copy {
            margin-inline-start: 5px;
            color: ${token.colorTextQuaternary};
            font-size: 10px;
            line-height: 1;
            vertical-align: 0;
            opacity: 0;
            transition: color 0.16s ease, opacity 0.16s ease;
        }

        .resource-copy:hover .ant-typography-copy,
        .resource-copy:focus-within .ant-typography-copy {
            color: ${token.colorPrimary};
            opacity: 1;
        }

        .resource-state {
            display: inline-flex;
            align-items: center;
            gap: 5px;
            line-height: 20px;
            white-space: nowrap;
            color: ${token.colorTextSecondary};
        }

        .resource-state i {
            width: 6px;
            height: 6px;
            flex: none;
            border-radius: 50%;
            background: currentColor;
        }

        .resource-state.is-success { color: ${token.colorSuccess}; }
        .resource-state.is-warning { color: ${token.colorWarning}; }
        .resource-state.is-error { color: ${token.colorError}; }

        .resource-time {
            display: block;
            overflow: hidden;
            line-height: 20px;
            text-overflow: ellipsis;
            white-space: nowrap;
        }

        @media (hover: none) {
            .resource-copy .ant-typography-copy { opacity: 0.55; }
        }
    `,
    form: css`
        padding-top: 8px;
        .ant-form-item:last-child { margin-bottom: 8px; }
    `,
}));

export default useStyles;
