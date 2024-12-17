<script lang="ts">
    import {PickDBFile, PickLogFile, RunImport} from '../../wailsjs/go/main/App.js'
    import {resultcollector} from '../../wailsjs/go/models'
    import Loader from "./Loader.svelte";

    let results: resultcollector.Result[] = []

    export let dbFile = ""
    let logFile = ""

    let loading = false

    function pickLog(): void {
        PickLogFile().then(result => logFile = result)
    }

    function runImport() {
        loading = true

        RunImport(logFile, dbFile).then((incomingResults) => {
            results = incomingResults
            loading = false
        }).catch(error => {
            alert(error)
            loading = false
        })
    }
</script>

<div class="params">
    <div class="param">
        <span>Log File</span>
        <span class="main">{logFile}</span>
        <button on:click={() => pickLog()} disabled={loading}>Select Log File</button>
    </div>
    <div class="param">
        <span>Database File</span>
        <span class="main">{dbFile}</span>
        <span class="placeholder"></span>
    </div>
    <button on:click={() => runImport()} disabled={dbFile === "" || logFile === "" || loading}>Run Importer</button>
</div>

{#if loading}
    <Loader/>
{/if}

<pre>
{results.map((result) => {
        return `${result.filename ? (result.filename + ":" + result.line) : "Generic"} - ${result.success ? " OK " : "FAIL"} - ${result.message}`
    }).join("\n")}
</pre>

<style>
    pre {
        width: calc(100% - 20px);
        max-height: 300px;
        overflow: auto;
        text-align: left;
        margin: 10px;
    }


    .param {
        display: flex;
    }

    .param span {
        width: 130px;
        text-align: left;
        padding-right: 10px;
    }
    .param span:nth-child(1) {
        text-align: right;
    }

    .param button, .param .placeholder {
        width: 130px;
        text-align: center;
        padding-right: 10px;
    }

    .param .main {
        flex-grow: 1;
    }

    .params {
        padding: 10px;
    }

    .params > button {
        margin-top: 20px;
        margin-bottom: 0;
    }
</style>
