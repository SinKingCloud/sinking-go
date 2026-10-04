import {createStyles} from "antd-style";

const useStyles = createStyles(({css, token}: any) => ({
    logTable: css`
        min-width: 0;

        .log-copy {
            display: block;
            min-width: 0;
            max-width: 100%;
            margin: 0;
            color: inherit;
            font: inherit;
            line-height: 20px;
            overflow-wrap: anywhere;
        }

        .log-copy .ant-typography-copy {
            margin-inline-start: 5px;
            color: ${token.colorTextQuaternary};
            font-size: 10px;
            line-height: 1;
            opacity: 0;
            vertical-align: 0;
            transition: color .16s ease, opacity .16s ease;
        }

        .log-copy:hover .ant-typography-copy,
        .log-copy .ant-typography-copy:focus-visible {
            color: ${token.colorPrimary};
            opacity: 1;
        }

        .log-type {
            max-width: 100%;
            display: inline-flex;
            align-items: center;
            gap: 5px;
            overflow: hidden;
            line-height: 20px;
            text-overflow: ellipsis;
            white-space: nowrap;
        }

        .log-type i {
            width: 6px;
            height: 6px;
            flex: none;
            border-radius: 50%;
            background: ${token.colorTextQuaternary};
        }

        .log-type.is-0 i {
            background: ${token.colorSuccess};
        }

        .log-type.is-1 i {
            background: ${token.colorPrimary};
        }

        .log-type.is-2 i {
            background: ${token.colorError};
        }

        .log-type.is-3 i {
            background: ${token.colorWarning};
        }

        .log-type.is-4 i {
            background: ${token.colorSuccess};
        }

        .log-text,
        .log-time {
            min-width: 0;
            display: block;
            overflow: hidden;
            line-height: 20px;
            text-overflow: ellipsis;
            white-space: nowrap;
        }

        @media (hover: none) {
            .log-copy .ant-typography-copy {
                opacity: .55;
            }
        }
    `,
}));

export default useStyles;
