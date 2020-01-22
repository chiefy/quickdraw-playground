<template>
  <q-page class="flex flex-center">
    <div class="hidden">
      <q-date
        subtitle="Start Date"
        v-model="startDate"
      />
      <q-date
        subtitle="End Date"
        v-model="endDate"
      />
    </div>
    <div class="q-pa-md q-gutter-sm">
      <q-table
        title="Quick Draw Data"
        :data="tableData"
        :columns="cols"
        :pagination.sync="pagination"
        :loading="loading"
        @request="getViewData"
        row-key="id"
        binary-state-sort
      >
      </q-table>
    </div>
  </q-page>
</template>

<script>
import { date } from 'quasar'

const tableCols = [
  {
    name: 'id',
    required: true,
    label: 'Draw #',
    aligh: 'left',
    field: row => row.id,
    sortable: true
  },
  {
    name: 'date',
    required: true,
    label: 'Date',
    aligh: 'right',
    align: 'center',
    field: row => row.drawDate,
    format: v => date.formatDate(v, 'MM-DD-YYYY'),
    sortable: false
  },
  {
    name: 'time',
    required: true,
    label: 'Time',
    aligh: 'center',
    field: row => row.drawTime,
    format: v => date.formatDate(v, 'hh:mm'),
    sortable: false
  },
  {
    name: 'picks',
    required: true,
    label: 'Winning Numbers',
    aligh: 'center',
    align: 'center',
    field: row => row.picks,
    format: val => `${val}`,
    sortable: false
  }
]

const apiURL = process.env.API || ''

export default {
  name: 'AllTimeData',
  data () {
    return {
      cols: tableCols,
      startDate: Date.now().toString(),
      endDate: Date.now().toString(),
      tableData: [],
      loading: true,
      pagination: {
        sortBy: 'id',
        descending: true,
        page: 1,
        rowsPerPage: 25,
        rowsNumber: 1
      }
    }
  },
  mounted () {
    this.getViewData({
      pagination: this.pagination,
      filter: undefined
    })
  },
  methods: {
    createDraw (draw) {
      return {
        id: draw.draw_number,
        drawDate: Date.parse(draw.draw_date),
        drawTime: Date.parse(draw.draw_time),
        picks: draw.winning_numbers,
        extra: draw.extra_multiplier,
        formatDate: this.formatDate.bind(this)
      }
    },
    formatDate (d) {
      return d.getMonth() + '-' + d.getDate() + '-' + d.getFullYear()
    },
    getURL (pagination) {
      const { page, rowsPerPage, sortBy, descending } = pagination
      const ascOrDesc = descending ? 'desc' : 'asc'
      return apiURL + '/draws/p' + page + '/s' + rowsPerPage + '/b' + sortBy + '/' + ascOrDesc
    },
    getViewData (props) {
      this.loading = true
      this.$axios.get(this.getURL(props.pagination))
        .then((response) => {
          const data = response.data || { draws: [] }
          const rows = this.lodash.map(data.draws, this.createDraw.bind(this))
          this.tableData.splice(0, this.tableData.length, ...rows)
          this.pagination.rowsNumber = data.total
          this.pagination.page = data.page
          this.pagination.rowsPerPage = data.page_size
          this.pagination.descending = props.pagination.descending
          this.pagination.sortBy = props.pagination.sortBy
          this.loading = false
        })
        .catch((e) => {
          this.loading = false
          this.$q.notify({
            color: 'negative',
            position: 'top',
            message: 'Loading Data Failed - ' + e,
            icon: 'report_problem'
          })
        })
    }
  }
}
</script>
