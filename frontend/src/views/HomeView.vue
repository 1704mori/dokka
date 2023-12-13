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
          General
        </button>
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
      </div>
    </div>
    <div class="flex flex-col gap-2">
      <div class="grid grid-cols-2 gap-1">
        <Input
          placeholder="Search..."
          type="text"
          v-model="search as string"
          @input="handleSearch"
        />
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
import { ref, h, watchEffect, onMounted } from "vue";
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

import { useContainer } from "@/lib/useContainer";
import { host } from "@/lib/env";
import { PlayCircle } from "lucide-vue-next";
import { toast } from "vue-sonner";

let eventSource: EventSource;

function startEventSource() {
  eventSource = new EventSource(`${host}/v1/containers`);
  eventSource.onopen = () => {
    console.log("event source opened");
  };

  eventSource.onerror = (error) => {
    console.log("event source error", error);
  };
}

function restartEventSource() {
  eventSource.close();
  startEventSource();
}

onMounted(() => {
  startEventSource();

  eventSource.onmessage = (event) => {
    console.log("event", event);
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
});

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
const search = ref<string | null>(null);

const loading = ref<{
  [containerId: string]: boolean;
}>({});

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

const handlePageChange = (page: number) => {
  console.log(page);
};

const handleSearch = (e: any) => {
  if (e.target.value.trim().length === 0) {
    data.value = ogData.value;
  }
};

const handleStopContainer = async (containerId: string) => {
  // loading.value[containerId] = true;
  // const { data, error } = await useContainer<{ message: string }>(
  //   `/v1/container/stop/${containerId}`,
  //   {
  //     method: "POST",
  //   },
  // );

  // if (error) {
  //   toast.error(error, { duration: 5000 });
  //   loading.value[containerId] = false;
  //   return;
  // }

  // toast.success(data?.message || "Container stopped successfully");
  // loading.value[containerId] = false;
  toast.promise(
    useContainer<{ message: string }>(`/v1/container/stop/${containerId}`, {
      method: "POST",
    }),
    {
      loading: "Stopping container...",
      success: ({ data }) => {
        return data?.message || "Container stopped successfully";
      },
      error: ({ error }) => error || "Something went wrong",
      duration: 10000, // the duration is not working the promises so the value is set to 10000
    }
  );
};

const handleStartContainer = async (containerId: string) => {
  // loading.value[containerId] = true;
  // const { data, error } = await useContainer<{ message: string }>(
  //   `/v1/container/start/${containerId}`,
  //   {
  //     method: "POST",
  //   },
  // );

  // if (error) {
  //   toast.error(error, { duration: 5000 });
  //   loading.value[containerId] = false;
  //   return;
  // }

  // toast.success(data?.message || "Container started successfully");
  // loading.value[containerId] = false;
  toast.promise(
    useContainer<{ message: string }>(`/v1/container/start/${containerId}`, {
      method: "POST",
    }),
    {
      loading: "Starting container...",
      success: ({ data }) => {
        return data?.message || "Container started successfully";
      },
      error: ({ error }) => error || "Something went wrong",
    }
  );
};

const handleRestartContainer = async (containerId: string) => {
  // loading.value[containerId] = true;
  // const { data, error } = await useContainer<{ message: string }>(
  //   `/v1/container/restart/${containerId}`,
  //   {
  //     method: "POST",
  //   },
  // );

  // if (error) {
  //   toast.error(error, { duration: 5000 });
  //   loading.value[containerId] = false;
  //   return;
  // }

  // toast.success(data?.message || "Container restarted successfully");
  // loading.value[containerId] = false;
  toast.promise(
    useContainer<{ message: string }>(`/v1/container/restart/${containerId}`, {
      method: "POST",
    }),
    {
      loading: "Restarting container...",
      success: ({ data }) => {
        return data?.message || "Container restarted successfully";
      },
      error: ({ error }) => error || "Something went wrong",
    }
  );
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
  if (search.value && search.value != "") {
    data.value = data.value.filter((f) =>
      new RegExp(search.value as string, "i").test(f.name)
    );
  }
});

const transformData = (data: Container[]) =>
  data.map((item) => ({
    id: `<p class="truncate max-w-[8rem]">${item.Id}</p>`,
    name: item.Names[0],
    status: `
      <div class="flex items-center gap-1">
        <div class="w-2 h-2 rounded-full ${
          !loading.value[item.Id] && item.State === "running"
            ? "bg-green-500"
            : item.State === "paused"
            ? "bg-yellow-500"
            : item.State === "exited" ||
              item.State === "dead" ||
              item.State === "removing"
            ? "bg-red-500"
            : "bg-gray-500"
        }"></div>
        ${
          !loading.value[item.Id]
            ? `<span>${item.State}</span>`
            : `<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" class="animate-spin" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="lucide lucide-loader-2"><path d="M21 12a9 9 0 1 1-6.219-8.56"/></svg>`
        }
      </div>`,
    image: item.Image,
    ports:
      item.Ports.length > 0
        ? item.Ports.map(
            (port) => `${port.PublicPort}:${port.PrivatePort}`
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
                  onClick: () =>
                    item.State == "running"
                      ? handleStopContainer(item.Id)
                      : handleStartContainer(item.Id),
                },
                [
                  h(item.State == "running" ? StopCircle : PlayCircle, {
                    size: 20,
                  }),
                ]
              ),
            ]),
            h(TooltipContent, [item.State == "running" ? "Stop" : "Start"]),
          ]),
        ]),
        h(TooltipProvider, [
          h(Tooltip, [
            h(TooltipTrigger, { asChild: true }, [
              h(
                "button",
                {
                  onClick: () => handleRestartContainer(item.Id),
                },
                [h(RotateCcw, { size: 20 })]
              ),
            ]),
            h(TooltipContent, ["Restart"]),
          ]),
        ]),
      ]);
    },
  }));
</script>
