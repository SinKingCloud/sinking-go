import {createStyles} from 'antd-style';

const useStyles = createStyles(({css, token}) => ({
    table: css`
        min-width: 0;

        .config-copy {
            display: block;
            min-width: 0;
            max-width: 100%;
            margin: 0;
            color: inherit;
            font: inherit;
            line-height: 20px;
            overflow-wrap: anywhere;
        }

        .config-copy .ant-typography-copy {
            margin-inline-start: 5px;
            color: ${token.colorTextQuaternary};
            font-size: 10px;
            opacity: 0;
            transition: color .16s ease, opacity .16s ease;
        }

        .config-copy:hover .ant-typography-copy,
        .config-copy .ant-typography-copy:focus-visible {
            color: ${token.colorPrimary};
            opacity: 1;
        }

        .config-name {
            max-width: 100%;
            padding: 0;
            border: 0;
            background: transparent;
            color: inherit;
            cursor: pointer;
            font: inherit;
            text-align: left;
            overflow-wrap: anywhere;
        }

        .config-name:hover,
        .config-name:focus-visible {
            color: ${token.colorPrimary};
        }

        .config-hash {
            display: flex;
            align-items: center;
            color: ${token.colorTextSecondary};
            font-family: ${token.fontFamilyCode};
            font-size: 11px;
        }

        .config-hash-value {
            display: block;
            flex: 1;
            min-width: 0;
            overflow: hidden;
            text-overflow: ellipsis;
            white-space: nowrap;
        }

        .config-hash .ant-typography-copy { flex: none; }

        .config-state {
            display: inline-flex;
            align-items: center;
            gap: 4px;
            line-height: 20px;
            white-space: nowrap;
        }

        .config-state i {
            width: 6px;
            height: 6px;
            border-radius: 50%;
            background: currentColor;
        }

        .config-state.enabled { color: ${token.colorSuccess}; }
        .config-state.disabled { color: ${token.colorWarning}; }

        .config-text {
            display: block;
            min-width: 0;
            overflow: hidden;
            line-height: 20px;
            text-overflow: ellipsis;
            white-space: nowrap;
        }

        @media (hover: none) {
            .config-copy .ant-typography-copy { opacity: .55; }
        }
    `,
    modal: css`
        .ant-modal {
            max-width: calc(100vw - 24px);
            outline: none;
        }
    `,
    editor: css`
        overflow: hidden;
        border: 1px solid ${token.colorBorderSecondary};
        border-radius: ${token.borderRadius}px;

        .ace_editor { border-radius: ${token.borderRadius}px !important; }
    `,
    loading: css`
        height: 380px;
        display: flex;
        align-items: center;
        justify-content: center;
    `,
}));

export default useStyles;
