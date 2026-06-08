<template>
  <div class="max-w-4xl mx-auto space-y-5">
    <!-- Header -->
    <header class="flex justify-between items-center border-b border-slate-200 dark:border-zinc-800/80 pb-3">
      <div>
        <h1 class="text-xl sm:text-2xl font-bold text-slate-900 dark:text-white">{{ t('my_billing') }}</h1>
        <p class="text-[11px] text-slate-500 dark:text-zinc-400 mt-0.5">{{ t('billing_desc') || 'Manage and view your historical payment records and receipts' }}</p>
      </div>
    </header>

    <!-- Invoices List Card -->
    <section class="backdrop-blur-md bg-white/70 dark:bg-zinc-900/40 border border-slate-200 dark:border-zinc-800/80 p-4 rounded-lg space-y-4 shadow-sm">
      <h2 class="text-xs font-bold text-slate-900 dark:text-white flex items-center space-x-2">
        <UIcon name="i-lucide-receipt" class="w-4 h-4 text-primary-500 dark:text-primary-400" />
        <span>{{ t('invoice_history') || 'Invoice History' }}</span>
      </h2>

      <div class="overflow-x-auto">
        <table class="w-full text-left border-collapse text-xs">
          <thead>
            <tr class="border-b border-slate-200 dark:border-zinc-800 text-slate-500 dark:text-zinc-400 font-medium">
              <th class="py-2 px-3 w-28">Invoice ID</th>
              <th class="py-2 px-3">Plan</th>
              <th class="py-2 px-3 w-32">Date</th>
              <th class="py-2 px-3 w-24">Amount</th>
              <th class="py-2 px-3 w-28 text-center">Status</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="invoice in mockInvoices" :key="invoice.id" class="border-b border-slate-200/50 dark:border-zinc-800/40 hover:bg-slate-100/50 dark:hover:bg-zinc-900/20 text-slate-700 dark:text-zinc-200">
              <td class="py-3 px-3 font-mono font-bold text-slate-650 dark:text-zinc-400">{{ invoice.id }}</td>
              <td class="py-3 px-3 font-semibold">{{ invoice.plan }}</td>
              <td class="py-3 px-3 text-slate-550 dark:text-zinc-400 font-mono">{{ invoice.date }}</td>
              <td class="py-3 px-3 font-mono font-semibold text-slate-800 dark:text-zinc-200">{{ invoice.amount }}</td>
              <td class="py-3 px-3 text-center">
                <UBadge :color="invoice.statusColor" variant="subtle" size="xs">
                  {{ invoice.status }}
                </UBadge>
              </td>
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

useSeoMeta({
  title: () => `${t('my_billing')} - HY-Board`
})

const mockInvoices = [
  {
    id: 'INV-2026-003',
    plan: 'Pro Plan - Monthly Subscription',
    date: '2026-06-01',
    amount: '$7.99',
    status: 'Paid',
    statusColor: 'primary'
  },
  {
    id: 'INV-2026-002',
    plan: 'Starter Plan - Monthly Subscription',
    date: '2026-05-01',
    amount: '$3.99',
    status: 'Paid',
    statusColor: 'primary'
  },
  {
    id: 'INV-2026-001',
    plan: 'Starter Plan - Monthly Subscription',
    date: '2026-04-01',
    amount: '$3.99',
    status: 'Paid',
    statusColor: 'primary'
  }
]
</script>
