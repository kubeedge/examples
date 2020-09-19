"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
function isTcpPath(path) {
    // tcp path must be:
    // tcp://<host>:<port>
    const regex = /^(?:tcp:\/\/)[\w.-]+[:][\d]+$/gm;
    return regex.test(path);
}
function parseTcpPath(path) {
    const str = path.replace("tcp://", "");
    return {
        host: str.substring(0, str.indexOf(":")),
        port: Number(str.substring(str.indexOf(":") + 1)),
    };
}
exports.default = { isTcpPath, parseTcpPath };
//# sourceMappingURL=socketPortUtils.js.map