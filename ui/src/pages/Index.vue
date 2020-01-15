<template>
  <q-page class="flex flex-center">
    <div class="q-pa-md q-gutter-sm">
      <draw-board
        v-for="p in nums"
        v-bind:key="p.num"
        v-bind:num="p.num"
        v-bind:picks="p.picks"
        v-bind:color="getColor(p.picks)"
        ref="drawBoard"
      ></draw-board>
    </div>
  </q-page>
</template>

<script>
import DrawBoard from 'components/DrawBoard'

export default {
  name: 'Main',
  components: {
    DrawBoard
  },
  data () {
    return {
      nums: [],
      leastPicks: 0,
      mostPicks: 0,
      pickDiff: 0
    }
  },
  methods: {
    getColor (numPicks) {
      let calc = (numPicks / this.pickDiff)
      console.log('picks = ' + calc + ' ' + numPicks)
      return 'blue-1'
    }
  },
  mounted () {
    return this.$axios.get('http://localhost:9090/freq/onemonth')
      .then((response) => {
        const addPick = (n) => {
          let numPicks = response.data[n]
          if (numPicks > this.mostPicks) {
            this.mostPicks = numPicks
          }
          if (numPicks < this.leastPicks || this.leastPicks === 0) {
            this.leastPicks = numPicks
          }
          this.nums.push({
            num: n,
            picks: numPicks
          })
        }
        this.lodash.each(this.lodash.range(1, 81), addPick.bind(this))
        this.pickDiff = this.mostPicks - this.leastPicks
      })
      .catch(() => {
        this.$q.notify({
          color: 'negative',
          position: 'top',
          message: 'Loading failed',
          icon: 'report_problem'
        })
      })
  }
}
</script>
