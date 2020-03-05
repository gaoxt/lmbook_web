<template>
  <a-list :data-source="data" :grid="{ xs: 1, sm: 2, md: 4, lg: 4, xl: 6, xxl: 3 }">
    <virtual-list
      style="height: 700px; overflow-y: auto;"
      :data-key="'id'"
      :data-sources="data"
      :data-component="itemComponent"
      :estimate-size="10"
    />
    <a-spin v-if="loading && !busy" class="demo-loading" />
  </a-list>
</template>

<script>
import Item from "./Item";
import VirtualList from "vue-virtual-scroll-list";

export default {
  name: "scroll",
  data() {
    return {
      data: [],
      loading: false,
      busy: false,
      itemComponent: Item,
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
    this.loading = true;
    await this.$store.dispatch("getBooks");
    this.loading = false;
    this.data = this.$store.state.idb.bookList.map((item, index) => ({
      ...item,
      index,
    }));
    console.log(this.data);
  },
  components: { "virtual-list": VirtualList },
};
</script>