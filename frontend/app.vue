<template>
  <UApp>
    <div class="min-h-screen bg-white dark:bg-zinc-950 text-slate-900 dark:text-zinc-100 font-sans selection:bg-primary-500/30 selection:text-primary-400 relative overflow-hidden transition-colors duration-300">
      <!-- Animated background glowing orbs for premium visual aesthetic -->
      <div class="absolute top-0 left-1/4 w-[500px] h-[500px] bg-primary-500/5 dark:bg-primary-500/10 rounded-full blur-[120px] pointer-events-none -translate-y-1/2"></div>
      <div class="absolute bottom-0 right-1/4 w-[600px] h-[600px] bg-primary-500/5 dark:bg-primary-500/10 rounded-full blur-[150px] pointer-events-none translate-y-1/2"></div>
      
      <div class="relative z-10">
        <NuxtLayout>
          <NuxtPage />
        </NuxtLayout>
      </div>
    </div>
  </UApp>
</template>

<style>
/* Non-linear Page transitions styling */
.page-enter-active,
.page-leave-active {
  transition: all 0.45s cubic-bezier(0.25, 1, 0.5, 1);
}
.page-enter-from {
  opacity: 0;
  transform: translateY(12px) scale(0.99);
}
.page-leave-to {
  opacity: 0;
  transform: translateY(-12px) scale(0.99);
}

/* Base custom scrollbar styling */
::-webkit-scrollbar {
  width: 6px;
  height: 6px;
}
::-webkit-scrollbar-track {
  background: transparent;
}
::-webkit-scrollbar-thumb {
  background: rgba(79, 125, 166, 0.2);
  border-radius: 4px;
}
::-webkit-scrollbar-thumb:hover {
  background: rgba(79, 125, 166, 0.4);
}
</style>
