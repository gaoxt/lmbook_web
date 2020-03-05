<template>
  <a-list-item slot="renderItem" slot-scope="source">
    <router-link :to="'/detail/' + source.id">
      <a-card hoverable :bordered="false">
        <img :src="source.HomeImg" slot="cover" />
        <a-card-meta :title="source.Name">
          <template slot="description" class="description_height">{{source.Abstract}}</template>
        </a-card-meta>
      </a-card>
    </router-link>
  </a-list-item>
</template>

<script>
import Vue from "vue";
import { List, Card, Spin } from "ant-design-vue";
// import axios from "axios";

Vue.use(List);
Vue.use(Card);
Vue.use(Spin);

export default {
  name: "item-component",
  props: {
    index: {
      // index of current item
      type: Number,
    },
    source: {
      // here is: {uid: 'unique_1', text: 'abc'}
      type: Object,
      default() {
        return {};
      },
    },
  },
};
</script>