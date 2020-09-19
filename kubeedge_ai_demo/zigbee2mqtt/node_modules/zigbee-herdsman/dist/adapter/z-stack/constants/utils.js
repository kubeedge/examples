"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.statusDescription = exports.getChannelMask = void 0;
const common_1 = require("./common");
function getChannelMask(channels) {
    const value = channels.reduce((mask, channel) => mask | (1 << channel), 0);
    return [value & 0xFF, (value >> 8) & 0xFF, (value >> 16) & 0xFF, (value >> 24) & 0xFF];
}
exports.getChannelMask = getChannelMask;
function statusDescription(code) {
    const hex = "0x" + code.toString(16).padStart(2, "0");
    return `(${hex}: ${common_1.ZnpCommandStatus[code] || "UNKNOWN"})`;
}
exports.statusDescription = statusDescription;
//# sourceMappingURL=utils.js.map