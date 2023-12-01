<template>
  <div class="flex items-center flex-wrap justify-center m-auto py-1">
    <button
      v-if="buttons"
      class="flex items-center px-2 m-1 h-8 rounded-lg transition-colors bg-900 hover:bg-neutral-900 disabled:cursor-not-allowed disabled:bg-900/30 disabled:dark:bg-neutral-900 disabled:opacity-60"
      :disabled="currentPage === 1"
      @click="handleButtonClick(currentPage - 1)"
    >
      <ArrowLeft />
      <span class="text-sm font-medium">Back</span>
    </button>

    <button
      v-for="(number, index) in range(1, Math.ceil(Number(total) / perPage))"
      :key="index"
      class="flex items-center justify-center w-min min-w-[2rem] h-8 p-2 m-1 font-medium cursor-pointer rounded-lg transition-colors"
      :class="{ 'bg-primary text-white': number === currentPage, 'hover:bg-900': number !== currentPage }"
      @click="handleButtonClick(number)"
    >
      {{ number }}
    </button>

    <button
      v-if="buttons"
      class="flex items-center px-2 m-1 h-8 rounded-lg transition-colors bg-900 hover:bg-neutral-900 disabled:cursor-not-allowed disabled:bg-900/30 disabled:dark:bg-neutral-900 disabled:opacity-60"
      :disabled="currentPage === Math.ceil(Number(total) / perPage) || total === 0"
      @click="handleButtonClick(currentPage + 1)"
    >
      <span class="text-sm font-medium">Next</span>
      <ArrowRight />
    </button>
  </div>
</template>

<script lang="ts">
import { ref, watchEffect } from 'vue';
import { ArrowLeft, ArrowRight } from 'lucide-vue-next'; // Assuming lucide-react is available as a Vue component

export default {
  props: {
    total: Number,
    perPage: {
      type: Number,
      default: 15,
    },
    buttons: {
      type: Boolean,
      default: true,
    },
    maxButtons: {
      type: Number,
      default: 5,
    },
    onPageClick: Function,
  },
  setup(props) {
    const currentPage = ref(1);

 
    const range = (start: number, end: number) => {
      const middle = Math.ceil(props.maxButtons / 2);
      const isEven = props.maxButtons % 2 === 0;
      const isTooCloseToStart = currentPage.value <= middle;
      const isTooCloseToEnd = currentPage.value >= end - middle + 1;

      if (end <= props.maxButtons) {
        return Array.from({ length: end }, (_, i) => i + 1);
      }

      if (isTooCloseToStart) {
        return Array.from({ length: props.maxButtons }, (_, i) => i + 1);
      }

      if (isTooCloseToEnd) {
        return Array.from({ length: props.maxButtons }, (_, i) => end - props.maxButtons + i + 1);
      }

      return Array.from(
        { length: props.maxButtons },
        (_, i) => currentPage.value - middle + i + (isEven ? 1 : 0)
      );
    };

    const handleButtonClick = (pageNumber: number) => {
      currentPage.value = pageNumber;
      props.onPageClick && props.onPageClick(pageNumber);

      window.scrollTo({
        top: 0,
        behavior: 'smooth',
      });
    };

    return {
      currentPage,
      range,
      handleButtonClick,
    };
  },
};
</script>
