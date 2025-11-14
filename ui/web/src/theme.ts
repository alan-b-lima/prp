const THEME_TEXT: Record<string, string> = {
    light: "🌞",
    dark: "🌑",
} as const

export function setup_theme(): void {
    const theme = document.getElementById("theme")
    if (theme !== null) {
        theme.addEventListener("click", listener_theme.bind(null, theme))
        listener_theme(theme)
    }
}

function listener_theme(element: HTMLElement): void {
    const html_classes = document.documentElement.classList
    let theme: string

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

        let light_dark = localStorage.getItem("theme")

        if (light_dark !== "light" && light_dark !== "dark") {
            if (window.matchMedia("(prefers-color-scheme: dark)").matches) {
                light_dark = "dark"
            } else {
                light_dark = "light"
            }
        }

        theme = light_dark
    }

    html_classes.add(theme)
    element.textContent = THEME_TEXT[theme]
    localStorage.setItem("theme", theme)
}
