<template>
  <q-page class="flex flex-center">
    <q-date
      subtitle="Start Date"
      v-model="startDate"
    />
    <q-date
      subtitle="End Date"
      v-model="endDate"
    />
    <div class="q-pa-md q-gutter-sm">
      <q-table
        title="All-Time Data"
        :data="getViewData()"
        :columns="cols"
        row-key="id"
      >
      </q-table>
    </div>
  </q-page>
</template>

<script>
const tableCols = [
  {
    name: 'id',
    required: true,
    label: 'Draw #',
    aligh: 'left',
    field: row => row.id,
    format: val => `{val}`,
    sortable: true
  },
  {
    name: 'date',
    required: true,
    label: 'Date',
    aligh: 'left',
    field: row => row.drawDate,
    format: val => `{val}`,
    sortable: false
  },
  {
    name: 'time',
    required: true,
    label: 'Time',
    aligh: 'left',
    field: row => row.drawTime,
    format: val => `{val}`,
    sortable: false
  },
  {
    name: 'picks',
    required: true,
    label: 'Winning Numbers',
    aligh: 'left',
    field: row => row.picks,
    format: val => `{val}`,
    sortable: false
  }
]

export default {
  name: 'AllTimeData',
  data () {
    return {
      cols: tableCols,
      startDate: Date.now().toString(),
      endDate: Date.now().toString()
    }
  },
  methods: {
    getViewData () {
      return this.$axios.get('http://localhost:9090/draws/' + this.$route.params.date)
        .then((response) => {
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
  }
}
</script>
