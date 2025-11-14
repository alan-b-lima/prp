import { MockUserGateway } from "./internal/domain/gateway/mock.ts"
import { Element } from "./module/jsxmm/jsxmm.ts"
import { setup_theme } from "./theme.ts"

async function main(): Promise<void> {
    setup_theme()

    if (location.pathname.endsWith("static.html")) {
        const api = new MockUserGateway()
        console.log(api)

        await CreateUsers(api)
        await Users(api)
    }
}

async function CreateUsers(api: user.Gateway) {
    await api.Create({ name: "Alan Lima", login: "alan-b-lima", password: "12345678" })
    await api.Create({ name: "Juan Ferreira", login: "juanzinho_bs", password: "12345678" })
    await api.Create({ name: "Luan Filipe", login: "cr-lf", password: "12345678" })
    await api.Create({ name: "Vitor Moisés", login: "vecto", password: "12345678" })
}

async function Users(api: user.Gateway) {
    const content = document.querySelector('.content')! as HTMLDivElement

    const resp = await api.List({ offset: 0, limit: 10 })
    if (resp instanceof Error) {
        return
    }

    for (const user of resp.records) {
        content.append(UserComponent(user))
    }
}

function UserComponent(user: user.Response) {
    return (
        Element("div", { className: "user" },
            Element("div", {}, user.uuid),
            Element("div", {}, user.name),
            Element("div", {}, user.login, " • ", user.level),
        )
    )
}


window.addEventListener("DOMContentLoaded", main)