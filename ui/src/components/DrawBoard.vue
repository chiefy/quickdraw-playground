<template>
  <div
    class="pick"
    v-bind:class="{'clear': num % 10 == 1 }"
  >
    <q-btn
      round
      :color="color"
      :label="num"
    >
      <q-tooltip>
        {{ picks }}
      </q-tooltip>
    </q-btn>
  </div>
</template>

<script>
import { QBtn, QTooltip } from 'quasar'
export default {
  name: 'DrawBoard',
  components: {
    QBtn,
    QTooltip
  },
  props: [
    'num',
    'picks',
    'color'
  ]
}
</script>
