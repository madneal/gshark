const createBus = () => {
    const events = new Map()

    const on = (event, cb) => {
        const handlers = events.get(event) || new Set()
        handlers.add(cb)
        events.set(event, handlers)
    }

    const off = (event, cb) => {
        if (!cb) {
            events.delete(event)
            return
        }
        events.get(event)?.delete(cb)
    }

    const emit = (event, ...args) => {
        events.get(event)?.forEach(cb => cb(...args))
    }

    return {
        on,
        off,
        emit,
        $on: on,
        $off: off,
        $emit: emit,
    }
}

export const bus = createBus()

const install = (app) => {
    app.config.globalProperties.$bus = bus
}

export default install
