const text = document.querySelector("script#config__json")?.textContent || "{}";

export interface Config {
  port: number;
  event_stream_interval: number;
}

const pageConfig = JSON.parse(text);

const config: Config = {
  ...pageConfig.env,
};

export const host = window.location.origin;
export default Object.freeze(config);
