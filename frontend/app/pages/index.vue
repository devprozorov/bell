<template>
  <section class="ui-panel p-10 mt-6">
    <div class="grid md:grid-cols-2 gap-10 items-center">
      <div>
        <h1 class="text-4xl md:text-5xl font-extrabold leading-tight">
          Видеоналитика задач и знаний<br>
          для вашей команды в стиле <span class="text-[color:var(--brand)]">BELL</span>
        </h1>
        <p class="mt-4 text-lg text-[color:var(--muted)]">
          Единая доска, вики и чат — в одном месте. Планируйте, фиксируйте знания и
          ускоряйте процессы с аккуратным UI и шустрым UX.
        </p>
        <div class="mt-6 flex gap-3">
          <NuxtLink to="/table" class="btn btn-primary">Перейти к доске</NuxtLink>
          <NuxtLink to="/wiki" class="btn btn-ghost">Открыть вики</NuxtLink>
        </div>
      </div>

      <div class="card p-6">
        <div class="aspect-video rounded-xl ui-panel flex items-center justify-center">
          <span class="text-sm text-[color:var(--muted)]">Здесь может быть промо-скрин/анимация</span>
        </div>
        <div class="mt-4 grid grid-cols-3 gap-3">
          <div class="ui-panel p-4 text-center">
            <div class="text-xl font-extrabold">3×</div>
            <div class="text-[color:var(--muted)] text-sm">быстрее поиск</div>
          </div>
          <div class="ui-panel p-4 text-center">
            <div class="text-xl font-extrabold">−45%</div>
            <div class="text-[color:var(--muted)] text-sm">времени на рутину</div>
          </div>
          <div class="ui-panel p-4 text-center">
            <div class="text-xl font-extrabold">+12%</div>
            <div class="text-[color:var(--muted)] text-sm">скорость релизов</div>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
definePageMeta({ ssr: false })
</script>
