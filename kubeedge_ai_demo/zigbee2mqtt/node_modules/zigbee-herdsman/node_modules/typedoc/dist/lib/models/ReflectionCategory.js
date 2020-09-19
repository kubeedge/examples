"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.ReflectionCategory = void 0;
class ReflectionCategory {
    constructor(title) {
        this.children = [];
        this.title = title;
        this.allChildrenHaveOwnDocument = (() => this.getAllChildrenHaveOwnDocument());
    }
    getAllChildrenHaveOwnDocument() {
        let onlyOwnDocuments = true;
        this.children.forEach((child) => {
            onlyOwnDocuments = onlyOwnDocuments && !!child.hasOwnDocument;
        });
        return onlyOwnDocuments;
    }
}
exports.ReflectionCategory = ReflectionCategory;
//# sourceMappingURL=ReflectionCategory.js.map