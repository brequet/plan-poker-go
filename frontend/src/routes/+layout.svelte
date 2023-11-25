<script>
	import '../app.css';

	const navItems = [
		{ name: 'Home', href: '/' },
		{ name: 'About', href: '/about' }
	];
</script>

<header class="bg-blue-500 py-2 flex flex-row">
	<div class="container mx-auto px-4 flex items-center">
		<a href="/" class="">
			<h1 class="text-white text-5xl" style="font-family: 'Montserrat Regular Alt1';">ESTIMAKE</h1>
		</a>
	</div>

	<nav class="sm:flex items-center justify-between bg-blue-500 py-4 sm:px-8">
		{#each navItems as navItem}
			<a
				href={navItem.href}
				class="text-white sm:text-lg font-semibold py-2 px-4 hover:underline hover:text-gray-100"
			>
				{navItem.name}
			</a>
		{/each}
	</nav>
</header>

<main class="flex-1 bg-gray-100 p-4">
	<slot />
</main>

<footer class="bg-blue-500 py-2 text-white text-center">
	<p>2023 No Company. No right reserved.</p>
</footer>
