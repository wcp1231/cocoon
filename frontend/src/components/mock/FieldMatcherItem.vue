<template>
  <div class="p-fluid formgrid grid">
    <div class="field">
      <label>{{ field }}</label>
      <div class="p-inputgroup">
        <template v-if="matcher.equals">
          <Button label="Equals" class="p-button-info" disabled="disabled" />
          <InputText :value="matcher.equals" disabled />
        </template>
        <template v-else-if="matcher.regex">
          <Button label="Regex" class="p-button-success" disabled="disabled" />
          <InputText :value="matcher.regex" disabled />
        </template>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent } from "vue";
import Button from "primevue/button";
import InputText from 'primevue/inputtext';

export default defineComponent({
  name: "FieldMatcherItem",
  props: {
    field: String,
    matcher: Object,
  },
  data() {
    return {
    };
  },
  computed: {
  },
  methods: {
  },
  components: {
    Button,
    InputText,
  },
});
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>

</style>
