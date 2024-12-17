<script lang="ts">
    import {QueryView} from '../../wailsjs/go/main/App.js'
    import {parser} from '../../wailsjs/go/models'
    import FilterPicker from "./FilterPicker.svelte";
    import Loader from "./Loader.svelte";

    let results: parser.LogLine[] = []

    export let dbFile = ""

    let timestampFrom = ""
    let timestampTo = ""
    let ip = ""

    let loading = false

    function query() {
        loading = true
        results = []

        QueryView(dbFile, ip, timestampFrom, timestampTo).then((incomingResults) => {
            results = incomingResults || []
            console.log(incomingResults)
            loading = false
        }).catch(error => {
            alert(error)
            loading = false
        })
    }
</script>

<div>
    <FilterPicker bind:timestampFrom={timestampFrom} bind:timestampTo={timestampTo} bind:ipAddress={ip}/>
    <button on:click={() => query()} disabled={dbFile === "" || loading}>Ergebnisse abrufen</button>
</div>


{#if loading}
    <Loader/>
{/if}
{#if results && results.length > 0}
    <div class="table">
        <table>
            <thead>
            <tr>
                <th>IP</th>
                <th>Zeitstempel</th>
                <th>Methode</th>
                <th>Pfad</th>
                <th>HTTP-Version</th>
                <th>Status</th>
                <th>Größe</th>
            </tr>
            </thead>
            <tbody>
            {#each results as result}
                <tr>
                    <td>{result.IP}</td>
                    <td>{result.Timestamp}</td>
                    <td>{result.Method}</td>
                    <td>{result.Path}</td>
                    <td>{result.HttpVesion}</td>
                    <td>{result.Status}</td>
                    <td>{result.Size}</td>
                </tr>
            {/each}
            </tbody>
        </table>
    </div>
{/if}

<style>
    .table {
        overflow:auto;
        width: calc(100% - 20px);
        height: 500px;
        margin: 10px;
    }

    table {
        min-width: 100%;
    }

    thead {
        position: sticky;
        top: 0;
        background-color: black;
    }
</style>