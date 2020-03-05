var store = {

    state: {
        aplayer: {},
        audio: [],
    },
    mutations: {
        setAplayer(state, aplayer) {
            state.aplayer = aplayer;
        },
        setAudioList(state, arr) {
            var existIndex = -1;
            for (let item in state.audio) {
                if (state.audio[item].name == arr.name) {
                    existIndex = Number(item);
                    break;
                }
            }
            if (existIndex < 0) {
                state.aplayer.audio.push(arr);
                state.aplayer.switch(arr.name)
            } else {
                state.aplayer.switch(existIndex)
            }
        },
    },
    actions: {
        async setAplayer({ commit }, aplayer) {
            commit('setAplayer', aplayer);
        },
        async saveAudioList({ commit }, arr) {
            commit('setAudioList', arr);
        },
        //TODO 一键加入播放列表、删除某一首、audio存入indexeddb、首页瀑布流、detail界面优化。

    },
}

export default store;