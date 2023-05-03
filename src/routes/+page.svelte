<script>
  import { onMount } from "svelte";

  let time = "";
  onMount(() => {
    const evtSrc = new EventSource("//localhost:8080/event");
    evtSrc.onmessage = (e) => {
      time = e.data;
    };
    evtSrc.onerror = (e) => {
      console.log("EventSource failed:", e);
    };
  });

  async function getTime() {
    const res = await fetch("//localhost:8080/time");
    // time = await res.text();
    if (res.status !== 200) {
      throw new Error(res.statusText);
    }
  }
</script>

<main>
  <h1>Server Sent Events</h1>
  <button on:click={getTime}>Get time</button>
  <p>Time: {time}</p>
</main>
