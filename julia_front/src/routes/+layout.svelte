<script>
	import { onMount } from 'svelte';
	import '../app.css';

	onMount(() => {
		const setVh = () => {
			const h = window.visualViewport ? window.visualViewport.height : window.innerHeight;
			document.documentElement.style.setProperty('--vh', `${h}px`);
		};
		setVh();
		window.addEventListener('resize', setVh);
		window.visualViewport?.addEventListener('resize', setVh);
		if ('serviceWorker' in navigator) {
			navigator.serviceWorker.register('/sw.js').catch(console.error);
		}
	});
</script>

<slot />
