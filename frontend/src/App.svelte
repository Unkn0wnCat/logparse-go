<script lang="ts">
  import {PickDBFile} from '../wailsjs/go/main/App.js'
  import Importer from "./components/Importer.svelte";
  import Viewer from "./components/Viewer.svelte";
  import IPCounter from "./components/IPCounter.svelte";
  import StatusCounter from "./components/StatusCounter.svelte";

  let dbFile: string = ""

  let activeTab: "importer"|"view"|"ip_counter"|"error_counter" = "importer"

  function pickDb(): void {
    PickDBFile().then(result => dbFile = result)
  }
</script>

<main>
  <h1>Log Parser</h1>

  <div class="dbselect">
    <span class="title">Datenbank</span>
    <span class="path">{dbFile || "Keine Datenbank geöffnet"}</span>
    <button on:click={() => pickDb()}>Auswählen</button>
  </div>

  <div class="tabs">
    <div class="tabbar">
      <button on:click={() => activeTab = "importer"} class:active={activeTab === "importer"}>Import</button>
      <button on:click={() => activeTab = "view"} class:active={activeTab === "view"}>Anzeigen</button>
      <button on:click={() => activeTab = "ip_counter"} class:active={activeTab === "ip_counter"}>IP-Analyse</button>
      <button on:click={() => activeTab = "error_counter"} class:active={activeTab === "error_counter"}>Status-Analyse</button>
    </div>
    {#if activeTab === "importer"}
      <div class="tab">
        <Importer {dbFile} />
      </div>
    {/if}
    {#if activeTab === "view"}
      <div class="tab">
        <Viewer {dbFile} />
      </div>
    {/if}
    {#if activeTab === "ip_counter"}
      <div class="tab">
        <IPCounter {dbFile} />
      </div>
    {/if}
    {#if activeTab === "error_counter"}
      <div class="tab">
        <StatusCounter {dbFile} />
      </div>
    {/if}
  </div>

</main>

<style>
  .dbselect {
    display: flex;
    flex-direction: column;
    border: 2px solid white;
    width: 300px;
    margin: 10px auto;
    border-radius: 5px;
  }

  .dbselect .title {
    font-weight: bold;
    font-size: 1.2em;
  }

  .dbselect .path {
    font-family: monospace;
    padding: 10px;
  }

  .dbselect button {
    border: none;
  }

  .tabs {
    display: flex;
    flex-direction: column;
    border: 2px solid white;
    border-radius: 5px;
    margin: 20px;
  }

  .tabbar {
    border-bottom: 2px solid white;
  }

  .tabbar button {
    border: none;
    border-right: 2px solid white;
    border-bottom: 4px solid transparent;
    font: inherit;
    color: inherit;
    background-color: transparent;
    margin: 0;
    padding: 10px;
    border-radius: 0;
    appearance: none;
  }

  .tabbar button.active {
    border-bottom-color: orange;
  }

  .tabbar button:last-child {
    border-right: none;
  }

  .input-box .btn {
    width: 60px;
    height: 30px;
    line-height: 30px;
    border-radius: 3px;
    border: none;
    margin: 0 0 0 20px;
    padding: 0 8px;
    cursor: pointer;
  }

  .input-box .btn:hover {
    background-image: linear-gradient(to top, #cfd9df 0%, #e2ebf0 100%);
    color: #333333;
  }

  .input-box .input {
    border: none;
    border-radius: 3px;
    outline: none;
    height: 30px;
    line-height: 30px;
    padding: 0 10px;
    background-color: rgba(240, 240, 240, 1);
    -webkit-font-smoothing: antialiased;
  }

  .input-box .input:hover {
    border: none;
    background-color: rgba(255, 255, 255, 1);
  }

  .input-box .input:focus {
    border: none;
    background-color: rgba(255, 255, 255, 1);
  }

</style>
