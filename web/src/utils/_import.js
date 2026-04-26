const modules = import.meta.glob('../view/**/*.vue')

export default file => {
    const normalizedFile = file.startsWith('/') ? file.slice(1) : file
    const loader = modules[`../${normalizedFile}`]

    if (!loader) {
        return () => import('@/view/error/index.vue')
    }

    return loader
}
