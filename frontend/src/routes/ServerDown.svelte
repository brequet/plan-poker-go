<script lang="ts">
	import { FrownIcon } from 'lucide-svelte';
</script>

	<div class="flex items-center justify-center bg-white p-6 rounded-lg shadow-lg">
		<div class="flex items-center justify-center">
			<FrownIcon class="ml-auto" color="black" />
		</div>
		<div class="mt-3 text-center sm:mt-0 sm:ml-4 sm:text-left ">
			<h3 class="text-lg sm:text-2xl leading-6 font-medium text-gray-900" id="modal-title">
				Service Unavailable
			</h3>
			<div class="mt-2">
				<p class="text-sm text-gray-500">
					We're sorry, but the service is currently unavailable. Please try again later.
				</p>
			</div>
		</div>
	</div>
