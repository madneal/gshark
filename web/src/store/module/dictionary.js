import { findSysDictionary } from '@/api/sysDictionary'

export const dictionary = {
    namespaced: true,
    state: {
        dictionaryMap: {},
    },
    mutations: {
        setDictionaryMap(state, dictionaryMap) {
            state.dictionaryMap = { ...state.dictionaryMap, ...dictionaryMap }
        },

    },
    actions: {
        // 从后台获取动态路由
        async getDictionary({ commit, state }, type) {
            if (state.dictionaryMap[type]) {
                return state.dictionaryMap[type]
            } else {
                const res = await findSysDictionary({ type })
                if (res.code == 0) {
                    const dictionaryMap = {}
                    const dict = []
                    res.data.resysDictionary.sysDictionaryDetails && res.data.resysDictionary.sysDictionaryDetails.map(item => {
                        dict.push({
                            label: item.label,
                            value: item.value
                        })
                    })
                    dictionaryMap[res.data.resysDictionary.type] = dict
                    commit("setDictionaryMap", dictionaryMap)
                    return state.dictionaryMap[type]
                }
            }
        }
    },
    getters:{
        getDictionary(state){
            return state.dictionaryMap
        }
    }
}