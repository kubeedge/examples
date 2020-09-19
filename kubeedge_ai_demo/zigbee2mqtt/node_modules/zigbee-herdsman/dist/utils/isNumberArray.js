"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
// eslint-disable-next-line
function isNumberArray(value) {
    if (value instanceof Array) {
        for (const item of value) {
            if (typeof item !== 'number') {
                return false;
            }
        }
        return true;
    }
    return false;
}
exports.default = isNumberArray;
//# sourceMappingURL=isNumberArray.js.map