<script lang="ts">
	import { Bar } from 'svelte-chartjs';

	export let average;
	export let estimatesGraphData;

	let graphOptions = {
		scales: {
			y: {
				ticks: {
					precision: 0
				}
			}
		}
	};
</script>

<!-- Stats half -->
<h3 class="text-xl font-bold mb-2">Stats</h3>

Average :<span class="font-bold">{average}</span>

<Bar data={estimatesGraphData} options={graphOptions} />
