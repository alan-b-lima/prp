function main() {
    const properties = {
        "darker": "--color-darker",
        "dark": "--color-dark",
        "neutral": "--color-neutral",
        "light": "--color-light",
        "lighter": "--color-lighter",
        "primary": "--color-primary",
        "secondary": "--color-secondary",
        "accent": "--color-accent",
    }

    const css = window.getComputedStyle(document.documentElement)

    for (const property in properties) {
        console.log(property, css.getPropertyValue(properties[property]))
    }
}

window.addEventListener("DOMContentLoaded", main)