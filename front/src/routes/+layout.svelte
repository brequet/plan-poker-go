<script>
	import '../app.css';
</script>

<header class="bg-blue-500 py-4 flex flex-row">
	<div class="container mx-auto px-4">
		<h1 class="text-white text-4xl font-bold">Poker Planning App</h1>
		<p class="text-white mt-2">Plan and estimate together!</p>
	</div>

    <nav>
        <a href='/'>HOME</a>
    </nav>
</header>

<main class="h-full">
	<slot />
</main>

<footer class="bg-blue-500 py-2 text-white text-center">
	<p>© 2023 No Company. No right reserved.</p>
</footer>
