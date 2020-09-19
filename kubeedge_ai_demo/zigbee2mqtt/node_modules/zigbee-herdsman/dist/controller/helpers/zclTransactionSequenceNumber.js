"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
class ZclTransactionSequenceNumber {
    static next() {
        this.number++;
        if (this.number > 255) {
            this.number = 1;
        }
        return this.number;
    }
}
ZclTransactionSequenceNumber.number = 1;
exports.default = ZclTransactionSequenceNumber;
//# sourceMappingURL=zclTransactionSequenceNumber.js.map