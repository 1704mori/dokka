<template>
  <Table
    :data="data"
    :columns="columns"
    :rowsPerPage="10"
    :totalData="total"
    @pageChanged="handlePageChange"
  />
</template>

<script lang="ts">
import { defineComponent, ref, watchEffect } from "vue";
import Table from "@/components/Table.vue";

// @ts-ignore
import type { Container } from "@/lib/types/container";

export default defineComponent({
  name: "HomeView",
  components: {
    Table,
  },
  setup() {
    const eventSource = new EventSource("http://localhost:7070/v1/containers");

    const columns = [
      {
        key: "id",
        label: "ID",
      },
      {
        key: "name",
        label: "Name",
      },
      {
        key: "status",
        label: "Status",
      },
      {
        key: "image",
        label: "Image",
      },
      {
        key: "ports",
        label: "Ports",
      },
      {
        key: "created_at",
        label: "Created At",
      },
    ];

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
      total.value = containerData.length;
    };

    const handlePageChange = (page: number) => {
      console.log(page);
    };

    const transformData = (data: Container[]) =>
      data.map((item) => {
        return {
          id: item.Id.slice(0, 12),
          name: item.Names[0],
          status: item.State,
          image: item.Image,
          ports:
            item.Ports.length > 0
              ? item.Ports.map(
                  (port) => `${port.PublicPort}:${port.PrivatePort}`
                ).join(", ")
              : "N/A",
          created_at: new Date(item.Created * 1000).toLocaleDateString(
            undefined,
            {
              year: "numeric",
              month: "long",
              day: "numeric",
            }
          ),
        };
      });

    return {
      columns,
      data,
      total,
      handlePageChange,
    };
  },
});
</script>
