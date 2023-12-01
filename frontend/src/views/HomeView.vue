<template>
  <div class="grid grid-cols-[15rem_1fr] gap-3">
    <div
      class="flex flex-col items-center gap-3 bg-white dark:bg-black border border-neutral-200 dark:border-neutral-900 px-2 py-3 rounded-md h-fit"
    >
      <span class="text-lg font-medium">Configuration</span>
      <div class="flex flex-col gap-2 font-medium w-full">
        <button
          class="bg-neutral-50 dark:bg-neutral-950 border border-neutral-200 dark:border-neutral-900 hover:bg-neutral-100 dark:hover:bg-neutral-900 px-3.5 py-1.5 rounded-md"
        >
          Table
        </button>
        <button
          class="bg-neutral-50 dark:bg-neutral-950 border border-neutral-200 dark:border-neutral-900 hover:bg-neutral-100 dark:hover:bg-neutral-900 px-3.5 py-1.5 rounded-md"
        >
          Authentication
        </button>
        <button
          class="bg-neutral-50 dark:bg-neutral-950 border border-neutral-200 dark:border-neutral-900 hover:bg-neutral-100 dark:hover:bg-neutral-900 px-3.5 py-1.5 rounded-md"
        >
          API
        </button>
      </div>
    </div>
    <div class="flex flex-col gap-2">
      <div class="grid grid-cols-2 gap-1">
        <Input placeholder="Search..." type="text" v-model="search" />
        <Select v-model="status">
          <SelectTrigger>
            <SelectValue placeholder="Status" />
          </SelectTrigger>
          <SelectContent>
            <SelectGroup>
              <SelectItem value="running">Running</SelectItem>
              <SelectItem value="paused">Paused</SelectItem>
              <SelectItem value="restarting">Restarting</SelectItem>
              <SelectItem value="exited">Exited</SelectItem>
              <SelectItem value="removing">Removing</SelectItem>
              <SelectItem value="created">Created</SelectItem>
              <SelectItem value="dead">Dead</SelectItem>
            </SelectGroup>
          </SelectContent>
        </Select>
      </div>

      <Table
        :data="data"
        :columns="columns"
        :rowsPerPage="10"
        :totalData="total"
        @pageChanged="handlePageChange"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, h, watchEffect } from "vue";
import Table from "@/components/Table.vue";
import Input from "@/components/Input.vue";
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectTrigger,
  SelectValue,
  // @ts-ignore
} from "@/components/Select";
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
  // @ts-ignore
} from "@/components/Tooltip";

import type { Container } from "@/lib/types/container";
import { StopCircle } from "lucide-vue-next";
import { RotateCcw } from "lucide-vue-next";

const eventSource = new EventSource("http://localhost:7070/v1/containers");

const columns = ref([
  { key: "id", label: "ID" },
  { key: "name", label: "Name" },
  { key: "status", label: "Status" },
  { key: "image", label: "Image" },
  { key: "ports", label: "Ports" },
  { key: "created_at", label: "Created At" },
  { key: "actions", label: "Actions" },
]);

const status = ref<string>("");
const search = ref<string>("");

const data = ref<
  {
    id: string;
    name: string;
    status: string;
    image: string;
    ports: string;
    created_at: string;
  }[]
>([]);
const ogData = ref<
  {
    id: string;
    name: string;
    status: string;
    image: string;
    ports: string;
    created_at: string;
  }[]
>([]);

const total = ref(0);

eventSource.onmessage = (event) => {
  const [type] = event.data.split(":");
  const [, _data] = event.data.split(`${type}:`);

  if (type == "error:") {
    console.log(_data);
    return;
  }

  const containerData: Container[] = JSON.parse(_data);
  data.value = transformData(containerData);
  ogData.value = transformData(containerData);
  total.value = containerData.length;
};

const handlePageChange = (page: number) => {
  console.log(page);
};

const handleStopContainer = async (containerId: string) => {
  const res = await fetch(
    `http://localhost:7070/v1/container/stop/${containerId}`,
    {
      method: "POST",
    },
  );
  const json = await res.json();
  console.log(json);
};

watchEffect(() => {
  if (status.value && status.value != "") {
    switch (status.value) {
      case "running":
        data.value = data.value.filter((f) => /running/.test(f.status));
        break;
      case "paused":
        data.value = data.value.filter((f) => /paused/.test(f.status));
        break;
      case "exited":
        data.value = data.value.filter((f) => /exited/.test(f.status));
        break;

      default:
        break;
    }
  }
});

watchEffect(() => {
  if (search.value == "") {
    console.log("search input is empty, resetting the data array");
    data.value = ogData.value;
  } else {
    data.value = data.value.filter((f) =>
      new RegExp(search.value, "i").test(f.name),
    );
  }
})

const transformData = (data: Container[]) =>
  data.map((item) => ({
    id: `<p class="truncate max-w-[8rem]">${item.Id}</p>`,
    name: item.Names[0],
    status: `
      <div class="flex items-center gap-1">
        <div class="w-2 h-2 rounded-full ${
          item.State === "running"
            ? "bg-green-500"
            : item.State === "paused"
              ? "bg-yellow-500"
              : item.State === "exited" ||
                  item.State === "dead" ||
                  item.State === "removing"
                ? "bg-red-500"
                : "bg-gray-500"
        }"></div>
        <span>${item.State}</span>
      </div>`,
    image: item.Image,
    ports:
      item.Ports.length > 0
        ? item.Ports.map(
            (port) => `${port.PublicPort}:${port.PrivatePort}`,
          ).join(", ")
        : "N/A",
    created_at: new Date(item.Created * 1000).toLocaleDateString(undefined, {
      year: "numeric",
      month: "long",
      day: "numeric",
    }),
    actions: () => {
      return h("div", { class: "flex justify-center items-center gap-2 " }, [
        h(TooltipProvider, [
          h(Tooltip, [
            h(TooltipTrigger, { asChild: true }, [
              h(
                "button",
                {
                  onClick: () => handleStopContainer(item.Id),
                },
                [h(StopCircle, { size: 20 })],
              ),
            ]),
            h(TooltipContent, ["Stop"]),
          ]),
        ]),
        h(TooltipProvider, [
          h(Tooltip, [
            h(TooltipTrigger, { asChild: true }, [
              h("button", [h(RotateCcw, { size: 20 })]),
            ]),
            h(TooltipContent, ["Restart"]),
          ]),
        ]),
      ]);
    },
  }));
</script>
