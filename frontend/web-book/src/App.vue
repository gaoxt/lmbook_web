<template>
  <div id="app">
    <aplayer :audio="audio" ref="aplayer" fixed />
    <router-view></router-view>
  </div>
</template>

<script>
import { mapState } from "vuex";

export default {
  computed: {
    ...mapState({
      audio: (state) => state.aplay.audio,
    }),
  },
  mounted() {
    console.log(this.$refs.aplayer);
    this.$store.commit("setAplayer", this.$refs.aplayer);
  },
};
</script>

<style>
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  margin-top: 60px;
}
</style>
