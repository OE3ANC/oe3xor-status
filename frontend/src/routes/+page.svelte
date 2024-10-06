<script lang="ts">

import M17Log from "$lib/components/custom/M17Log.svelte";
import MMDVMStatus from "$lib/components/custom/MMDVMStatus.svelte";
import { ScrollArea } from "$lib/components/ui/scroll-area/index.ts";
import ThemeSwitcher from "$lib/components/custom/ThemeSwitcher.svelte";
</script>

<div class="grid grid-cols-5 grid-rows-1 gap-4">
    <div ><MMDVMStatus></MMDVMStatus></div>
    <div class="col-span-4">
            <M17Log></M17Log>

    </div>
</div>
