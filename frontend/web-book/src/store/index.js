import Vue from 'vue'
import Vuex from 'vuex'

import idb from './modules/idb';
import aplay from './modules/aplay';

Vue.use(Vuex)

export default new Vuex.Store({
    modules: {
        ...idb,
        aplay
    },
    debug: true,
})

