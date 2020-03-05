<template>
  <a-list itemLayout="horizontal" :dataSource="bookDetail">
    <a-list-item slot="renderItem" slot-scope="item">
      <a-list-item-meta :description="item.AudioAbstract">
        <a slot="title" @click="addPlayList(item)">{{item.Title}}</a>
      </a-list-item-meta>
    </a-list-item>
  </a-list>
</template>

<script>
import Vue from "vue";
import { List } from "ant-design-vue";
import axios from "axios";
import config from "../config";

Vue.use(List);

export default {
  computed: {
    bookDetail() {
      return this.$store.state.idb.bookDetail.Details;
    },
  },
  async beforeMount() {
    await this.$store.dispatch(
      "checkExpire",
      "detail_" + this.$route.params.id
    );
    if (this.$store.state.idb.isExpire) {
      let res = await this.fetchData();
      await this.$store.dispatch("saveBookDetail", {
        id: this.$route.params.id,
        res,
      });
    }
    await this.$store.dispatch("getBookDetail", this.$route.params.id);
  },
  methods: {
    fetchData() {
      return axios
        .get(config.url.detail + this.$route.params.id)
        .then(function (res) {
          return res.data.Data;
        });
    },
    async addPlayList(item) {
      let detail = this.$store.state.idb.bookDetail;
      let aplay = {
        name: detail.Name + " " + item.Title,
        artist: detail.Author,
        url: item.FilePath,
        cover: detail.HomeImg,
      };
      this.$store.dispatch("saveAudioList", aplay);
    },
  },
};
</script>
<style></style>