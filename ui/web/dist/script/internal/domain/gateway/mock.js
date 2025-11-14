export class MockUserGateway {
    #users;
    #actor;
    constructor() {
        const users = sessionStorage.getItem("users");
        if (users !== null) {
            this.#users = JSON.parse(users);
        }
        else {
            this.#users = [];
        }
        this.#actor = null;
    }
    async List(req) {
        const records = this.#users.slice(req.offset, req.offset + req.limit);
        const resp = {
            offset: req.offset,
            length: records.length,
            records: records,
            total_records: this.#users.length,
        };
        return resp;
    }
    async Get(req) {
        for (const user of this.#users) {
            if (req.uuid === user.uuid) {
                return user;
            }
        }
        return new Error("user not found");
    }
    async GetByLogin(req) {
        for (const user of this.#users) {
            if (req.login === user.login) {
                return user;
            }
        }
        return new Error("user not found");
    }
    async Create(req) {
        const resp_get = await this.GetByLogin({ login: req.login });
        if (!(resp_get instanceof Error)) {
            return new Error("user found");
        }
        const resp = {
            uuid: crypto.randomUUID(),
            name: req.name,
            login: req.login,
            level: "user",
        };
        this.#users.push(resp);
        sessionStorage.setItem("users", JSON.stringify(this.#users));
        return resp;
    }
    async Patch(req) {
        const resp = await this.Get({ uuid: req.uuid });
        if (resp instanceof Error) {
            return resp;
        }
        if (req.name !== undefined) {
            resp.name = req.name;
        }
        if (req.login !== undefined) {
            resp.login = req.login;
        }
        sessionStorage.setItem("users", JSON.stringify(this.#users));
        return resp;
    }
    async Delete(req) {
        for (let i = 0; i < this.#users.length; i++) {
            if (req.uuid !== this.#users[i].uuid) {
                continue;
            }
            this.#users = [...this.#users.slice(0, i), ...this.#users.slice(i + 1)];
        }
        sessionStorage.setItem("users", JSON.stringify(this.#users));
        return new Error("user not found");
    }
    async Authenticate(req) {
        const resp_get = await this.GetByLogin({ login: req.login });
        if (resp_get instanceof Error) {
            return resp_get;
        }
        this.#actor = resp_get;
        const resp = {
            uuid: crypto.randomUUID(),
            user: resp_get.uuid,
            expires: new Date(Date.now()),
        };
        return resp;
    }
    async Me(req) {
        if (this.#actor === null) {
            return new Error("unautheticated user");
        }
        return this.#actor;
    }
}
//# sourceMappingURL=mock.js.map