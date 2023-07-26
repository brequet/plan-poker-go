<script>
	import '../app.css';
</script>

<header class="bg-blue-500 py-4 flex flex-row">
	<div class="container mx-auto px-4">
		<h1 class="text-white text-4xl font-bold">Poker Planning App</h1>
		<p class="text-white mt-2">Plan and estimate together!</p>
	</div>

	<nav class="flex items-center justify-between bg-blue-500 py-4 px-8">
		<a
			href="/"
			class="text-white text-lg font-semibold py-2 px-4 hover:underline hover:text-gray-100"
		>
			Home
		</a>
	</nav>
</header>

<main class="flex-1 bg-gray-100 p-4">
	<slot />
</main>

<footer class="bg-blue-500 py-2 text-white text-center">
	<p>© 2023 No Company. No right reserved.</p>
</footer>
