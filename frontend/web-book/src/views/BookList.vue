<template>
  <div
    v-infinite-scroll="handleInfiniteOnLoad"
    class="demo-infinite-container"
    :infinite-scroll-disabled="busy"
    :infinite-scroll-distance="10"
  >
    <a-list :data-source="data" :grid="{ xs: 1, sm: 2, md: 4, lg: 4, xl: 6, xxl: 3 }">
      <a-list-item slot="renderItem" slot-scope="item">
        <router-link :to="'/detail/' + item.id">
          <a-card hoverable :bordered="false">
            <img :src="item.HomeImg" slot="cover" />
            <a-card-meta :title="item.Name">
              <template slot="description" class="description_height">{{item.Abstract}}</template>
            </a-card-meta>
          </a-card>
        </router-link>
      </a-list-item>
      <a-spin v-if="loading && !busy" class="demo-loading" />
    </a-list>
  </div>
</template>

<script>
import Vue from "vue";
import { List, Card, Spin } from "ant-design-vue";
import axios from "axios";
import infiniteScroll from "vue-infinite-scroll";
import config from "../config";

Vue.use(List);
Vue.use(Card);
Vue.use(Spin);

export default {
  directives: { infiniteScroll },
  data() {
    return {
      data: [],
      loading: false,
      busy: false,
      page: 1,
      pageSize: 18,
    };
  },
  async beforeMount() {
    await this.$store.dispatch("checkExpire", "books");
    if (this.$store.state.idb.isExpire) {
      let res = await this.fetchData();
      await this.$store.dispatch("saveBook", res);
    }
    await this.$store.dispatch("getBooks");
    this.handleInfiniteOnLoad();
  },
  methods: {
    cutData() {
      let startIndex = (this.page - 1) * this.pageSize;
      let endIndex = startIndex + this.pageSize;
      this.page++;
      return this.$store.state.idb.bookList.slice(startIndex, endIndex);
    },
    fetchData() {
      return axios.get(config.url.list).then(function (res) {
        return res.data.Data;
      });
    },
    handleInfiniteOnLoad() {
      this.loading = true;
      let data = this.data;
      if (data.length >= this.$store.state.idb.bookList.length) {
        this.busy = true;
        this.loading = false;
        return;
      }
      this.data = data
        .concat(this.cutData())
        .map((item, index) => ({ ...item, index }));
      this.loading = false;
    },
  },
};
</script>
<style>
.ant-card-meta-description {
  height: 40px;
  color: rgba(0, 0, 0, 0.45);
}
.demo-loading {
  position: absolute;
  bottom: 40px;
  width: 100%;
  text-align: center;
}
</style>