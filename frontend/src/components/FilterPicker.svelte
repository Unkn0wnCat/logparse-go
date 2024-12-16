<script lang="ts">
    export let ipAddress = ""
    export let timestampFrom = ""
    export let timestampTo = ""

    export let withCodes = false
    export let codes = ""
</script>

<div class="filters">
    <div class="filter">
        <span>IP-Addresse</span>
        <input type="text" bind:value={ipAddress} />
    </div>
    <div class="filter">
        <span>Von</span>
        <input type="datetime-local" bind:value={timestampFrom} step="1" />
    </div>
    <div class="filter">
        <span>Bis</span>
        <input type="datetime-local" bind:value={timestampTo} step="1" />
    </div>
    {#if withCodes}
        <div class="filter">
            <span>Codes</span>
            <input type="text" bind:value={codes} placeholder="200,404" />
        </div>
    {/if}
</div>

<style>
    .filter {
        display: flex;
    }

    .filter span {
        width: 130px;
        text-align: right;
        padding-right: 10px;
    }

    .filter input {
        flex-grow: 1;
    }

    .filters {
        padding: 10px;
    }
</style>