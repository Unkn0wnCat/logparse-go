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

<div>
    <div>
        <span>Log File</span>
        <span>{logFile}</span>
        <button on:click={() => pickLog()} disabled={loading}>Select Log File</button>
    </div>
    <div>
        <span>Database File</span>
        <span>{dbFile}</span>
    </div>
    <button on:click={() => runImport()} disabled={dbFile === "" || logFile === "" || loading}>Run Importer</button>
</div>

{#if loading}
    <Loader/>
{/if}

<pre>
    {results.map((result) => {
        return `${result.filename}:${result.line} - ${result.success ? " OK " : "FAIL"} - ${result.message}`
    }).join("\n")}
</pre>
