"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.compact = void 0;
function compact(options) {
    return options.fn(this)
        .split('\n')
        .map(line => line.trim())
        .join('')
        .replace(/&nbsp;/g, ' ')
        .trim();
}
exports.compact = compact;
//# sourceMappingURL=compact.js.map