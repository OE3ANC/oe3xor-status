<script lang="ts">
    import { onMount } from 'svelte';

    import { Button } from "$lib/components/ui/button/index.js";
    import * as Card from "$lib/components/ui/card/index.js";
    import { Label } from "$lib/components/ui/label/index.js";
    import ThemeSwitcher from "$lib/components/custom/ThemeSwitcher.svelte";
    import PocketBase from 'pocketbase'

    const pb = new PocketBase('/');

    let lastMode = "No Data" , lastTimestamp = ""

    onMount(async () => {

        pb.collection('mmdvm_mode').getList(1, 1, {
            sort: "-created"
        }).then(record => {
            console.log("Subscription update:", record.items[0])
            lastMode = record.items[0].mode;
            lastTimestamp = record.items[0].created;
        });

        pb.collection('mmdvm_mode').subscribe('*', function (e) {
            lastMode = e.record.mode;
            lastTimestamp = e.record.created;
        }, { /* other options like expand, custom headers, etc. */});

    })
</script>

<Card.Root class="w-[350px]">
    <Card.Header>
        <Card.Title>OE3XOR Mode</Card.Title>
    </Card.Header>
    <Card.Content>
        <h1 title="{lastTimestamp}">{lastMode}</h1>
    </Card.Content>
    <Card.Footer class="flex justify-between">
        <ThemeSwitcher></ThemeSwitcher>
    </Card.Footer>
</Card.Root>