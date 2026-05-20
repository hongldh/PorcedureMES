<template>
  <DocPage title="采购退货单" :api="purchaseReturnApi" partyKey="supplierId" partyLabel="供应商"
    partyRel="supplier" :partyApi="getSupplierAll" priceLabel="退货单价" confirmLabel="确认退货"
    defaultPriceField="purchasePrice" />
</template>
<script setup>
import DocPage from '../_DocPage.vue'
import { purchaseReturnApi } from '@/api/inventory/documents'
import { getSupplierAll } from '@/api/inventory/partner'
defineOptions({ name: 'InvPurchaseReturn' })
</script>
