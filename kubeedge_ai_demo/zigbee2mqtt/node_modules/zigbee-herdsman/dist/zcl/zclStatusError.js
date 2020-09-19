"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
const status_1 = __importDefault(require("./definition/status"));
class ZclStatusError extends Error {
    constructor(code) {
        super(`Status '${status_1.default[code]}'`);
        this.code = code;
    }
}
exports.default = ZclStatusError;
//# sourceMappingURL=zclStatusError.js.map