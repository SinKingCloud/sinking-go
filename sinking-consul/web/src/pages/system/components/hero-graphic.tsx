import React from "react";
import {createStyles} from "antd-style";

const useStyles = createStyles(({css, token}: any) => ({
    graphic: css`
        width: 100%;
        height: 100%;
        display: block;
        overflow: visible;
        color: ${token.colorPrimary};

        .flow-line {
            fill: none;
            stroke-linecap: round;
            stroke-linejoin: round;
        }

        .flow-primary {
            stroke-width: 1.6;
        }

        .flow-secondary {
            stroke: currentColor;
            stroke-width: 1;
            stroke-opacity: .16;
        }

        .flow-dash {
            stroke: currentColor;
            stroke-width: 1;
            stroke-dasharray: 3 10;
            stroke-opacity: .22;
        }

        .flow-signal {
            stroke: currentColor;
            stroke-width: 2.1;
            stroke-dasharray: 15 190;
            stroke-opacity: .48;
            animation: heroSignal 16s linear infinite;
        }

        .node-halo {
            fill: currentColor;
            fill-opacity: .055;
            stroke: currentColor;
            stroke-opacity: .13;
        }

        .node-core {
            fill: currentColor;
            fill-opacity: .52;
        }

        .flow-scene {
            transform-box: fill-box;
            transform-origin: center;
            animation: heroDrift 16s ease-in-out infinite;
        }

        .flow-node-a,
        .flow-node-c {
            animation: heroNode 7s ease-in-out infinite;
        }

        .flow-node-c {
            animation-delay: -3.5s;
        }

        .setting-orbit,
        .setting-trace {
            fill: none;
            stroke: currentColor;
            stroke-linecap: round;
        }

        .setting-orbit {
            stroke-dasharray: 4 10;
            stroke-opacity: .16;
        }

        .setting-trace {
            stroke-width: 1.3;
        }

        .setting-sheet {
            fill: currentColor;
            fill-opacity: .04;
            stroke: currentColor;
            stroke-opacity: .18;
        }

        .setting-sheet.front {
            fill-opacity: .09;
            stroke-opacity: .34;
        }

        .setting-chip {
            fill: currentColor;
            fill-opacity: .3;
        }

        .setting-chip.muted {
            fill-opacity: .1;
        }

        .setting-point-halo {
            fill: currentColor;
            fill-opacity: .045;
            stroke: currentColor;
            stroke-opacity: .12;
        }

        .setting-point-core {
            fill: currentColor;
            fill-opacity: .42;
        }

        .setting-stack {
            transform-box: fill-box;
            transform-origin: center;
            animation: settingFloat 18s ease-in-out infinite;
        }

        .setting-point-a,
        .setting-point-b {
            animation: heroNode 8s ease-in-out infinite;
        }

        .setting-point-b {
            animation-delay: -4s;
        }

        @keyframes heroSignal {
            to {
                stroke-dashoffset: -205;
            }
        }

        @keyframes heroDrift {
            0%, 100% {
                transform: translate3d(0, 0, 0);
            }
            50% {
                transform: translate3d(6px, -3px, 0);
            }
        }

        @keyframes heroNode {
            0%, 100% {
                opacity: .62;
            }
            50% {
                opacity: 1;
            }
        }

        @keyframes settingFloat {
            0%, 100% {
                transform: translate3d(0, 0, 0);
            }
            50% {
                transform: translate3d(-4px, 3px, 0);
            }
        }

        @media (prefers-reduced-motion: reduce) {
            .flow-scene,
            .flow-signal,
            .flow-node-a,
            .flow-node-c,
            .setting-stack,
            .setting-point-a,
            .setting-point-b {
                animation: none;
            }
        }
    `,
}));

const TaskGraphic = ({styles, id}: any): React.ReactNode => (
    <svg className={`${styles.graphic} hero-flow-graphic`} viewBox="0 0 480 132"
         aria-hidden="true" focusable="false">
        <defs>
            <linearGradient id={`${id}-flow`} x1="38" y1="90" x2="456" y2="52"
                            gradientUnits="userSpaceOnUse">
                <stop stopColor="currentColor" stopOpacity="0"/>
                <stop offset=".24" stopColor="currentColor" stopOpacity=".24"/>
                <stop offset=".58" stopColor="currentColor" stopOpacity=".7"/>
                <stop offset="1" stopColor="currentColor" stopOpacity=".04"/>
            </linearGradient>
            <radialGradient id={`${id}-aura`} cx="0" cy="0" r="1"
                            gradientTransform="translate(322 63) rotate(90) scale(60 132)"
                            gradientUnits="userSpaceOnUse">
                <stop stopColor="currentColor" stopOpacity=".12"/>
                <stop offset=".55" stopColor="currentColor" stopOpacity=".035"/>
                <stop offset="1" stopColor="currentColor" stopOpacity="0"/>
            </radialGradient>
        </defs>
        <g className="flow-scene">
            <ellipse className="flow-aura" cx="322" cy="63" rx="132" ry="60"
                     fill={`url(#${id}-aura)`}/>
            <path className="flow-line flow-secondary flow-detail"
                  d="M34 108C116 87 139 29 229 35S346 107 468 52"/>
            <path className="flow-line flow-dash flow-detail"
                  d="M60 30C143 28 174 84 255 86S371 38 454 72"/>
            <path className="flow-line flow-primary" stroke={`url(#${id}-flow)`}
                  d="M26 88C106 84 143 49 215 51S324 94 457 58"/>
            <path className="flow-line flow-signal"
                  d="M26 88C106 84 143 49 215 51S324 94 457 58"/>
            <g className="flow-node flow-node-a" transform="translate(154 57)">
                <circle className="node-halo" r="13"/>
                <circle className="node-core" r="3.4"/>
            </g>
            <g className="flow-node flow-node-b flow-secondary" transform="translate(282 78)">
                <circle className="node-halo" r="10"/>
                <circle className="node-core" r="2.8"/>
            </g>
            <g className="flow-node flow-node-c" transform="translate(390 67)">
                <circle className="node-halo" r="14"/>
                <circle className="node-core" r="3.6"/>
            </g>
        </g>
    </svg>
);

const SettingGraphic = ({styles, id}: any): React.ReactNode => (
    <svg className={`${styles.graphic} hero-setting-graphic`} viewBox="0 0 480 132"
         aria-hidden="true" focusable="false">
        <defs>
            <linearGradient id={`${id}-edge`} x1="76" y1="101" x2="438" y2="28"
                            gradientUnits="userSpaceOnUse">
                <stop stopColor="currentColor" stopOpacity="0"/>
                <stop offset=".44" stopColor="currentColor" stopOpacity=".28"/>
                <stop offset="1" stopColor="currentColor" stopOpacity=".04"/>
            </linearGradient>
            <radialGradient id={`${id}-aura`} cx="0" cy="0" r="1"
                            gradientTransform="translate(316 64) rotate(90) scale(60 138)"
                            gradientUnits="userSpaceOnUse">
                <stop stopColor="currentColor" stopOpacity=".11"/>
                <stop offset=".62" stopColor="currentColor" stopOpacity=".03"/>
                <stop offset="1" stopColor="currentColor" stopOpacity="0"/>
            </radialGradient>
        </defs>
        <ellipse className="setting-aura" cx="316" cy="64" rx="138" ry="60"
                 fill={`url(#${id}-aura)`}/>
        <path className="setting-orbit setting-detail"
              d="M70 102C139 77 175 31 254 30S370 77 454 43"/>
        <path className="setting-trace" stroke={`url(#${id}-edge)`}
              d="M96 93C163 86 193 63 238 63"/>
        <g className="setting-stack">
            <rect className="setting-sheet setting-detail" x="238" y="23" width="142" height="82" rx="22"
                  transform="rotate(-7 309 64)"/>
            <rect className="setting-sheet setting-detail" x="270" y="22" width="132" height="84" rx="22"
                  transform="rotate(6 336 64)"/>
            <rect className="setting-sheet front" x="250" y="27" width="146" height="78" rx="21"/>
            <rect className="setting-chip" x="270" y="46" width="42" height="8" rx="4"/>
            <rect className="setting-chip muted" x="321" y="46" width="54" height="8" rx="4"/>
            <rect className="setting-chip muted" x="270" y="64" width="68" height="8" rx="4"/>
            <rect className="setting-chip" x="347" y="64" width="28" height="8" rx="4"/>
            <rect className="setting-chip muted" x="270" y="82" width="32" height="8" rx="4"/>
            <rect className="setting-chip muted" x="311" y="82" width="64" height="8" rx="4"/>
        </g>
        <g className="setting-point setting-point-a" transform="translate(181 75)">
            <circle className="setting-point-halo" r="12"/>
            <circle className="setting-point-core" r="3.2"/>
        </g>
        <g className="setting-point setting-point-b setting-detail" transform="translate(430 38)">
            <circle className="setting-point-halo" r="9"/>
            <circle className="setting-point-core" r="2.6"/>
        </g>
    </svg>
);

const HeroGraphic = ({variant = "task"}: any): React.ReactNode => {
    const {styles} = useStyles();
    const id = React.useId().replace(/:/g, "");
    return variant === "setting"
        ? <SettingGraphic styles={styles} id={`setting-${id}`}/>
        : <TaskGraphic styles={styles} id={`task-${id}`}/>;
};

export default React.memo(HeroGraphic);
