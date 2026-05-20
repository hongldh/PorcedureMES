<template>
  <DocPage title="销售出库单" :api="salesOrderApi" partyKey="customerId" partyLabel="客户"
    partyRel="customer" :partyApi="getInvCustomerAll" priceLabel="售价" confirmLabel="确认出库"
    defaultPriceField="salePrice" />
</template>
<script setup>
import DocPage from '../_DocPage.vue'
import { salesOrderApi } from '@/api/inventory/documents'
import { getInvCustomerAll } from '@/api/inventory/partner'
defineOptions({ name: 'InvSalesOrder' })
</script>
