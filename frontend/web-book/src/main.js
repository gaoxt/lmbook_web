import Vue from 'vue';
import router from './router.js'

import App from './App';
import store from "./store/index"
import Aplayer from '@moefe/vue-aplayer'



Vue.config.productionTip = false;

Vue.use(Aplayer, {
  defaultCover: "https://github.com/u3u.png",
  productionTip: false,
});


new Vue({
  router,
  store,
  render: h => h(App)
}).$mount('#app');