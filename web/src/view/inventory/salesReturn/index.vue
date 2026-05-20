<template>
  <DocPage title="销售退货单" :api="salesReturnApi" partyKey="customerId" partyLabel="客户"
    partyRel="customer" :partyApi="getInvCustomerAll" priceLabel="退货单价" confirmLabel="确认退入"
    defaultPriceField="salePrice" />
</template>
<script setup>
import DocPage from '../_DocPage.vue'
import { salesReturnApi } from '@/api/inventory/documents'
import { getInvCustomerAll } from '@/api/inventory/partner'
defineOptions({ name: 'InvSalesReturn' })
</script>
