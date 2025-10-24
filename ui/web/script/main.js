const THEME_TEXT = {
    "light": "🌞",
    "dark": "🌑",
}

function main() {
    const light_dark = document.getElementById("light-dark")
    if (light_dark !== null) {
        light_dark.addEventListener("click", listener_light_dark.bind(null, light_dark))
        listener_light_dark(light_dark)
    }

    const tabs = document.querySelector(".tabs")
    if (tabs !== null) {
        window.addEventListener("resize", listener_is_overflowing.bind(null, tabs, "overflowing"));
    }
}

function listener_light_dark(element) {
    const html_classes = document.documentElement.classList
    let theme

    LookForTheme:
    {
        if (html_classes.contains("light")) {
            html_classes.remove("light")

            theme = "dark"
            break LookForTheme
        }

        if (html_classes.contains("dark")) {
            html_classes.remove("dark")

            theme = "light"
            break LookForTheme
        }

        let light_dark = localStorage.getItem("light-dark")

        if (light_dark !== "light" && light_dark !== "dark") {
            if (window.matchMedia && window.matchMedia("(prefers-color-scheme: dark)").matches) {
                light_dark = "dark"
            } else {
                light_dark = "light"
            }
        }

        theme = light_dark
    }

    html_classes.add(theme)
    element.textContent = THEME_TEXT[theme]
    localStorage.setItem("light-dark", theme)
}

function listener_is_overflowing(element, clazz) {
    const rect = element.getBoundingClientRect()

    const top = Math.min(0, rect.top) / rect.height >= 0
    const bottom = (rect.bottom - window.innerHeight) / rect.height >= 0

    if (top && bottom) {
        element.classList.add(clazz);
    } else {
        element.classList.remove(clazz);
    }
}

window.addEventListener("DOMContentLoaded", main)