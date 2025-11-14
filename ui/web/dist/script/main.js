import { MockUserGateway } from "./internal/domain/gateway/mock.js";
import { Element } from "./module/jsxmm/jsxmm.js";
import { setup_theme } from "./theme.js";
async function main() {
    setup_theme();
    if (location.pathname.endsWith("static.html")) {
        const api = new MockUserGateway();
        await CreateUsers(api);
        await Users(api);
    }
}
async function CreateUsers(api) {
    await api.Create({ name: "Alan Lima", login: "alan-b-lima", password: "12345678" });
    await api.Create({ name: "Juan Ferreira", login: "juanzinho_bs", password: "12345678" });
    await api.Create({ name: "Luan Filipe", login: "cr-lf", password: "12345678" });
    await api.Create({ name: "Vitor Moisés", login: "vecto", password: "12345678" });
}
async function Users(api) {
    const content = document.querySelector('.content');
    const resp = await api.List({ offset: 0, limit: 10 });
    if (resp instanceof Error) {
        return;
    }
    for (const user of resp.records) {
        console.log(user);
        content.append(UserComponent(user));
    }
}
function UserComponent(user) {
    return (Element("div", { className: "user" }, Element("div", {}, user.uuid), Element("div", {}, user.name), Element("div", {}, user.login, " • ", user.level)));
}
window.addEventListener("DOMContentLoaded", main);
//# sourceMappingURL=main.js.map