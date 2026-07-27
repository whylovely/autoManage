<script lang="ts" setup>
import { onBeforeUnmount, ref } from 'vue'
import HelloWorld from './components/HelloWorld.vue'
import { EventsOn } from '../wailsjs/runtime/runtime'

type DueReminder = {
  reminder: {
    Title: string
  }
  dueByDate: boolean
  dueByOdometer: boolean
}

const dueReminder = ref<DueReminder | null>(null)

const unsubscribe = EventsOn('reminder:due', (payload: DueReminder) => {
  dueReminder.value = payload
  console.info('Due reminder:', payload)
})

onBeforeUnmount(() => {
  unsubscribe()
})
</script>

<template>
  <aside v-if="dueReminder" class="reminder-toast" role="status">
    <strong>Напоминание: {{ dueReminder.reminder.Title }}</strong>
    <span v-if="dueReminder.dueByDate">Срок по дате скоро наступит.</span>
    <span v-if="dueReminder.dueByOdometer">Срок по пробегу скоро наступит.</span>
    <button type="button" aria-label="Закрыть уведомление" @click="dueReminder = null">×</button>
  </aside>
  <img id="logo" alt="Wails logo" src="./assets/images/logo-universal.png"/>
  <HelloWorld/>
</template>

<style>
#logo {
  display: block;
  width: 50%;
  height: 50%;
  margin: auto;
  padding: 10% 0 0;
  background-position: center;
  background-repeat: no-repeat;
  background-size: 100% 100%;
  background-origin: content-box;
}

.reminder-toast {
  position: fixed;
  top: 1rem;
  right: 1rem;
  z-index: 1;
  display: grid;
  gap: 0.35rem;
  max-width: 20rem;
  padding: 1rem 2.5rem 1rem 1rem;
  border-radius: 0.5rem;
  background: #fff4cc;
  color: #362d00;
  text-align: left;
  box-shadow: 0 0.5rem 1.5rem rgb(0 0 0 / 25%);
}

.reminder-toast button {
  position: absolute;
  top: 0.4rem;
  right: 0.55rem;
  border: 0;
  background: transparent;
  color: inherit;
  font-size: 1.4rem;
  cursor: pointer;
}
</style>
