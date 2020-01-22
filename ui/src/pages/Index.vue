<template>
  <q-page class="flex flex-center">
    <div class="q-pa-md q-gutter-sm">
      <h3>{{ viewType }}</h3>
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

const numColors = 10
const apiURL = process.env.API || ''

export default {
  name: 'Main',
  components: {
    DrawBoard
  },
  data () {
    return this.defaultData()
  },
  watch: {
    $route (to, from) {
      this.getViewData()
    }
  },
  methods: {
    defaultData () {
      return {
        nums: [],
        leastPicks: 0,
        mostPicks: 0,
        pickDiff: 0,
        viewType: null
      }
    },
    getColor (numPicks) {
      let pickGroups = Math.ceil(this.pickDiff / numColors)
      let colorIdx = numColors - Math.round((this.mostPicks - numPicks) / pickGroups)
      colorIdx = colorIdx === 0 ? 1 : colorIdx
      return 'blue-' + colorIdx
    },
    resetData () {
      Object.assign(this.$data, this.defaultData())
    },
    getViewData () {
      this.resetData()
      this.viewType = this.$route.meta.viewTypeTitles[this.$route.params.viewType]
      return this.$axios.get(apiURL + '/freq/' + this.$route.params.viewType)
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
            message: 'Loading Data Failed :(',
            icon: 'report_problem'
          })
        })
    }
  },
  mounted () {
    this.getViewData()
  }
}
</script>
