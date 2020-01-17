<template>
  <q-layout view="lHh Lpr lFf">
    <q-header elevated>
      <q-toolbar>
        <q-btn
          flat
          dense
          round
          @click="leftDrawerOpen = !leftDrawerOpen"
          icon="menu"
          aria-label="Menu"
        />
        <q-toolbar-title>Quick Draw Explorer</q-toolbar-title>
      </q-toolbar>
    </q-header>

    <q-drawer
      v-model="leftDrawerOpen"
      show-if-above
      bordered
      content-class="bg-grey-2"
    >
      <q-scroll-area class="fit">
        <q-list
          v-for="(menuItem, index) in menuList"
          :key="index"
        >
          <q-item
            clickable
            exact
            :to="menuItem.link"
            :active="menuItem.label === 'Outbox'"
            v-ripple
          >
            <q-item-section avatar>
              <q-icon :name="menuItem.icon" />
            </q-item-section>
            <q-item-section>
              {{ menuItem.label }}
            </q-item-section>
          </q-item>
          <q-separator v-if="menuItem.separator" />
        </q-list>
      </q-scroll-area>
    </q-drawer>
    <q-page-container>
      <router-view />
    </q-page-container>
  </q-layout>
</template>

<script>
const menuList = [
  {
    icon: 'trending_up',
    label: '24 Hours',
    link: '/totals/oneday',
    separator: true
  },
  {
    icon: 'timeline',
    label: 'One Week',
    link: '/totals/oneweek',
    separator: true
  },
  {
    icon: 'timeline',
    label: 'One Month',
    link: '/totals/onemonth',
    separator: true
  },
  {
    icon: 'timeline',
    label: 'All-Time',
    separator: true
  },
  {
    icon: 'view_list',
    label: 'Table Explorer',
    link: '/query',
    separator: true
  }
]
export default {
  name: 'MyLayout',

  data () {
    return {
      leftDrawerOpen: false,
      menuList: menuList
    }
  }
}
</script>
