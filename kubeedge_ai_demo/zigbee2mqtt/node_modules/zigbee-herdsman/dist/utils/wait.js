"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
function wait(milliseconds) {
    return new Promise((resolve) => {
        setTimeout(() => resolve(), milliseconds);
    });
}
exports.default = wait;
//# sourceMappingURL=wait.js.map