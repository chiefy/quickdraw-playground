
const viewTypeTitles = {
  'oneday': 'Last 24 Hours',
  'oneweek': 'Last Seven Days',
  'onemonth': 'Last Month',
  'alltime': 'All-Time'
}

const routes = [
  {
    path: '/',
    component: () => import('layouts/MainLayout.vue'),
    children: [
      {
        path: 'totals/:viewType?',
        meta: {
          viewTypeTitles
        },
        component: () => import('pages/Index.vue')
      },
      { path: 'query/', component: () => import('pages/Table.vue') }
    ]
  }
]

// Always leave this as last one
if (process.env.MODE !== 'ssr') {
  routes.push({
    path: '*',
    component: () => import('pages/Error404.vue')
  })
}

export default routes
