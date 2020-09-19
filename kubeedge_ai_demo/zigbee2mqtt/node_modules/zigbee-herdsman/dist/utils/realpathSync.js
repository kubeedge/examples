"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
/* istanbul ignore file */
const fs_1 = __importDefault(require("fs"));
/* Only used for mocking purposes */
function realpathSync(path) {
    return fs_1.default.realpathSync(path);
}
exports.default = realpathSync;
//# sourceMappingURL=realpathSync.js.map