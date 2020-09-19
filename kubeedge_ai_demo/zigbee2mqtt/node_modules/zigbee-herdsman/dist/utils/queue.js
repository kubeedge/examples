"use strict";
var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    function adopt(value) { return value instanceof P ? value : new P(function (resolve) { resolve(value); }); }
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : adopt(result.value).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
};
Object.defineProperty(exports, "__esModule", { value: true });
class Queue {
    constructor(concurrent = 1) {
        this.jobs = [];
        this.concurrent = concurrent;
    }
    execute(func, key = null) {
        return new Promise((resolve, reject) => {
            this.jobs.push({ key, func, running: false, resolve, reject });
            this.executeNext();
        });
    }
    executeNext() {
        return __awaiter(this, void 0, void 0, function* () {
            const job = this.getNext();
            if (job) {
                job.running = true;
                try {
                    const result = yield job.func();
                    this.jobs.splice(this.jobs.indexOf(job), 1);
                    job.resolve(result);
                    this.executeNext();
                }
                catch (error) {
                    this.jobs.splice(this.jobs.indexOf(job), 1);
                    job.reject(error);
                    this.executeNext();
                }
            }
        });
    }
    getNext() {
        if (this.jobs.filter((j) => j.running).length > (this.concurrent - 1)) {
            return null;
        }
        for (let i = 0; i < this.jobs.length; i++) {
            const job = this.jobs[i];
            if (!job.running && (!job.key || !this.jobs.find((j) => j.key === job.key && j.running))) {
                return job;
            }
        }
        return null;
    }
    clear() {
        this.jobs = [];
    }
}
exports.default = Queue;
//# sourceMappingURL=queue.js.map