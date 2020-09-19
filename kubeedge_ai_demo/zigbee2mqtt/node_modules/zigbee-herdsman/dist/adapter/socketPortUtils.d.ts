declare function isTcpPath(path: string): boolean;
declare function parseTcpPath(path: string): {
    host: string;
    port: number;
};
declare const _default: {
    isTcpPath: typeof isTcpPath;
    parseTcpPath: typeof parseTcpPath;
};
export default _default;
//# sourceMappingURL=socketPortUtils.d.ts.map