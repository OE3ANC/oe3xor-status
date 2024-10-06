<script lang="ts">
    import * as Table from "$lib/components/ui/table/index.js";
    import ThemeSwitcher from "$lib/components/custom/ThemeSwitcher.svelte";
    import * as Card from "$lib/components/ui/card/index.js";
    import { onMount } from 'svelte';
    import PocketBase from 'pocketbase'

    const pb = new PocketBase('/');


    let qsos = []

    onMount(async () => {

        pb.collection('mmdvm_qso_m17').getList(1, 10, {
            sort: "-created"
        }).then(record => {
            qsos = record.items;
        });

        pb.collection('mmdvm_qso_m17').subscribe('*', function (e) {
            pb.collection('mmdvm_qso_m17').getList(1, 10, {
                sort: "-created"
            }).then(record => {
                qsos = record.items;
            });
        }, { /* other options like expand, custom headers, etc. */});
    })
</script>

<Card.Root class="">
    <Card.Header>
        <Card.Title>M17 Log</Card.Title>
    </Card.Header>
    <Card.Content>
        <Table.Root>
            <Table.Header>
                <Table.Row>
                    <Table.Head class="w-[100px]">FROM</Table.Head>
                    <Table.Head class="w-[100px]">TO</Table.Head>
                    <Table.Head class="w-[150px]">Bit Error Rate</Table.Head>
                    <Table.Head class="w-[100px]">Source</Table.Head>
                    <Table.Head class="w-[100px]">Type</Table.Head>
                    <Table.Head class="w-[100px]">Duration</Table.Head>
                </Table.Row>
            </Table.Header>
            <Table.Body>
                {#each qsos as qso, i (i)}
                    <Table.Row>
                        <Table.Cell class="font-medium">{qso.source_callsign}</Table.Cell>
                        <Table.Cell class="font-medium">{qso.destination_callsign}</Table.Cell>
                        <Table.Cell>{qso.ber}</Table.Cell>
                        <Table.Cell>{qso.source}</Table.Cell>
                        <Table.Cell>{qso.traffic_type}</Table.Cell>
                        <Table.Cell>{qso.duration}</Table.Cell>
                    </Table.Row>
                {/each}
            </Table.Body>
        </Table.Root>
    </Card.Content>
    <Card.Footer class="flex justify-between">
    </Card.Footer>
</Card.Root>

