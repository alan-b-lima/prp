export function Element(tag, properties = {}, ...children) {
    const element = document.createElement(tag);
    replace(element, properties);
    element.append(...children);
    return element;
}
export function Style(element, style) {
    replace(element.style, style);
}
function replace(base, replacement) {
    for (const key in replacement) {
        if (!(key in base)) {
            console.error(`${key} not present in ${base} element`);
        }
        if (typeof replacement[key] === "object") {
            replace(base[key], replacement[key]);
        }
        else {
            base[key] = replacement[key];
        }
    }
}
//# sourceMappingURL=jsxmm.js.map