import { setup_theme } from "./theme.js";
function main() {
    setup_theme();
    const tabs = document.querySelector(".tabs");
    if (tabs !== null) {
        window.addEventListener("resize", listener_is_overflowing.bind(null, tabs, "overflowing"));
    }
}
function listener_is_overflowing(element, clazz) {
    const rect = element.getBoundingClientRect();
    const top = Math.min(0, rect.top) / rect.height >= 0;
    const bottom = (rect.bottom - window.innerHeight) / rect.height >= 0;
    if (top && bottom) {
        element.classList.add(clazz);
    }
    else {
        element.classList.remove(clazz);
    }
}
window.addEventListener("DOMContentLoaded", main);
