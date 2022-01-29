<template>
  <Layout />
</template>

<script lang="ts">
import { defineComponent } from "vue";
import Layout from "./components/Layout.vue";

export default defineComponent({
  name: "App",
  components: {
    Layout,
  },
});
</script>

<style>
body {
  margin: 0;
}
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}
</style>
