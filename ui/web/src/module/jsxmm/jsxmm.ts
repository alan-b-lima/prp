type Properties<T extends keyof HTMLElementTagNameMap> = RecursivePartial<HTMLElementTagNameMap[T]> & Record<string, any>

export function Element<T extends keyof HTMLElementTagNameMap>(tag: T, properties: Properties<T> = {}, ...children: (Node | string)[]): HTMLElementTagNameMap[T] {
    const element = document.createElement(tag)
    replace(element, properties)

    element.append(...children)
    return element
}

export function Style(element: HTMLElement, style: Partial<CSSStyleDeclaration>): void {
    replace(element.style, style)
}

function replace(base: Record<PropertyKey, any>, replacement: Record<PropertyKey, any>): void {
    for (const key in replacement) {
        if (!(key in base)) {
            console.error(`${key} not present in ${base} element`)
            // continue regardless, that's a problem for the caller
        }

        if (typeof replacement[key] === "object") {
            replace(base[key], replacement[key])
        } else {
            base[key] = replacement[key]
        }
    }
}

type CreateTupleOfLength<T extends number, U extends never[] = []> = U['length'] extends T ? U : CreateTupleOfLength<T, [...U, never]>

type Decrement<T extends number> = CreateTupleOfLength<T> extends [...infer Head, any] ? Head['length'] : 0

type RecursivePartial<T extends Object, C extends number = 1> = { [K in keyof T]?: T[K] extends Object ? (C extends 0 ? T[K] : RecursivePartial<T[K], Decrement<C>>) : T[K] }
