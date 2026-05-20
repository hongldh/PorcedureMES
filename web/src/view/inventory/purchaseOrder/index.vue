<template>
  <DocPage title="采购入库单" :api="purchaseOrderApi" partyKey="supplierId" partyLabel="供应商"
    partyRel="supplier" :partyApi="getSupplierAll" priceLabel="单价" confirmLabel="确认入库"
    defaultPriceField="purchasePrice" />
</template>
<script setup>
import DocPage from '../_DocPage.vue'
import { purchaseOrderApi } from '@/api/inventory/documents'
import { getSupplierAll } from '@/api/inventory/partner'
defineOptions({ name: 'InvPurchaseOrder' })
</script>
