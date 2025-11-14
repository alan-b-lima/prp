namespace user {
    interface Gateway {
        List(req: ListRequest): PromiseResult<ListResponse>
        Get(req: GetRequest): PromiseResult<Response>
        GetByLogin(req: GetByLoginRequest): PromiseResult<Response>
        Create(req: CreateRequest): PromiseResult<Response>
        Patch(req: PatchRequest): PromiseResult<Response>
        Delete(req: DeleteRequest): PromiseResult<void>
        Authenticate(req: AuthRequest): PromiseResult<AuthResponse>
    }

    type ListRequest = {
        offset: number
        limit: number
    }

    type GetRequest = {
        uuid: string
    }

    type GetByLoginRequest = {
        login: string
    }

    type CreateRequest = {
        name: string
        login: string
        password: string
    }

    type PatchRequest = {
        uuid: string
        name?: string
        login?: string
        password?: string
    }

    type DeleteRequest = {
        uuid: string
    }

    type AuthRequest = {
        login: string
        password: string
    }

    type ContextRequest = {
        uuid: string
    }

    type ListResponse = {
        offset: number
        length: number
        records: Response[]
        total_records: number
    }

    type AuthResponse = {
        uuid: string
        user: string
        expires: Date
    }

    type Response = {
        uuid: string
        name: string
        login: string
        level: string
    }
}