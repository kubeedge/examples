"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
class Entity {
    static injectDatabase(database) {
        Entity.database = database;
    }
    static injectAdapter(adapter) {
        Entity.adapter = adapter;
    }
}
Entity.database = null;
Entity.adapter = null;
exports.default = Entity;
//# sourceMappingURL=entity.js.map