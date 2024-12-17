<script lang="ts">
    import {QueryStatusCounts} from '../../wailsjs/go/main/App.js'
    import {database} from '../../wailsjs/go/models'
    import FilterPicker from "./FilterPicker.svelte";
    import Loader from "./Loader.svelte";

    let results: database.StatusCount[] = []

    export let dbFile = ""

    let timestampFrom = ""
    let timestampTo = ""
    let ip = ""
    let codes = ""

    let loading = false

    function query() {
        loading = true
        results = []

        QueryStatusCounts(dbFile, ip, timestampFrom, timestampTo, codes).then((incomingResults) => {
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
    <FilterPicker withCodes={true} bind:codes={codes} bind:timestampFrom={timestampFrom} bind:timestampTo={timestampTo} bind:ipAddress={ip}/>
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
                <th>Status</th>
                <th>Anzahl</th>
            </tr>
            </thead>
            <tbody>
            {#each results as result}
                <tr>
                    <td>{result.Status}</td>
                    <td>{result.Count}</td>
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

