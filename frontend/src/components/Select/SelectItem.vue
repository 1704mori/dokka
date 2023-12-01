<script setup lang="ts">
import {
  SelectItem,
  SelectItemIndicator,
  type SelectItemProps,
  SelectItemText,
} from "radix-vue";
import { Check } from "lucide-vue-next";
import { cn } from "@/lib/utils";

const props = defineProps<SelectItemProps & { class?: string }>();
</script>

<template>
  <SelectItem
    v-bind="props"
    :class="
      cn(
        'relative flex w-full cursor-pointer select-none items-center rounded-md py-1.5 pl-8 pr-2 text-sm outline-none hover:bg-neutral-100 focus:bg-neutral-100 focus:text-neutral-900 data-[disabled]:pointer-events-none data-[disabled]:opacity-50 dark:hover:bg-neutral-800 dark:focus:bg-neutral-800 dark:focus:text-neutral-50',
        props.class,
      )
    "
  >
    <span class="absolute left-2 flex h-3.5 w-3.5 items-center justify-center">
      <SelectItemIndicator>
        <Check class="h-4 w-4" />
      </SelectItemIndicator>
    </span>

    <SelectItemText>
      <slot />
    </SelectItemText>
  </SelectItem>
</template>
