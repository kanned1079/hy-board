<template>
  <div class="max-w-4xl mx-auto space-y-5">
    <!-- Header -->
    <header class="flex justify-between items-center border-b border-slate-200 dark:border-zinc-800/80 pb-3">
      <div>
        <h1 class="text-xl sm:text-2xl font-bold text-slate-900 dark:text-white">{{ t('my_invite') }}</h1>
        <p class="text-[11px] text-slate-500 dark:text-zinc-400 mt-0.5">{{ t('invite_desc') || 'Invite friends to join and earn recurrent commission rewards on their plan purchases' }}</p>
      </div>
    </header>

    <!-- Stats Cards -->
    <div class="grid grid-cols-1 sm:grid-cols-3 gap-4">
      <div class="backdrop-blur-md bg-white/70 dark:bg-zinc-900/40 border border-slate-200 dark:border-zinc-800/80 p-4 rounded-lg space-y-2 shadow-sm">
        <span class="text-[9px] font-semibold text-slate-450 dark:text-zinc-500 uppercase tracking-wider">{{ t('total_referrals') || 'Total Referrals' }}</span>
        <div class="text-xl font-extrabold text-slate-900 dark:text-white font-mono">4</div>
      </div>
      <div class="backdrop-blur-md bg-white/70 dark:bg-zinc-900/40 border border-slate-200 dark:border-zinc-800/80 p-4 rounded-lg space-y-2 shadow-sm">
        <span class="text-[9px] font-semibold text-slate-450 dark:text-zinc-500 uppercase tracking-wider">{{ t('total_commissions') || 'Total Commissions' }}</span>
        <div class="text-xl font-extrabold text-primary-600 dark:text-primary-400 font-mono">$12.50</div>
      </div>
      <div class="backdrop-blur-md bg-white/70 dark:bg-zinc-900/40 border border-slate-200 dark:border-zinc-800/80 p-4 rounded-lg space-y-2 shadow-sm">
        <span class="text-[9px] font-semibold text-slate-450 dark:text-zinc-500 uppercase tracking-wider">{{ t('unpaid_balance') || 'Unpaid Balance' }}</span>
        <div class="text-xl font-extrabold text-slate-950 dark:text-white font-mono">$4.50</div>
      </div>
    </div>

    <!-- Referral Link -->
    <section class="backdrop-blur-md bg-white/70 dark:bg-zinc-900/40 border border-slate-200 dark:border-zinc-800/80 p-5 rounded-lg space-y-3 shadow-sm">
      <h3 class="text-xs font-bold text-slate-900 dark:text-white">{{ t('referral_link') || 'Your Referral Invitation Link' }}</h3>
      <p class="text-[10px] text-slate-550 dark:text-zinc-400">Share this link to earn a 20% commission on every payment made by invited friends.</p>
      
      <div class="flex flex-col sm:flex-row gap-3 pt-1">
        <UInput
          :model-value="referralUrl"
          readonly
          class="flex-1 font-mono text-xs"
          color="primary"
          size="md"
          icon="i-lucide-link"
        />
        <UButton
          color="primary"
          size="md"
          icon="i-lucide-copy"
          @click="copyReferralLink"
        >
          {{ t('copy_link') }}
        </UButton>
      </div>
    </section>

    <!-- Referred Users list -->
    <section class="backdrop-blur-md bg-white/70 dark:bg-zinc-900/40 border border-slate-200 dark:border-zinc-800/80 p-4 rounded-lg space-y-4 shadow-sm">
      <h3 class="text-xs font-bold text-slate-900 dark:text-white flex items-center space-x-1.5">
        <UIcon name="i-lucide-users" class="w-4 h-4 text-primary-500" />
        <span>{{ t('referred_users_list') || 'Referred Users' }}</span>
      </h3>

      <div class="overflow-x-auto">
        <table class="w-full text-left border-collapse text-xs">
          <thead>
            <tr class="border-b border-slate-200 dark:border-zinc-800 text-slate-500 dark:text-zinc-400 font-medium">
              <th class="py-2 px-3">Email</th>
              <th class="py-2 px-3 w-32">Join Date</th>
              <th class="py-2 px-3 w-32 text-right">Commission Reward</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="user in mockReferrals" :key="user.email" class="border-b border-slate-200/50 dark:border-zinc-800/40 hover:bg-slate-100/50 dark:hover:bg-zinc-900/20 text-slate-700 dark:text-zinc-200">
              <td class="py-2.5 px-3">{{ user.email }}</td>
              <td class="py-2.5 px-3 text-slate-500 font-mono">{{ user.date }}</td>
              <td class="py-2.5 px-3 text-right font-mono font-bold text-primary-600 dark:text-primary-400">{{ user.reward }}</td>
            </tr>
          </tbody>
        </table>
      </div>
    </section>
  </div>
</template>

<script setup>
definePageMeta({
  layout: 'user'
})

const { t } = useI18n()
const toast = useToast()

const referralUrl = computed(() => {
  if (typeof window !== 'undefined') {
    return `${window.location.origin}/register?aff=10324`
  }
  return '/register?aff=10324'
})

useSeoMeta({
  title: () => `${t('my_invite')} - HY-Board`
})

const copyReferralLink = () => {
  navigator.clipboard.writeText(referralUrl.value)
  toast.add({
    id: 'referral_copied',
    title: t('copied') || 'Copied',
    description: 'Referral link has been copied to clipboard.',
    color: 'primary',
    timeout: 2000
  })
}

const mockReferrals = [
  { email: 'alex****@gmail.com', date: '2026-05-15', reward: '+$3.19' },
  { email: 'jame****@outlook.com', date: '2026-05-02', reward: '+$1.59' },
  { email: 'miki****@live.jp', date: '2026-04-22', reward: '+$6.12' },
  { email: 'huan****@gmail.com', date: '2026-04-18', reward: '+$1.60' }
]
</script>
